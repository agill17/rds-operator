package dbcluster

import (
	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) restoreAndUpdateState(cr *kubev1alpha1.DBCluster, cluster rdsLib.RDS) error {

	err := cluster.Delete()
	if err != nil {
		logrus.Errorf("Error while re-healing db cluster instance from snapshot: %v", err)
		return err
	}

	// handle restoring/creating and reconcile if still creating
	if err := r.handlePhases(cr, *cr.Spec.CreateClusterFromSnapshot.DBClusterIdentifier); err != nil {
		return err
	}
	cr.Status.Created = true
	cr.Status.RestoredFromSnapshotName = *cr.Spec.CreateClusterFromSnapshot.SnapshotIdentifier
	_, cr.Status.DescriberClusterOutput = lib.DbClusterExists(
		&lib.RDSGenerics{ClusterID: *cr.Spec.CreateClusterFromSnapshot.DBClusterIdentifier, RDSClient: r.rdsClient})
	return lib.UpdateCrStatus(r.client, cr)
}
