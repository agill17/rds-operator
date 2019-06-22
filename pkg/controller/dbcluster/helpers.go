package dbcluster

import (
	"context"
	"errors"

	"github.com/agill17/rds-operator/pkg/rdsLib"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func getClusterSecretName(crName string) string {
	return crName + "-secret"
}

func validateRequiredInput(cr *kubev1alpha1.DBCluster) error {
	if cr.Spec.CreateClusterSpec == nil && cr.Spec.CreateClusterFromSnapshot == nil {
		return errors.New("CreateClusterSpecEmptyError")
	}

	if cr.Region == "" {
		return errors.New("regionCannotBeEmptyError")
	}

	if cr.Spec.DeleteSpec == nil {
		return errors.New("deleteClusterSpecCannotBeEmptyError")
	}

	return nil
}

func getActionType(cr *kubev1alpha1.DBCluster) rdsLib.RDSAction {
	if cr.GetDeletionTimestamp() != nil && len(cr.GetFinalizers()) > 0 {
		return rdsLib.DELETE
	} else if cr.Spec.CreateClusterFromSnapshot != nil {
		return rdsLib.RESTORE
	} else if cr.Spec.CreateClusterSpec != nil {
		return rdsLib.CREATE
	}

	return rdsLib.UNKNOWN
}

func getDBClusterID(cr *kubev1alpha1.DBCluster, actionType rdsLib.RDSAction) string {
	switch actionType {
	case rdsLib.CREATE:
		return *cr.Spec.CreateClusterSpec.DBClusterIdentifier
	case rdsLib.RESTORE:
		return *cr.Spec.CreateClusterFromSnapshot.DBClusterIdentifier
	case rdsLib.DELETE:
		if cr.Spec.CreateClusterFromSnapshot != nil {
			return *cr.Spec.CreateClusterFromSnapshot.DBClusterIdentifier
		} else if cr.Spec.CreateClusterSpec != nil {
			return *cr.Spec.CreateClusterSpec.DBClusterIdentifier
		}
	}
	return ""
}

func (r *ReconcileDBCluster) createExternalSvc(cr *kubev1alpha1.DBCluster) error {
	svc := getClusterSvc(cr)
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, svc, func(runtime.Object) error {
		controllerutil.SetControllerReference(cr, svc, r.scheme)
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// when useCredentialsFrom is true, no need to deploy a new secret
// else deploy a secret
func useCredentialsFrom(cr *kubev1alpha1.DBCluster) bool {
	if cr.Spec.CredentialsFrom.UsernameKey != "" && cr.Spec.CredentialsFrom.PasswordKey != "" {
		return true
	}
	return false
}
