package dbsubnetgroup

import (
	"context"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBSubnetGroup) updateCrStatus(cr *agillv1alpha1.DBSubnetGroup) error {
	if err := r.client.Status().Update(context.TODO(), cr); err != nil {
		logrus.Errorf("Failed to update cr status in namespace: %v -- ERROR: %v", cr.Namespace, err)
		return err
	}
	return nil
}
