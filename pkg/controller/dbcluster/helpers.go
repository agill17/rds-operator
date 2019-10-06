package dbcluster

import (
	"github.com/agill17/rds-operator/pkg/utils"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/davecgh/go-spew/spew"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)


func (r *ReconcileDBCluster)getActionType(cr *kubev1alpha1.DBCluster) rdsLib.RDSAction {
	if cr.GetDeletionTimestamp() != nil && len(cr.GetFinalizers()) > 0 {
		return rdsLib.DELETE
	} else if cr.ClusterSpec.SnapshotIdentifier != nil {
		return rdsLib.RESTORE
	}

	return rdsLib.CREATE
}

func getDBClusterID(cr *kubev1alpha1.DBCluster) string {
	return *cr.ClusterSpec.DBClusterIdentifier
}

// when useCredentialsFrom is true, no need to deploy a new secret
// else deploy a secret
func useCredentialsFrom(cr *kubev1alpha1.DBCluster) bool {
	if cr.ClusterSpec.CredentialsFrom.UsernameKey != "" && cr.ClusterSpec.CredentialsFrom.PasswordKey != "" {
		return true
	}
	return false
}

func (r *ReconcileDBCluster) updateClusterStatusInCr(cr *kubev1alpha1.DBCluster) error {
	var err error
	if !cr.Status.Created {
		cr.Status.Created = true
		_, cr.Status.DescriberClusterOutput, err = utils.DbClusterExists(utils.RDSGenerics{RDSClient: r.rdsClient, ClusterID: getDBClusterID(cr)})
		if err != nil {
			return err
		}
		spew.Dump(cr.Status)
		return utils.UpdateCrStatus(r.client, cr)
	}
	return nil
}
