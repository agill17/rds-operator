package dbcluster

import (
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

// responsible for creating, deleting and restoring db cluster based on actionType
func (r *ReconcileDBCluster) crud(cr *kubev1alpha1.DBCluster,
	clusterObj rdsLib.RDS, actionType rdsLib.RDSAction) error {

	switch actionType {

	// fresh install
	case rdsLib.CREATE:

		if !cr.Status.Created {

			err := r.createItAndUpdateState(cr, clusterObj)
			if err != nil {
				return err
			}

		}

	// delete event
	case rdsLib.DELETE:
		return r.deleteAndUpdateState(cr, clusterObj)

	// restore from snapshot
	case rdsLib.RESTORE:
		if !cr.Status.Created {
			logrus.Infof("Recreate cluster requested for namespace: %v", cr.Namespace)
			return r.restoreAndUpdateState(cr, clusterObj)
		}
	}

	return nil
}
