package dbcluster

import (
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
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
				switch err.(type) {
				case *lib.ErrorResourceCreatingInProgress:
					// print this out so its helpful when looking at logs
					logrus.Errorf("Namespace: %v | CR: %v | Msg: Cluster still in creating phase. Reconciling to check again.", cr.Namespace, cr.Name)
				}
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
