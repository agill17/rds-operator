package dbinstance

import (

	// h "cloud.google.com/go/bigquery/benchmarks"
	"context"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/controller/lib"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileDBInstance) restoreFromSnapshot(cr *agillv1alpha1.DBInstance, dbID string, request reconcile.Request) error {
	var snapID string
	var err error
	instance := &agillv1alpha1.DBInstance{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		logrus.Errorf("Error while getting instance cr from cluster: %v", err)
		return err
	}

	logrus.Warnf("Namespace: %v | DB Identifier: %v | Msg: re-creating new rds instance from latest available snapshot available", cr.Namespace, dbID)
	// if this was a cluster ( let cluster controller re-create it )
	if instance.Spec.DBClusterIdentifier != "" && instance.Spec.RehealFromLatestSnapshot {
		if _, err := r.createNewDBInstance(instance, dbID, r.createDBInstanceInput(instance, dbID), request); err != nil {
			logrus.Errorf("ERROR while creating new instance as part of cluster reheal policy: %v", err)
			return err
		}
	} else if instance.Spec.RehealFromLatestSnapshot {
		snapID, _ = lib.GetLatestSnapID(dbID, instance.Namespace)
		_, err = r.rdsClient.RestoreDBInstanceFromDBSnapshot(r.restoreFromSnapInput(instance, dbID, snapID))
		if err != nil {
			logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: ERROR While restoring db from snapshot: %v", cr.Namespace, dbID, err)
			return err
		}
	}
	lib.WaitForExistence("available", dbID, cr.Namespace, r.rdsClient)

	return err
}

// restore instance from s3 -- TODO
