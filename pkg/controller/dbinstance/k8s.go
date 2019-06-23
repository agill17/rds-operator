package dbinstance

import (
	"context"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	if !isPartOfCluster(cr) {
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

// func (r *ReconcileDBInstance) getCreateJobInput(cr *kubev1alpha1.DBInstance, jobCmd []string) *batchv1.Job {
// 	input := &batchv1.Job{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "init-rds",
// 			Namespace: cr.Namespace,
// 			Labels:    cr.GetLabels(),
// 		},
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Job",
// 			APIVersion: "batch/v1",
// 		},
// 		Spec: batchv1.JobSpec{
// 			Template: corev1.PodTemplateSpec{
// 				Spec: corev1.PodSpec{
// 					RestartPolicy: "Never",
// 					Containers: []corev1.Container{{
// 						Name:    "init-rds-container",
// 						Image:   cr.Spec.InitDB.Image,
// 						Command: jobCmd,
// 						Env: []corev1.EnvVar{
// 							{
// 								Name: "DATABASE_USERNAME",
// 								ValueFrom: &corev1.EnvVarSource{
// 									SecretKeyRef: &corev1.SecretKeySelector{
// 										LocalObjectReference: corev1.LocalObjectReference{
// 											Name: getSecretName(cr),
// 										},
// 										Key: "",
// 									},
// 								},
// 							},
// 							{
// 								Name: "DATABASE_PASSWORD",
// 								ValueFrom: &corev1.EnvVarSource{
// 									SecretKeyRef: &corev1.SecretKeySelector{
// 										LocalObjectReference: corev1.LocalObjectReference{
// 											Name: getSecretName(cr),
// 										},
// 										Key: "",
// 									},
// 								},
// 							}},
// 					}},
// 				},
// 			},
// 		},
// 	}
// 	// setup ownerReference
// 	controllerutil.SetControllerReference(cr, input, r.scheme)

// 	return input

// }

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
