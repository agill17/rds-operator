package dbinstance

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func (r *ReconcileDBInstance) crud(cr *kubev1alpha1.DBInstance, actionType rdsLib.RDSAction) error {
	insID := getInstanceID(cr)
	statusCreated := cr.Status.Created
	insObj := rdsLib.NewInstance(
		r.rdsClient,
		cr.Spec.CreateInstanceSpec,
		cr.Spec.DeleteInstanceSpec,
		cr.Spec.RestoreInstanceFromSnap,
		cr, r.client, insID,
	)

	switch actionType {

	// fresh install
	case rdsLib.CREATE:

		if !statusCreated {
			if err := r.waitForClusterIfNeeded(cr); err != nil {
				return err
			}
			if err := insObj.Create(); err != nil {
				return err
			}
		}

	case rdsLib.DELETE:
		logrus.Warnf("Namespace: %v | Instance CR: %v | Delete Event detected", cr.Namespace, cr.Name)

		err := insObj.Delete()
		if err != nil {
			return err
		}

		cr.SetFinalizers([]string{})

		return lib.UpdateCr(r.client, cr)
	case rdsLib.RESTORE:

		if !statusCreated {
			if err := insObj.Restore(); err != nil {
				return err
			}
		}

	}

	if !statusCreated {

		if err := rdsLib.SyncAndReconcileIfNotReady(insObj); err != nil {
			return err
		}

		cr.Status.Created = true
		_, cr.Status.RDSInstanceStatus = lib.DBInstanceExists(
			&lib.RDSGenerics{
				RDSClient:  r.rdsClient,
				InstanceID: insID,
			},
		)
		if err := lib.UpdateCrStatus(r.client, cr); err != nil {
			return err
		}
	}

	return nil
}
