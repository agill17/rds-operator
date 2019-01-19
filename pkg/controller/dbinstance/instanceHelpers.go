package dbinstance

import (
	"context"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/controller/lib"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func crHasDBStatus(cr *agillv1alpha1.DBInstance) bool {
	var crHasSTatus bool
	if cr.Status.DeployedInitially {
		// cr has db status -- meaning db had already been deployed for this cr.
		crHasSTatus = true
	}
	return crHasSTatus
}

func (r *ReconcileDBInstance) dbInstanceExists(dbIdentifier string) (bool, *rds.DescribeDBInstancesOutput) {
	exists := true

	output, err := r.rdsClient.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &dbIdentifier,
	})

	if err != nil && err.(awserr.Error).Code() == rds.ErrCodeDBInstanceNotFoundFault {
		exists = false
	}
	return exists, output
}

/*
	Cases to handle:
	1. DB exists in AWS AND CR has no status about it -- dont do anything but just log that msg
	2. DB does not exist in AWS AND CR has no status about any deployment -- create a new fresh DB
	3. DB does not exist in AWS AND CR has status that it was deployed atleast once -- reheal from snapshot if defined to do so, notify
	4. DB does exists in AWS AND CR has status about it -- dont do anything but just log that msg
*/
func (r *ReconcileDBInstance) createDbInstanceIfNotExists(cr *agillv1alpha1.DBInstance, dbID, ns string, request reconcile.Request) error {
	var err error

	dbInput := r.createDBInstanceInput(cr, dbID)
	dbExistsInAws, _ := r.dbInstanceExists(dbID)
	crHasStatus := crHasDBStatus(cr)

	// ALL cases from above
	if dbExistsInAws && !crHasStatus {
		// 1. DB exists in AWS AND CR has no status about it -- dont do anything but just log that msg
		logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: dbID already exists in AWS! Please create a new DB by updating the name under metadata of the cr since the operator is not sure if it deployed it", cr.Namespace, dbID)
	} else if !dbExistsInAws && !crHasStatus {
		// 2. DB does not exist and CR has no status about any deployment -- create a new fresh DB
		if _, err := r.createNewDBInstance(cr, dbID, dbInput, request); err != nil {
			return err
		}
	} else if !dbExistsInAws && crHasStatus {
		// 3. DB does not exist in AWS AND CR has status that it was deployed atleast once -- reheal from snapshot if defined to do so, notify
		logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: CR is marked as db is deployed but does not exist in AWS!!!", cr.Namespace, dbID)
		if err := r.restoreFromSnapshot(cr, dbID, request); err != nil {
			return err
		}
	} else if dbExistsInAws && crHasStatus {
		// 4. DB exists in AWS AND CR has status about it -- dont do anything but just log that msg
		logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: CR has DB marked as deployed and also exists in aws", cr.Namespace, dbID)

	}
	return err
}

func (r *ReconcileDBInstance) createNewDBInstance(cr *agillv1alpha1.DBInstance, dbID string, dbInput *rds.CreateDBInstanceInput, request reconcile.Request) (*rds.CreateDBInstanceOutput, error) {
	var createOut *rds.CreateDBInstanceOutput
	var err error
	logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Creating DB", cr.Namespace, dbID)
	if cr.Spec.DBClusterIdentifier != "" {
		logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Part of cluster: %v -- checking if its available first", cr.Namespace, dbID, cr.Spec.DBClusterIdentifier)
		// lib.WaitForExistence("available", cr.Spec.DBClusterIdentifier, cr.Namespace, r.rdsClient)
	}
	createOut, err = r.rdsClient.CreateDBInstance(dbInput)

	if err != nil {
		logrus.Errorf("Error while creating DB: %v", err)
		return createOut, err
	} else {
		spew.Dump(createOut)
		err = lib.WaitForExistence("available", dbID, cr.Namespace, r.rdsClient)
		logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Updating CR status with DB status...", cr.Namespace, dbID)

		instance := &agillv1alpha1.DBInstance{}
		err := r.client.Get(context.TODO(), request.NamespacedName, instance)
		if err != nil {
			return nil, err
		}
		instance.Status.DeployedInitially = true
		_, instance.Status.RDSInstanceStatus = r.dbInstanceExists(dbID)
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return nil, err
		}
	}
	return createOut, err
}
