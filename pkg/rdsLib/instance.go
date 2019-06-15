package rdsLib

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

type instance struct {
	rdsClient         *rds.RDS
	createIn          *rds.CreateDBInstanceInput
	deleteIn          *rds.DeleteDBInstanceInput
	restoreFromSnapIn *rds.RestoreDBInstanceFromDBSnapshotInput
}

func NewInstance(rdsClient *rds.RDS,
	createIn *rds.CreateDBInstanceInput, deleteIn *rds.DeleteDBInstanceInput,
	restoreIn *rds.RestoreDBInstanceFromDBSnapshotInput) RDS {

	return &instance{
		rdsClient:         rdsClient,
		createIn:          createIn,
		deleteIn:          deleteIn,
		restoreFromSnapIn: restoreIn,
	}
}

// Create Instance
func (i *instance) Create() error {
	if exists := i.instanceExists(); !exists {
		if _, err := i.rdsClient.CreateDBInstance(i.createIn); err != nil {
			logrus.Errorf("Failed to create new DB Instance: %v", err)
			return err
		}
	}
	return nil
}

// Delete Instance
func (i *instance) Delete() error {

	if exists := i.instanceExists(); exists {
		if _, err := i.rdsClient.DeleteDBInstance(i.deleteIn); err != nil {
			logrus.Errorf("Failed to delete DB Instance: %v", err)
			return err
		}
	}
	return nil
}

// Restore Instance
func (i *instance) Restore() error {

	if exists := i.instanceExists(); !exists {
		if _, err := i.rdsClient.RestoreDBInstanceFromDBSnapshot(i.restoreFromSnapIn); err != nil {
			logrus.Errorf("Failed to restore DB cluster from snapshot :%v", err)
			return err
		}
	}
	return nil
}

func (i *instance) instanceExists() bool {
	var insID string
	if i.createIn != nil {
		insID = *i.createIn.DBInstanceIdentifier
	} else if i.restoreFromSnapIn != nil {
		insID = *i.restoreFromSnapIn.DBInstanceIdentifier
	}

	exists, _ := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: i.rdsClient, InstanceID: insID})

	return exists
}
