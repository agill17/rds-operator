package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) restoreClusterFromSnap(cr *kubev1alpha1.DBCluster, installType string) error {

	// decide which restore input and snapID to use?
	/*
		1. If creating a fresh one from an existing snapshot
		2. If restoring a already deployed db by the operator then get the latest snapID
	*/
	var restoreInput *rds.RestoreDBClusterFromSnapshotInput
	var snapID string
	if installType == "newInstallFromSnapshot" {
		restoreInput = cr.CreateFromSnapshot
		snapID = *restoreInput.SnapshotIdentifier
	} else if installType == "newInstall" {
		snapID, _ = getLatestClusterSnapID(*cr.Spec.DBClusterIdentifier, cr.Namespace, cr.Region)
		restoreInput = GetRestoreClusterDBFromSnapInput(cr, *cr.Spec.DBClusterIdentifier, snapID)
	}

	spew.Dump(restoreInput)
	logrus.Infof("Using snapID: %v", snapID)

	if snapID != "" {
		_, err := r.rdsClient.RestoreDBClusterFromSnapshot(restoreInput)
		if err != nil && err.(awserr.Error).Code() != "DBClusterAlreadyExistsFault" {
			logrus.Errorf("Error while re-healing db cluster instance from snapshot: %v", err)
			return err
		}

		// handle restoring/creating and reconcile if still creating
		if err := r.handlePhases(cr); err != nil {
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
	}
	return nil
}
