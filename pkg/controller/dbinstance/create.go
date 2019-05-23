package dbinstance

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/davecgh/go-spew/spew"
)

func (r *ReconcileDBInstance) createNewDBInstance(cr *kubev1alpha1.DBInstance) (*rds.CreateDBInstanceOutput, error) {

	var createOut *rds.CreateDBInstanceOutput
	var err error

	if err := r.waitForClusterIfNeeded(cr); err != nil {
		return nil, err
	}

	// only create if it does not exist or not already being created
	if exists, _ := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: r.rdsClient, InstanceID: *cr.Spec.DBInstanceIdentifier}); !exists {
		logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Initial -- Creating DB", cr.Namespace, *cr.Spec.DBInstanceIdentifier)
		createOut, err = r.rdsClient.CreateDBInstance(cr.Spec)
		if err != nil {
			logrus.Errorf("Error while creating DB: %v", err)
			return createOut, err
		}
		spew.Dump(createOut)
	}

	if err := r.handlePhases(cr); err != nil {
		return nil, err
	}

	cr.Status.DeployedInitially = true
	_, rdsInstanceStatus := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: r.rdsClient, InstanceID: *cr.Spec.DBInstanceIdentifier})
	cr.Status.RDSInstanceStatus = rdsInstanceStatus

	// update status
	if err = r.updateResourceStatus(cr); err != nil {
		logrus.Errorf("Failed to update cr status for DBInstance: %v", err)
		return nil, err
	}
	logrus.Infof("CreateDBInstance was successful. Updated status")

	return createOut, err
}

func (r *ReconcileDBInstance) waitForClusterIfNeeded(cr *kubev1alpha1.DBInstance) error {
	var err error
	// when cluster is still not available, this will throw ErrorClusterCreatingInProgress
	// only run this when this DBInstance is part of a DBCluster
	if cr.Spec.DBClusterIdentifier != nil && !cr.Status.DBClusterMarkedAvail {
		logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Part of cluster: %v -- checking if its available first", cr.Namespace, *cr.Spec.DBClusterIdentifier, *cr.Spec.DBClusterIdentifier)
		err = r.dbClusterReady(cr)
		if err != nil {
			return err
		}
		cr.Status.DBClusterMarkedAvail = true
		if err := r.updateResourceStatus(cr); err != nil {
			return err
		}
	}
	return err
}
