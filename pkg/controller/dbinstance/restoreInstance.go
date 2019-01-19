package dbinstance

import (

	// h "cloud.google.com/go/bigquery/benchmarks"
	"context"
	"strings"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/controller/dbcluster"
	"github.com/agill17/rds-operator/pkg/controller/lib"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileDBInstance) restoreFromSnapshot(cr *agillv1alpha1.DBInstance, dbID string, request reconcile.Request) error {

	var err error
	instance := &agillv1alpha1.DBInstance{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		logrus.Errorf("Error while getting instance cr from cluster: %v", err)
		return err
	}
	var snapID string

	newDBIDExists, _ := r.dbInstanceExists(dbID)
	if cr.Spec.RehealFromLatestSnapshot && !newDBIDExists {

		logrus.Warnf("Namespace: %v | DB Identifier: %v | Msg: re-creating new rds instance from latest available snapshot available", cr.Namespace, dbID)
		// if this was a cluster
		if instance.Spec.DBClusterIdentifier != "" {
			snapID, _ = lib.GetLatestClusterSnapID(instance.Spec.DBClusterIdentifier, instance.Namespace, "us-east-1")
			clusterObj := &agillv1alpha1.DBCluster{}
			err := r.client.Get(context.TODO(), types.NamespacedName{
				Name:      strings.Replace(instance.Spec.DBClusterIdentifier, instance.Namespace+"-", "", 1),
				Namespace: instance.Namespace}, clusterObj)
			if err != nil {
				logrus.Errorf("Error while getting  cluster spec cr when rehealing from snapshot: %v", err)
			}
			restoreClusterInstanceInput := dbcluster.GetRestoreClusterDBFromSnapInput(clusterObj, *clusterObj.Status.RDSClusterStatus.DBClusters[0].DBClusterIdentifier, snapID)
			if _, err := r.rdsClient.RestoreDBClusterFromSnapshot(restoreClusterInstanceInput); err != nil {
				logrus.Errorf("Error while re-healing db cluster instance from snapshot: %v", err)
			}
			r.createNewDBInstance(instance, dbID, r.createDBInstanceInput(instance, dbID), request)
		} else {
			snapID, _ = lib.GetLatestSnapID(dbID, instance.Namespace)
			_, err = r.rdsClient.RestoreDBInstanceFromDBSnapshot(r.restoreFromSnapInput(instance, dbID, snapID))
			if err != nil {
				logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: ERROR While restoring db from snapshot: %v", cr.Namespace, dbID, err)
				return err
			}
		}

		lib.WaitForExistence("available", dbID, cr.Namespace, r.rdsClient)
		_, cr.Status.RDSInstanceStatus = r.dbInstanceExists(dbID)
		instance.Status.RestoredFromSnapshotName = snapID
		instance.Status.UpdateKubeFiles = true
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return err
		}
	}

	return err
}

// restore instance from s3 -- TODO
