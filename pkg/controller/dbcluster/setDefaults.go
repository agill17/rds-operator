package dbcluster

import (
	"fmt"
	"strings"
	"time"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/lib/dbHelpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) setUpDefaultsIfNeeded(cr *kubev1alpha1.DBCluster, installType dbHelpers.DBInstallType) error {

	if err := r.setUpCredentialsIfNeeded(cr); err != nil {
		return err
	}

	if err := r.setCRDeleteClusterID(cr); err != nil {
		return err
	}

	if err := r.setCRDeleteSpecSnapName(cr, installType); err != nil {
		return err
	}

	return nil
}

// used when deleteSpec.DBClusterID is not set, than use the one provided within cr.createClusterSpec
// this is so we can make the clusterID optional in deleteSpec
func (r *ReconcileDBCluster) setCRDeleteClusterID(cr *kubev1alpha1.DBCluster) error {
	if cr.Spec.DeleteSpec.DBClusterIdentifier == nil {
		id := r.getDBClusterIDFromSpec(cr)
		cr.Spec.DeleteSpec.DBClusterIdentifier = &id
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR while setting up DeleteSpec.DBClusterIdentifier: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBCluster) setCRDeleteSpecSnapName(cr *kubev1alpha1.DBCluster, installType dbHelpers.DBInstallType) error {
	var clusterID string

	if !(*cr.Spec.DeleteSpec.SkipFinalSnapshot) && cr.Spec.DeleteSpec.FinalDBSnapshotIdentifier == nil {
		switch installType {
		case dbHelpers.CREATE:
			clusterID = *cr.Spec.CreateClusterSpec.DBClusterIdentifier
		case dbHelpers.RESTORE:
			clusterID = *cr.Spec.CreateClusterFromSnapshot.DBClusterIdentifier
		}

		currentTime := time.Now().Format("2006-01-02:03-02-44")
		snashotName := fmt.Sprintf("%v-%v", clusterID, strings.Replace(currentTime, ":", "-", -1))
		cr.Spec.DeleteSpec.FinalDBSnapshotIdentifier = aws.String(snashotName)

		return lib.UpdateCr(r.client, cr)
	}
	return nil
}

func (r *ReconcileDBCluster) getDBClusterIDFromSpec(cr *kubev1alpha1.DBCluster) string {
	if cr.Spec.CreateClusterSpec != nil {
		return *cr.Spec.CreateClusterSpec.DBClusterIdentifier
	} else if cr.Spec.CreateClusterFromSnapshot != nil {
		return *cr.Spec.CreateClusterFromSnapshot.DBClusterIdentifier
	}
	return ""
}

func (r *ReconcileDBCluster) setCRUsername(cr *kubev1alpha1.DBCluster) error {
	if cr.Spec.CreateClusterSpec.MasterUsername == nil {
		u := lib.RandStringBytes(9)
		cr.Spec.CreateClusterSpec.MasterUsername = &u
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR while setting up username: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBCluster) setCRPassword(cr *kubev1alpha1.DBCluster) error {
	if cr.Spec.CreateClusterSpec.MasterUserPassword == nil {
		p := lib.RandStringBytes(9)
		cr.Spec.CreateClusterSpec.MasterUserPassword = &p
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR while setting up credentials: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBCluster) setUpCredentialsIfNeeded(cr *kubev1alpha1.DBCluster) error {
	if err := r.setCRUsername(cr); err != nil {
		return err
	}

	if err := r.setCRPassword(cr); err != nil {
		return err
	}

	return nil
}
