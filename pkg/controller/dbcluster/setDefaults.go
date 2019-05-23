package dbcluster

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

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

	if err := r.setMoreTags(cr); err != nil {
		return err
	}

	return nil
}

// used when deleteSpec.DBClusterID is not set, than use the one provided within cr.createClusterSpec
// this is so we can make the clusterID optional in deleteSpec
func (r *ReconcileDBCluster) setCRDeleteClusterID(cr *kubev1alpha1.DBCluster) error {
	if cr.Spec.DeleteSpec.DBClusterIdentifier == nil {
		id := *cr.Spec.CreateClusterSpec.DBClusterIdentifier
		cr.Spec.DeleteSpec.DBClusterIdentifier = &id
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR while setting up DeleteSpec.DBClusterIdentifier: %v", err)
			return err
		}
	}
	return nil
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

func (r *ReconcileDBCluster) setMoreTags(cr *kubev1alpha1.DBCluster) error {
	if exists := r.checkIfTagContainsCRName(cr); !exists {
		tags := []*rds.Tag{
			{Key: aws.String("CLUSTER_CR_NAME"), Value: aws.String(cr.Name)},
		}
		cr.Spec.CreateClusterSpec.Tags = append(cr.Spec.CreateClusterSpec.Tags, tags...)
		if err := lib.UpdateCr(r.client, cr); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileDBCluster) checkIfTagContainsCRName(cr *kubev1alpha1.DBCluster) bool {
	exists := false
	for _, e := range cr.Spec.CreateClusterSpec.Tags {
		if *e.Key == "CLUSTER_CR_NAME" {
			exists := true
			return exists
		}
	}
	return exists
}
