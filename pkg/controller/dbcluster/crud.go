package dbcluster

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/davecgh/go-spew/spew"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

// responsible for creating, deleting and restoring db cluster based on actionType
func (r *ReconcileDBCluster) crud(cr *kubev1alpha1.DBCluster,
	clusterObj rdsLib.RDS, actionType rdsLib.RDSAction) error {
	clusterID := getDBClusterID(cr, actionType)
	statusCreated := cr.Status.Created

	switch actionType {

	// fresh install
	case rdsLib.CREATE:

		if !statusCreated {
			if err := clusterObj.Create(); err != nil {
				return err
			}
		}

	// delete event
	case rdsLib.DELETE:

		err := clusterObj.Delete()
		if err != nil {
			return err
		}
		cr.SetFinalizers([]string{})
		return lib.UpdateCr(r.client, cr)

	// restore from snapshot
	case rdsLib.RESTORE:

		if !statusCreated {
			if err := clusterObj.Restore(); err != nil {
				return err
			}
		}

	}

	if !statusCreated {
		// return err if not ready in AWS yet
		if err := r.handlePhases(cr, clusterID); err != nil {
			return err
		}
		cr.Status.Created = true
		_, cr.Status.DescriberClusterOutput = lib.DbClusterExists(
			&lib.RDSGenerics{RDSClient: r.rdsClient,
				ClusterID: clusterID})
		if err := lib.UpdateCrStatus(r.client, cr); err != nil {
			return err
		}
		spew.Dump(cr.Status)
	}

	return nil
}
