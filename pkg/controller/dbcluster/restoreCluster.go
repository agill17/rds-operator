package dbcluster

import (
	"context"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileDBCluster) restoreClusterFromSnap(request reconcile.Request, clusterID, snapID string) error {
	instance := &agillv1alpha1.DBCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		logrus.Errorf("Error while getting DBCluster spec cr when rehealing from snapshot: %v", err)
		return err
	}
	if _, err := r.rdsClient.RestoreDBClusterFromSnapshot(GetRestoreClusterDBFromSnapInput(instance, clusterID, snapID)); err != nil {
		logrus.Errorf("Error while re-healing db cluster instance from snapshot: %v", err)
		return err
	}
	instance.Status.RehealedFromSnapshot = snapID
	return err
}
