package dbinstance

import (
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"

	"github.com/agill17/rds-operator/pkg/lib"

	// h "cloud.google.com/go/bigquery/benchmarks"
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func (r *ReconcileDBInstance) deleteAndUpdateState(cr *kubev1alpha1.DBInstance, insObj rdsLib.RDS) error {
	logrus.Warnf("Namespace: %v | Instance CR: %v | Delete Event detected", cr.Namespace, cr.Name)

	err := insObj.Delete()
	if err != nil {
		return err
	}

	cr.SetFinalizers([]string{})

	return lib.UpdateCr(r.client, cr)
}
