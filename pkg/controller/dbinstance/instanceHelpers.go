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
