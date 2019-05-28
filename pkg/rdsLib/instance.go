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
	restoreIn *rds.RestoreDBInstanceFromDBSnapshotInput) (RDS, error) {

	return &instance{
		rdsClient:         rdsClient,
		createIn:          createIn,
		deleteIn:          deleteIn,
		restoreFromSnapIn: restoreIn,
	}, nil
}

// Create Instance
func (i *instance) Create() error {
	exists, _ := lib.DBInstanceExists(
		&lib.RDSGenerics{
			RDSClient:  i.rdsClient,
			InstanceID: *i.createIn.DBInstanceIdentifier,
		},
	)
	if !exists {
		if _, err := i.rdsClient.CreateDBInstance(i.createIn); err != nil {
			logrus.Errorf("Failed to create new DB Instance: %v", err)
			return err
		}
	}
	return nil
}

// Delete Instance
func (i *instance) Delete() error {
	exists, _ := lib.DBInstanceExists(
		&lib.RDSGenerics{
			RDSClient:  i.rdsClient,
			InstanceID: *i.deleteIn.DBInstanceIdentifier,
		},
	)
	if exists {
		if _, err := i.rdsClient.DeleteDBInstance(i.deleteIn); err != nil {
			logrus.Errorf("Failed to delete DB Instance: %v", err)
			return err
		}
	}
	return nil
}

// Restore Instance
func (i *instance) Restore() error {
	exists, _ := lib.DBInstanceExists(
		&lib.RDSGenerics{
			RDSClient:  i.rdsClient,
			InstanceID: *i.restoreFromSnapIn.DBInstanceIdentifier,
		},
	)
	if !exists {
		if _, err := i.rdsClient.RestoreDBInstanceFromDBSnapshot(i.restoreFromSnapIn); err != nil {
			logrus.Errorf("Failed to restore DB cluster from snapshot :%v", err)
			return err
		}
	}
	return nil
}
