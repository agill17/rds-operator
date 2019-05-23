package dbsubnetgroup

import (
	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBSubnetGroup) createSubnetGroup(cr *agillv1alpha1.DBSubnetGroup) error {

	logrus.Infof("Creating DBSubnetGroup %v in namespace: %v", cr.Name, cr.Namespace)
	_, err := r.rdsClient.CreateDBSubnetGroup(cr.Spec)
	if err != nil {
		logrus.Errorf("Something went wrong while creating db subnet group: %v", err)
		return err
	}

	// update status
	cr.Status.Created = true
	if err := r.updateCrStatus(cr); err != nil {
		return err
	}
	return nil

}
