package dbinstance

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func validateSpecBasedOnType(cr *kubev1alpha1.DBInstance, actionType rdsLib.RDSAction) error {

	switch actionType {
	case rdsLib.CREATE:
		return validateCreateInstanceSpec(cr)
	case rdsLib.RESTORE:
		return validateCreateFromSnapSpec(cr)
	}
	return validateDeleteInstanceSpec(cr)
}

func validateCreateInstanceSpec(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.CreateInstanceSpec == nil {
		logrus.Errorf("Namespace: %v | CR Name: %v | spec.createInstanceSpecEmptyError", cr.Namespace, cr.Name)
		return errors.New("spec.createInstanceSpecEmptyError")
	}
	return nil
}

func validateCreateFromSnapSpec(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.RestoreInstanceFromSnap == nil {
		logrus.Errorf("Namespace: %v | CR Name: %v | spec.createInstanceFromSnapshotEmptyError", cr.Namespace, cr.Name)
		return errors.New("spec.createInstanceFromSnapshotEmptyError")
	}
	return nil
}

func validateDeleteInstanceSpec(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.DeleteInstanceSpec == nil {
		logrus.Errorf("Namespace: %v | CR Name: %v | spec.deleteInstanceSpecEmptyError", cr.Namespace, cr.Name)
		return errors.New("spec.deleteInstanceSpecEmptyError")
	}
	return nil
}
