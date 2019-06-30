package dbcluster

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func (r *ReconcileDBCluster) crud(cr *kubev1alpha1.DBCluster, actionType rdsLib.RDSAction) error {
	clusterID := getDBClusterID(cr)
	logrus.Infof("ClusterID ~~> %v", clusterID)
	statusCreated := cr.Status.Created

	// returns cluster struct which is also part of rds interface
	// so we can call all funcs that are part of the interface as long as cluster satifies the interface
	// cluster obj implements and satifies RDS interface by implementing all methods of that interface
	clusterObj := rdsLib.NewCluster(r.rdsClient,
		cr.Spec.CreateClusterSpec,
		cr.Spec.DeleteSpec,
		cr.Spec.CreateClusterFromSnapshot, cr, r.client, clusterID)

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
		if err := rdsLib.SyncAndReconcileIfNotReady(clusterObj); err != nil {
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
