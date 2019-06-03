package dbinstance

import (
	"github.com/agill17/rds-operator/pkg/lib"

	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

// restore instances from snapshot are ALWAYS standalone
func (r *ReconcileDBInstance) restoreAndUpdateState(cr *kubev1alpha1.DBInstance, insObj rdsLib.RDS) error {

	err := insObj.Restore()
	if err != nil {
		return err
	}

	err = r.handlePhases(cr)
	if err != nil {
		return err
	}

	cr.Status.Created = true
	cr.Status.RestoredFromSnapshotName = *cr.Spec.RestoreInstanceFromSnap.DBSnapshotIdentifier

	return lib.UpdateCrStatus(r.client, cr)
}
