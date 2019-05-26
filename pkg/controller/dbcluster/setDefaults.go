package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) setUpDefaultsIfNeeded(cr *kubev1alpha1.DBCluster) error {

	if err := r.setUpCredentialsIfNeeded(cr); err != nil {
		return err
	}

	if err := r.setCRDeleteClusterID(cr); err != nil {
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
