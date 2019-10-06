package dbsubnetgroup

import (
	"fmt"

	"github.com/sirupsen/logrus"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/utils"
)

func (r *ReconcileDBSubnetGroup) setCRDefaultsIfNeeded(cr *agillv1alpha1.DBSubnetGroup) error {

	if cr.Spec.DBSubnetGroupName == nil {
		cr.Spec.DBSubnetGroupName = &cr.Name
		logrus.Warnf("DBSubnetGroupName is empty, will try to create one using cr.Name: %v", cr.Name)
		if err := utils.UpdateCr(r.client, cr); err != nil {
			return err
		}
	}

	if cr.Spec.DBSubnetGroupDescription == nil {
		desc := fmt.Sprintf("CustomResoure Name: %v Inside Namespace: %v", cr.Name, cr.Namespace)
		cr.Spec.DBSubnetGroupDescription = &desc
		if err := utils.UpdateCr(r.client, cr); err != nil {
			return err
		}
	}

	return nil
}
