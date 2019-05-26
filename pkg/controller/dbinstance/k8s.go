package dbinstance

import (
	"context"
	"time"

	batchv1 "k8s.io/api/batch/v1"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileDBInstance) getSvcObj(cr *kubev1alpha1.DBInstance) *corev1.Service {
	s := &corev1.Service{
		TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getSvcName(cr),
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:         "ExternalName",
			ExternalName: *cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address,
		},
	}
	// setup ownerReference
	controllerutil.SetControllerReference(cr, s, r.scheme)

	return s
}

func (r *ReconcileDBInstance) getSecretObj(cr *kubev1alpha1.DBInstance, masterUsername, masterPassword, secretName string) *corev1.Secret {
	data := map[string][]byte{
		"DATABASE_ID":       []byte(*cr.Spec.CreateInstanceSpec.DBInstanceIdentifier),
		"DATABASE_ENDPOINT": []byte(*cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address),
	}

	if cr.Spec.CreateInstanceSpec.DBName != nil {
		data["DATABASE_NAME"] = []byte(*cr.Spec.CreateInstanceSpec.DBName)
	}

	// if associated to cluster, than creds are managed by DBCluster secret
	if isStandAlone(cr) {
		data["DATABASE_USER"] = []byte(masterUsername)
		data["DATABASE_PASSWORD"] = []byte(masterPassword)
	}

	s := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: cr.Namespace,
			Labels:    cr.GetLabels(),
		},
		Type: corev1.SecretType("Opaque"),
		Data: data,
	}
	// setup ownerReference
	controllerutil.SetControllerReference(cr, s, r.scheme)

	return s

}

func (r *ReconcileDBInstance) getCreateJobInput(cr *kubev1alpha1.DBInstance, jobCmd []string) *batchv1.Job {
	input := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "init-rds",
			Namespace: cr.Namespace,
			Labels:    cr.GetLabels(),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: "Never",
					Containers: []corev1.Container{{
						Name:    "init-rds-container",
						Image:   cr.Spec.InitDB.Image,
						Command: jobCmd,
						Env: []corev1.EnvVar{
							{
								Name: "DATABASE_USERNAME",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: getSecretName(cr),
										},
										Key: "",
									},
								},
							},
							{
								Name: "DATABASE_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: getSecretName(cr),
										},
										Key: "",
									},
								},
							}},
					}},
				},
			},
		},
	}
	// setup ownerReference
	controllerutil.SetControllerReference(cr, input, r.scheme)

	return input

}

func (r *ReconcileDBInstance) updateK8sFiles(cr *kubev1alpha1.DBInstance) error {

	// update svc
	svcName := getSvcName(cr)
	secretName := getSecretName(cr)
	svcObj, _ := r.getSvcFromCluster(svcName, cr)

	logrus.Infof("Namesapce: %v | DB Identifier: %v | Msg: Updating External Service as DB Endpoint as changed", cr.Namespace, *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier)
	svcObj.Spec.ExternalName = *cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address
	if err := r.client.Update(context.TODO(), svcObj); err != nil {
		logrus.Errorf("Failed while updating service as UpdateKubeFiles was required: %v", err)
	}

	// update secret
	secretObj, _ := r.getSecretFromCluster(cr, secretName)
	logrus.Infof("Namesapce: %v | DB Identifier: %v | Msg: Updating Secret as DB Endpoint as changed", cr.Namespace, *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier)
	secretObj.Data["DATABASE_ENDPOINT"] = []byte(*cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address)
	if err := r.client.Update(context.TODO(), secretObj); err != nil {
		logrus.Errorf("Failed while updating secret as UpdateKubeFiles was required: %v", err)
	}

	cr.Status.UpdateKubeFiles = false
	err := r.updateResourceStatus(cr)
	if err != nil {
		logrus.Errorf("Failed to update cr status for DBInstance: %v", err)
		return err
	}

	return nil
}

func (r *ReconcileDBInstance) getSvcFromCluster(svcName string, cr *kubev1alpha1.DBInstance) (*corev1.Service, error) {
	svcObj := &corev1.Service{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: svcName, Namespace: cr.Namespace}, svcObj); err != nil {
		logrus.Errorf("ERROR While getting externalSvc: %v", err)
		return nil, err
	}
	return svcObj, nil
}

func (r *ReconcileDBInstance) getSecretFromCluster(cr *kubev1alpha1.DBInstance, secretName string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: cr.Namespace}, secret); err != nil {
		logrus.Errorf("ERROR While getting secret: %v", err)
		return nil, err
	}
	return secret, nil
}

func (r *ReconcileDBInstance) createExternalNameSvc(cr *kubev1alpha1.DBInstance) error {

	if err := r.client.Create(context.TODO(), r.getSvcObj(cr)); err != nil && !errors.IsAlreadyExists(err) && !errors.IsForbidden(err) {
		logrus.Errorf("Namespace: %v | Msg: ERROR while creating RDS Service: %v", cr.Namespace, err)
		return err
	}

	return nil
}

func (r *ReconcileDBInstance) createSecret(cr *kubev1alpha1.DBInstance) error {
	secretObj := r.getSecretObj(cr, cr.Status.Username, cr.Status.Password, getSecretName(cr))
	if err := r.client.Create(context.TODO(), secretObj); err != nil && !errors.IsAlreadyExists(err) && !errors.IsForbidden(err) {
		logrus.Errorf("Error while creating secret object: %v", err)
		return err
	}
	return nil
}

// check if defined and proceed with a timeout type loop
func (r *ReconcileDBInstance) createInitDBJob(instance *kubev1alpha1.DBInstance) error {
	var err error

	if instance.Spec.InitDB.Image != "" {
		timeout := instance.Spec.InitDB.Timeout
		jobCmd := instance.Spec.InitDB.Command

		input := r.getCreateJobInput(instance, jobCmd)
		if len(instance.Spec.InitDB.NodeSelector) != 0 {
			input.Spec.Template.Spec.NodeSelector = instance.Spec.InitDB.NodeSelector
		}
		if instance.Spec.InitDB.ImagePullSecret != "" {
			input.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{Name: instance.Spec.InitDB.ImagePullSecret}}
		}

		// 1. Create job
		err = r.client.Create(context.TODO(), input)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Error while creating init-db job: %v", err)
			return err
		}

		// get the job deployed inside cluster
		job := &batchv1.Job{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: "init-rds", Namespace: instance.Namespace}, job)

		// if initDBJob is not marked as successfull and initDBJob has not ran into a timeout yet and cr spec states to waitUntilJobCompleted/timedOut
		if err == nil && !instance.Status.InitJobSuccessfull && instance.Spec.InitDB.WaitTillCompleted && !instance.Status.InitJobTimedOut {

			for start := time.Now(); ; {
				if err := r.client.Get(context.TODO(), types.NamespacedName{Name: "init-rds", Namespace: instance.Namespace}, job); err != nil {
					logrus.Errorf("Unable to find init-rds job inside %v namespace", instance.Namespace)
				} else {
					if job.Status.Succeeded == 1 {
						logrus.Infof("Namespace: %v | Msg: Job was successfully completed before timeout! %v", instance.Namespace, time.Since(start))
						instance.Status.InitJobSuccessfull = true
						err = r.updateResourceStatus(instance)
						if err != nil {
							logrus.Errorf("Failed to update cr status for DBInstance: %v", err)
							return err
						}
						break
					}
				}
				logrus.Infof("Namespace: %v | Msg: Waiting for job to either timeout or be successfull.", instance.Namespace)
				logrus.Infof("Namespace: %v | Msg: sleeping for 5 secs before next job status check", instance.Namespace)
				time.Sleep(5 * time.Second)
				if time.Since(start) > time.Duration(timeout)*time.Second {
					logrus.Errorf("Namespace: %v | Msg: Timed out waiting for job to have 1 pod successful: %v", instance.Namespace, time.Since(start))
					instance.Status.InitJobTimedOut = true
					err = r.updateResourceStatus(instance)
					if err != nil {
						logrus.Errorf("Failed to update cr status for DBInstance: %v", err)
						return err
					}
					break
				}
			}
		}

	}
	return err

}
