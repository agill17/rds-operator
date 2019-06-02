package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) deleteAndUpdateState(cr *kubev1alpha1.DBCluster, clusterObj rdsLib.RDS) error {
	logrus.Warnf("Namespace: %v | CLuster CR: %v | Delete Event detected", cr.Namespace, cr.Name)
	err := clusterObj.Delete()
	if err != nil {
		return err
	}
	cr.SetFinalizers([]string{})
	if err := lib.UpdateCr(r.client, cr); err != nil {
		return err
	}
	return nil
}
