package dbcluster

import (
	"errors"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/lib/dbHelpers"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) restoreAndUpdateState(cr *kubev1alpha1.DBCluster, cluster *dbHelpers.Cluster) error {

	if cluster.RestoreFromSnapInput.SnapshotIdentifier == nil ||
		cluster.RestoreFromSnapInput.DBClusterIdentifier == nil {
		return errors.New("RestoreDBClusterInsufficientParameterError")
	}

	err := dbHelpers.InstallRestoreDelete(cluster, dbHelpers.RESTORE)
	if err != nil {
		logrus.Errorf("Error while re-healing db cluster instance from snapshot: %v", err)
		return err
	}

	// handle restoring/creating and reconcile if still creating
	if err := r.handlePhases(cr); err != nil {
		return err
	}

	cr.Status.RestoredFromSnapshotName = *cluster.RestoreFromSnapInput.SnapshotIdentifier
	cr.Status.SecretUpdateNeeded = true
	cr.Status.RestoreNeeded = false
	_, cr.Status.DescriberClusterOutput = lib.DbClusterExists(
		&lib.RDSGenerics{ClusterID: *cluster.CreateInput.DBClusterIdentifier, RDSClient: r.rdsClient})
	return lib.UpdateCrStatus(r.client, cr)
}
