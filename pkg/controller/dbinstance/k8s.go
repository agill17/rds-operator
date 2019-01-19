package dbinstance

import (
	"context"
	"time"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/controller/lib"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileDBInstance) getSvcObj(cr *agillv1alpha1.DBInstance, dbID, svcName string) *corev1.Service {
	yes := true
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{{
				Name:               cr.Name,
				APIVersion:         "agill.apps/v1alpha1",
				Kind:               "DBInstance",
				UID:                cr.UID,
				BlockOwnerDeletion: &yes,
				Controller:         &yes,
			}},
		},
		Spec: corev1.ServiceSpec{
			Type:         "ExternalName",
			ExternalName: *cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address,
		},
	}
}

func (r *ReconcileDBInstance) getSecretObj(cr *agillv1alpha1.DBInstance, dbID, dbName, masterUsername, masterPassword string) *corev1.Secret {
	o := true
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.DBSecretName,
			Namespace: cr.Namespace,
			Labels:    cr.GetLabels(),
			OwnerReferences: []metav1.OwnerReference{{
				Name:               cr.Name,
				APIVersion:         "agill.apps/v1alpha1",
				Kind:               "DBInstance",
				UID:                cr.GetUID(),
				BlockOwnerDeletion: &o,
				Controller:         &o,
			}},
		},
		Type: corev1.SecretType("Opaque"),
		Data: map[string][]byte{
			"DATABASE_ID":       []byte(dbID),
			"DATABASE_NAME":     []byte(dbName),
			"DATABASE_ENDPOINT": []byte(*cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address),
			"DATABASE_USERNAME": []byte(masterUsername),
			"DATABASE_PASSWORD": []byte(masterPassword),
		},
	}

}

func getCreateJobInput(cr *agillv1alpha1.DBInstance, jobCmd []string) *batchv1.Job {
	o := true
	input := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "init-rds",
			Namespace: cr.Namespace,
			Labels:    cr.GetLabels(),
			OwnerReferences: []metav1.OwnerReference{{
				Name:               cr.Name,
				APIVersion:         "agill.apps/v1alpha1",
				Kind:               "DBInstance",
				UID:                cr.GetUID(),
				BlockOwnerDeletion: &o,
				Controller:         &o,
			}},
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
						Image:   cr.Spec.InitDBJob.Image,
						Command: jobCmd,
					}},
				},
			},
		},
	}

	return input

}

func (r *ReconcileDBInstance) updateK8sFiles(cr *agillv1alpha1.DBInstance, dbID, svcName, secretName string, request reconcile.Request) (reconcile.Result, error) {
	instance := &agillv1alpha1.DBInstance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if instance.Status.UpdateKubeFiles {

		// update svc
		svcObj, _ := r.getSvcFromCluster(request, svcName, cr)

		logrus.Infof("Namesapce: %v | DB Identifier: %v | Msg: Updating External Service as DB Endpoint as changed", cr.Namespace, dbID)
		svcObj.Spec.ExternalName = *cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address

		// update secret
		secretObj, _ := r.getSecretFromCluster(request)

		logrus.Infof("Namesapce: %v | DB Identifier: %v | Msg: Updating Secret as DB Endpoint as changed", cr.Namespace, dbID)
		secretObj.Data["DATABASE_ENDPOINT"] = []byte(*cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address)

		// TODO - update secret

		instance.Status.UpdateKubeFiles = false
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, err
		}

	}

	return reconcile.Result{}, nil
}

func (r *ReconcileDBInstance) getSvcFromCluster(request reconcile.Request, svcName string, cr *agillv1alpha1.DBInstance) (*corev1.Service, error) {
	svcObj := &corev1.Service{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: svcName, Namespace: cr.Namespace}, svcObj); err != nil {
		logrus.Errorf("ERROR While getting externalSvc: %v", err)
		return nil, err
	}
	return svcObj, nil
}

func (r *ReconcileDBInstance) getSecretFromCluster(request reconcile.Request) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, secret); err != nil {
		logrus.Errorf("ERROR While getting secret: %v", err)
		return nil, err
	}
	return secret, nil
}

func (r *ReconcileDBInstance) createExternalNameSvc(cr *agillv1alpha1.DBInstance, dbID string, request reconcile.Request) (reconcile.Result, error) {
	var svcName = cr.Name + "rds-service"
	if cr.Spec.ExternalSvcName != "" {
		svcName = cr.Spec.ExternalSvcName
	}

	if err := r.client.Create(context.TODO(), r.getSvcObj(cr, dbID, svcName)); err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: ERROR while creating RDS Service: %v", cr.Namespace, dbID, err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileDBInstance) createSecret(cr *agillv1alpha1.DBInstance, dbID, dbName, username, password string, request reconcile.Request) (reconcile.Result, error) {
	if err := r.client.Create(context.TODO(), r.getSecretObj(cr, dbID, dbName, username, password)); err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Error while creating secret object: %v", err)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

// check if defined and proceed with a timeout type loop
func (r *ReconcileDBInstance) createInitDBJob(cr *agillv1alpha1.DBInstance, request reconcile.Request) error {
	var err error
	instance := &agillv1alpha1.DBInstance{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return err
	}
	if cr.Spec.InitDBJob.Image != "" {
		timeout := cr.Spec.InitDBJob.Timeout
		if timeout == 0 {
			timeout = lib.DefaultTimeoutForJob
		}
		jobCmd := lib.GetImportJobCmd(
			cr.Spec.Engine,
			cr.Spec.DBName,
			cr.Spec.MasterUsername,
			cr.Spec.MasterPassword,
			*cr.Status.RDSInstanceStatus.DBInstances[0].Endpoint.Address,
			cr.Spec.InitDBJob.SQLFile)

		input := getCreateJobInput(cr, jobCmd)
		if len(cr.Spec.InitDBJob.NodeSelector) != 0 {
			input.Spec.Template.Spec.NodeSelector = cr.Spec.InitDBJob.NodeSelector
		}
		if cr.Spec.InitDBJob.ImagePullSecret != "" {
			input.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{Name: cr.Spec.InitDBJob.ImagePullSecret}}
		}

		/*
			1. Create job
			2. If waitTillCompleted; run a timer loop
			3. Break in loop if timed out or pods.succeeded == 1

		*/

		// 1. Create job
		err = r.client.Create(context.TODO(), input)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Error while creating init-db job: %v", err)
			return err
		}

		// get the job deployed inside cluster
		job := &batchv1.Job{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: "init-rds", Namespace: cr.Namespace}, job)

		// if initDBJob is not marked as successfull and initDBJob has not ran into a timeout yet and cr spec states to waitUntilJobCompleted/timedOut
		if err == nil && !cr.Status.InitJobSuccessfull && cr.Spec.InitDBJob.WaitTillCompleted && !cr.Status.InitJobTimedOut {

			for start := time.Now(); ; {
				if err := r.client.Get(context.TODO(), types.NamespacedName{Name: "init-rds", Namespace: cr.Namespace}, job); err != nil {
					logrus.Errorf("Unable to find init-rds job inside %v namespace", cr.Namespace)
				} else {
					if job.Status.Succeeded == 1 {
						logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Job was successfully completed before timeout! %v", cr.Namespace, lib.SetDBID(cr.Namespace, cr.Name), time.Since(start))
						instance.Status.InitJobSuccessfull = true
						if err := r.client.Update(context.TODO(), instance); err != nil {
							return err
						}
						break
					}
				}
				logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Waiting for job to either timeout or be successfull.", cr.Namespace, lib.SetDBID(cr.Namespace, cr.Name))
				logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: sleeping for 5 secs before next job status check", cr.Namespace, lib.SetDBID(cr.Namespace, cr.Name))
				time.Sleep(5 * time.Second)
				if time.Since(start) > time.Duration(timeout)*time.Second {
					logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: Timed out waiting for job to have 1 pod successful: %v", cr.Namespace, lib.SetDBID(cr.Namespace, cr.Name), time.Since(start))
					instance.Status.InitJobTimedOut = true
					if err := r.client.Update(context.TODO(), instance); err != nil {
						return err
					}
					break
				}
			}
		}

	}
	return err

}
