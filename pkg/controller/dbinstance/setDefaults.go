package dbinstance

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/utils"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBInstance) setUsername(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.CreateInstanceSpec.MasterUsername == nil {
		// meaning this is a standalone deployment
		u := utils.RandStringBytes(9)
		cr.Spec.CreateInstanceSpec.MasterUsername = &u
		if err := utils.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBInstance CR while setting up username: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBInstance) setRegion(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.Region == "" {
		cr.Spec.Region = "us-east-1"
		if err := utils.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBInstance CR while setting up password: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBInstance) setPassword(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.CreateInstanceSpec.MasterUserPassword == nil {
		// meaning this is a standalone deployment
		u := utils.RandStringBytes(9)
		cr.Spec.CreateInstanceSpec.MasterUserPassword = &u
		if err := utils.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBInstance CR while setting up password: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBInstance) setDeleteInsID(cr *kubev1alpha1.DBInstance, actionType rdsLib.RDSAction) error {

	switch actionType {
	case rdsLib.CREATE:
		cr.Spec.DeleteInstanceSpec.DBInstanceIdentifier = cr.Spec.CreateInstanceSpec.DBInstanceIdentifier
	case rdsLib.RESTORE:
		cr.Spec.DeleteInstanceSpec.DBInstanceIdentifier = cr.Spec.RestoreInstanceFromSnap.DBInstanceIdentifier
	}

	if err := utils.UpdateCrStatus(r.client, cr); err != nil {
		return err
	}

	return utils.UpdateCr(r.client, cr)
}

// only used when instance is not associated to cluster
func (r *ReconcileDBInstance) setCRDefaultsIfNeeded(cr *kubev1alpha1.DBInstance, actionType rdsLib.RDSAction) error {
	if err := r.setDeleteInsID(cr, actionType); err != nil {
		return err
	}

	// when NOT associated with cluster, username and password
	// must be part of CreateDBInstanceInput
	if !isPartOfCluster(cr) {

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

	return nil
}
