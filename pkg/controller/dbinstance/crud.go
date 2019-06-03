package dbinstance

import (
	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func (r *ReconcileDBInstance) crud(cr *kubev1alpha1.DBInstance, insObj rdsLib.RDS, actionType rdsLib.RDSAction) error {

	switch actionType {

	case rdsLib.CREATE:
		if !cr.Status.Created {
			if _, err := r.createNewDBInstance(cr, insObj); err != nil {
				return err
			}
		}

	case rdsLib.DELETE:
		return r.deleteAndUpdateState(cr, insObj)
	case rdsLib.RESTORE:
		return r.restoreAndUpdateState(cr, insObj)
	}

	return nil
}
