package dbHelpers

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

type DBInstallType string

const (
	CREATE  DBInstallType = "new"
	RESTORE DBInstallType = "fromSnapshot"
	DELETE  DBInstallType = "delete"
	UNKNOWN DBInstallType = "unknown"
)

type Cluster struct {
	RDSClient            *rds.RDS
	CreateInput          *rds.CreateDBClusterInput
	DeleteInput          *rds.DeleteDBClusterInput
	RestoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput
}

type Instance struct {
	RDSClient            *rds.RDS
	CreateInput          *rds.CreateDBInstanceInput
	DeleteInput          *rds.DeleteDBInstanceInput
	RestoreFromSnapInput *rds.RestoreDBInstanceFromDBSnapshotInput
}

type DB interface {
	Create() error
	Delete() error
	Restore() error
}

// Create Cluster
func (dh *Cluster) Create() error {
	exists, _ := lib.DbClusterExists(
		&lib.RDSGenerics{
			RDSClient: dh.RDSClient,
			ClusterID: *dh.CreateInput.DBClusterIdentifier,
		},
	)

	if !exists {
		if _, err := dh.RDSClient.CreateDBCluster(dh.CreateInput); err != nil {
			logrus.Errorf("Failed to create new DB Cluster, %v", err)
			return err
		}
	}

	return nil
}

// Delete Cluster
func (dh *Cluster) Delete() error {

	exists, _ := lib.DbClusterExists(
		&lib.RDSGenerics{
			RDSClient: dh.RDSClient,
			ClusterID: *dh.CreateInput.DBClusterIdentifier,
		},
	)
	if !exists {
		if _, err := dh.RDSClient.DeleteDBCluster(dh.DeleteInput); err != nil {
			logrus.Errorf("Failed to delete DB cluster: %v", err)
			return err
		}
	}
	return nil
}

// Restore Cluster
func (dh *Cluster) Restore() error {
	exists, _ := lib.DbClusterExists(
		&lib.RDSGenerics{
			RDSClient: dh.RDSClient,
			ClusterID: *dh.CreateInput.DBClusterIdentifier,
		},
	)
	if !exists {
		if _, err := dh.RDSClient.RestoreDBClusterFromSnapshot(dh.RestoreFromSnapInput); err != nil {
			logrus.Errorf("Failed to restore DB cluster from snapshot :%v", err)
			return err
		}
	}
	return nil
}

// // HandlePhases Cluster
// func (dh *Cluster) HandlePhases() error {
// 	exists, _ := lib.DbClusterExists(
// 		&lib.RDSGenerics{
// 			RDSClient: dh.RDSClient,
// 			ClusterID: *dh.CreateInput.DBClusterIdentifier,
// 		},
// 	)
// 	if exists {

// 	}
// 	return nil
// }

// Create Instance
func (dh *Instance) Create() error {
	exists, _ := lib.DBInstanceExists(
		&lib.RDSGenerics{
			RDSClient:  dh.RDSClient,
			InstanceID: *dh.CreateInput.DBInstanceIdentifier,
		},
	)
	if !exists {
		if _, err := dh.RDSClient.CreateDBInstance(dh.CreateInput); err != nil {
			logrus.Errorf("Failed to create new DB Instance: %v", err)
			return err
		}
	}
	return nil
}

// Delete Instance
func (dh *Instance) Delete() error {
	exists, _ := lib.DBInstanceExists(
		&lib.RDSGenerics{
			RDSClient:  dh.RDSClient,
			InstanceID: *dh.CreateInput.DBInstanceIdentifier,
		},
	)
	if !exists {
		if _, err := dh.RDSClient.DeleteDBInstance(dh.DeleteInput); err != nil {
			logrus.Errorf("Failed to delete DB Instance: %v", err)
			return err
		}
	}
	return nil
}

// Restore Instance
func (dh *Instance) Restore() error {
	exists, _ := lib.DBInstanceExists(
		&lib.RDSGenerics{
			RDSClient:  dh.RDSClient,
			InstanceID: *dh.CreateInput.DBInstanceIdentifier,
		},
	)
	if !exists {
		if _, err := dh.RDSClient.RestoreDBInstanceFromDBSnapshot(dh.RestoreFromSnapInput); err != nil {
			logrus.Errorf("Failed to restore DB cluster from snapshot :%v", err)
			return err
		}
	}
	return nil
}

func InstallRestoreDelete(dbInput DB, installType DBInstallType) error {
	switch installType {
	case CREATE:
		return dbInput.Create()
	case DELETE:
		return dbInput.Delete()
	case RESTORE:
		return dbInput.Restore()
	}

	return nil
}
