package dbinstance

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBInstance) setUsername(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.CreateInstanceSpec.MasterUsername == nil {
		// meaning this is a standalone deployment
		u := lib.RandStringBytes(9)
		cr.Spec.CreateInstanceSpec.MasterUsername = &u
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBInstance CR while setting up username: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBInstance) setRegion(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.Region == "" {
		cr.Spec.Region = "us-east-1"
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBInstance CR while setting up password: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBInstance) setPassword(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.CreateInstanceSpec.MasterUserPassword == nil {
		// meaning this is a standalone deployment
		u := lib.RandStringBytes(9)
		cr.Spec.CreateInstanceSpec.MasterUserPassword = &u
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBInstance CR while setting up password: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBInstance) setJobBackOffLimit(cr *kubev1alpha1.DBInstance) error {
	cr.Spec.InitDB.BackOffLimit = 6
	return lib.UpdateCr(r.client, cr)
}

func (r *ReconcileDBInstance) setInitDBJobDefaults(cr *kubev1alpha1.DBInstance) error {
	return r.setJobBackOffLimit(cr)
}

// only used when instance is not associated to cluster
func (r *ReconcileDBInstance) setCRDefaultsIfNeeded(cr *kubev1alpha1.DBInstance) error {

	if isStandAlone(cr) {

		// update username if needed
		if err := r.setUsername(cr); err != nil {
			return err
		}

		// update username if needed
		if err := r.setPassword(cr); err != nil {
			return err
		}

	}

	if err := r.setRegion(cr); err != nil {
		return err
	}

	if cr.Spec.InitDB.Image != "" {
		if err := r.setInitDBJobDefaults(cr); err != nil {
			return err
		}
	}

	return nil
}
