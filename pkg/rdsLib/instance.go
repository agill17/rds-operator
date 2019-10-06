package rdsLib

import (
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/utils"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

type instance struct {
	rdsClient         *rds.RDS
	createIn          *rds.CreateDBInstanceInput
	deleteIn          *rds.DeleteDBInstanceInput
	restoreFromSnapIn *rds.RestoreDBInstanceFromDBSnapshotInput
	instanceID        string
	runtimeObj        *kubev1alpha1.DBInstance
	k8sClient         client.Client
}

func NewInstance(rdsClient *rds.RDS,
	createIn *rds.CreateDBInstanceInput, deleteIn *rds.DeleteDBInstanceInput,
	restoreIn *rds.RestoreDBInstanceFromDBSnapshotInput, cr *kubev1alpha1.DBInstance, client client.Client, insID string) RDS {

	return &instance{
		rdsClient:         rdsClient,
		createIn:          createIn,
		deleteIn:          deleteIn,
		restoreFromSnapIn: restoreIn,
		instanceID:        insID,
		runtimeObj:        cr,
		k8sClient:         client,
	}
}

// Create Instance
func (i *instance) Create() error {
	exists, _ := utils.DBInstanceExists(utils.RDSGenerics{RDSClient: i.rdsClient, InstanceID: i.instanceID})
	if !exists {

		// if instance is part of cluster, check cluster existence before
		partOfCluster := i.createIn.DBClusterIdentifier != nil
		if partOfCluster {
			clusterExists, _ , _:= utils.DbClusterExists(utils.RDSGenerics{RDSClient: i.rdsClient, ClusterID: *i.createIn.DBClusterIdentifier})
			if !clusterExists {
				return utils.ErrorResourceCreatingInProgress{Message: "ClusterForDBInstanceNotFoundError"}
			}
		}

		if _, err := i.rdsClient.CreateDBInstance(i.createIn); err != nil {
			logrus.Errorf("Failed to create new DB Instance: %v", err)
			return err
		}
	}
	return nil
}

// Delete Instance
func (i *instance) Delete() error {
	exists, _ := utils.DBInstanceExists(utils.RDSGenerics{RDSClient: i.rdsClient, InstanceID: i.instanceID})
	if exists {
		if _, err := i.rdsClient.DeleteDBInstance(i.deleteIn); err != nil {
			logrus.Errorf("Failed to delete DB Instance: %v", err)
			return err
		}
		return utils.RemoveFinalizer(i.runtimeObj, i.k8sClient, utils.DBInstanceFinalizer)
	}
	return nil
}

// Restore Instance from a snapshot
func (i *instance) Restore() error {
	exists, _ := utils.DBInstanceExists(utils.RDSGenerics{RDSClient: i.rdsClient, InstanceID: i.instanceID})
	if !exists && !strings.Contains(*i.restoreFromSnapIn.Engine, "aurora-") && i.restoreFromSnapIn.DBSnapshotIdentifier != nil {
		if _, err := i.rdsClient.RestoreDBInstanceFromDBSnapshot(i.restoreFromSnapIn); err != nil {
			logrus.Errorf("Failed to restore DB cluster from snapshot :%v", err)
			return err
		}
	}
	return nil
}


// SyncAwsStatusWithCRStatus returns resource state in aws
func (i *instance) SyncAwsStatusWithCRStatus() (string, error) {

	exists, out := utils.DBInstanceExists(utils.RDSGenerics{RDSClient: i.rdsClient, InstanceID: i.instanceID})
	currentLocalPhase := i.runtimeObj.Status.CurrentPhase

	if exists {
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in AWS: %v", i.runtimeObj.Name, i.runtimeObj.Namespace, *out.DBInstances[0].DBInstanceStatus)
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in CR: %v", i.runtimeObj.Name, i.runtimeObj.Namespace, currentLocalPhase)

		if currentLocalPhase != strings.ToLower(*out.DBInstances[0].DBInstanceStatus) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", i.runtimeObj.Namespace)
			i.runtimeObj.Status.CurrentPhase = strings.ToLower(*out.DBInstances[0].DBInstanceStatus)
			if err := utils.UpdateCrStatus(i.k8sClient, i.runtimeObj); err != nil {
				return "", err
			}
		}
	}
	return i.runtimeObj.Status.CurrentPhase, nil

}

func (i *instance) setTimestampInSnapshotName() {
	if i.deleteIn.FinalDBSnapshotIdentifier != nil && !*i.deleteIn.SkipFinalSnapshot {
		currentTime := time.Now().Format("2006-01-02:03-02-44")
		snashotName := fmt.Sprintf("%v-%v", i.deleteIn.DBInstanceIdentifier, strings.Replace(currentTime, ":", "-", -1))
		i.deleteIn.FinalDBSnapshotIdentifier = &snashotName
	}
}
