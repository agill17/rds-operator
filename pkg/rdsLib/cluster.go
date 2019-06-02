package rdsLib

import (
	"errors"

	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

type cluster struct {
	rdsClient            *rds.RDS
	createInput          *rds.CreateDBClusterInput
	deleteInput          *rds.DeleteDBClusterInput
	restoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput
	object               runtime.Object
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
	if exists, _ := dh.clusterExists(); !exists {
		if _, err := dh.rdsClient.CreateDBCluster(dh.createInput); err != nil {
			logrus.Errorf("Failed to create new DB Cluster, %v", err)
			return err
		}
	}

	return nil
}

// Delete Cluster
func (dh *cluster) Delete() error {

	if exists, _ := dh.clusterExists(); exists {
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
	if exists, _ := dh.clusterExists(); !exists {

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

// GetAWSStatus gets cluster status
func (dh *cluster) GetAWSStatus() RDS_RESOURCE_STATE {
	_, state := dh.clusterExists()
	return state
}

// return bool ( exist / not exist ) and a remote status of the resource
func (dh *cluster) clusterExists() (bool, RDS_RESOURCE_STATE) {
	var clID string
	if dh.createInput != nil {
		clID = *dh.createInput.DBClusterIdentifier
	} else if dh.restoreFromSnapInput != nil {
		clID = *dh.restoreFromSnapInput.DBClusterIdentifier
	}

	exists, out := lib.DbClusterExists(
		&lib.RDSGenerics{
			RDSClient: dh.rdsClient,
			ClusterID: clID,
		},
	)

	return exists, parseRemoteStatus(*out.DBClusters[0].Status)
}
