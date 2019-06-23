package dbcluster

import (
	"errors"

	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

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
	if cr.Spec.CreateClusterSpec != nil {
		return *cr.Spec.CreateClusterSpec.DBClusterIdentifier
	} else if cr.Spec.CreateClusterFromSnapshot != nil {
		return *cr.Spec.CreateClusterFromSnapshot.DBClusterIdentifier
	}
	return ""
}

// when useCredentialsFrom is true, no need to deploy a new secret
// else deploy a secret
func useCredentialsFrom(cr *kubev1alpha1.DBCluster) bool {
	if cr.Spec.CredentialsFrom.UsernameKey != "" && cr.Spec.CredentialsFrom.PasswordKey != "" {
		return true
	}
	return false
}
