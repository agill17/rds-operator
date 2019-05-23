package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

// Used when a CR was used to create a DB cluster, but that no longer exists,
// so restore from latest available snapshot -- meaning installType must be a CREATE_INSTALL_NEW
func (r *ReconcileDBCluster) recoverDeletedClusterFromSnapshot(cr *kubev1alpha1.DBCluster) error {

	snapID, _ := getLatestClusterSnapID(*cr.Spec.DBClusterIdentifier, cr.Namespace, cr.Region)
	restoreInput := GetRestoreClusterDBFromSnapInput(cr, *cr.Spec.DBClusterIdentifier, snapID)
	spew.Dump(restoreInput)
	logrus.Infof("Using snapID: %v", snapID)

	if snapID != "" {
		_, err := r.rdsClient.RestoreDBClusterFromSnapshot(restoreInput)
		if err != nil && err.(awserr.Error).Code() != "DBClusterAlreadyExistsFault" {
			logrus.Errorf("Error while re-healing db cluster instance from snapshot: %v", err)
			return err
		}

		// handle restoring/creating and reconcile if still creating
		if err := r.handlePhases(cr, *cr.Spec.DBClusterIdentifier); err != nil {
			return err
		}

		cr.Status.RestoredFromSnapshotName = snapID
		cr.Status.SecretUpdateNeeded = true
		cr.Status.RestoreNeeded = false
		_, cr.Status.DescriberClusterOutput = lib.DbClusterExists(&lib.RDSGenerics{ClusterID: *cr.Spec.DBClusterIdentifier, RDSClient: r.rdsClient})
		if err := r.updateCrStats(cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR status: %v", err)
			return err
		}
	} else {
		logrus.Errorf("Could not find any latest snapshot for this CR: %v to recover DBCluster", cr.Name)
	}
	return nil

}
