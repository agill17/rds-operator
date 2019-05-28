package rdsLib

import (
	"errors"

	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

type cluster struct {
	rdsClient            *rds.RDS
	createInput          *rds.CreateDBClusterInput
	deleteInput          *rds.DeleteDBClusterInput
	restoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput
}

func NewCluster(rdsClient *rds.RDS, createInput *rds.CreateDBClusterInput,
	deleteInput *rds.DeleteDBClusterInput, restoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput) RDS {
	return &cluster{
		rdsClient:            rdsClient,
		createInput:          createInput,
		restoreFromSnapInput: restoreFromSnapInput,
		deleteInput:          deleteInput,
	}
}

// Create Cluster
func (dh *cluster) Create() error {
	if !dh.clusterExists() {
		if _, err := dh.rdsClient.CreateDBCluster(dh.createInput); err != nil {
			logrus.Errorf("Failed to create new DB Cluster, %v", err)
			return err
		}
	}

	return nil
}

// Delete Cluster
func (dh *cluster) Delete() error {

	if dh.clusterExists() {
		if _, err := dh.rdsClient.DeleteDBCluster(dh.deleteInput); err != nil {
			logrus.Errorf("Failed to delete DB cluster: %v", err)
			return err
		}
		logrus.Warnf("Successfully Deleted DB Cluster: %v", *dh.deleteInput.DBClusterIdentifier)
	}
	return nil
}

// Restore Cluster
func (dh *cluster) Restore() error {
	if !dh.clusterExists() {

		if dh.restoreFromSnapInput.DBClusterIdentifier == nil ||
			dh.restoreFromSnapInput.SnapshotIdentifier == nil {
			logrus.Errorf("Restore DBClusterIdentifier and SnapshotIdentifier cannot be empty")
			return errors.New("RestoreDBClusterInsufficientParameterError")
		}

		if _, err := dh.rdsClient.RestoreDBClusterFromSnapshot(dh.restoreFromSnapInput); err != nil {
			logrus.Errorf("Failed to restore DB cluster from snapshot :%v", err)
			return err
		}
	}
	return nil
}

func (dh *cluster) clusterExists() bool {
	var clID string
	if dh.createInput != nil {
		clID = *dh.createInput.DBClusterIdentifier
	} else if dh.restoreFromSnapInput != nil {
		clID = *dh.restoreFromSnapInput.DBClusterIdentifier
	}

	exists, _ := lib.DbClusterExists(
		&lib.RDSGenerics{
			RDSClient: dh.rdsClient,
			ClusterID: clID,
		},
	)
	return exists
}
