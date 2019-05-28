package rdsLib

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

type Cluster struct {
	RDSClient            *rds.RDS
	CreateInput          *rds.CreateDBClusterInput
	DeleteInput          *rds.DeleteDBClusterInput
	RestoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput
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
			ClusterID: *dh.DeleteInput.DBClusterIdentifier,
		},
	)
	if exists {
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
			ClusterID: *dh.RestoreFromSnapInput.DBClusterIdentifier,
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
