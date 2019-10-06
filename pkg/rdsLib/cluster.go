package rdsLib

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"os/exec"
	"strings"
	"time"

	"k8s.io/api/core/v1"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/utils"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

type cluster struct {
	k8sClient            client.Client
	rdsClient            *rds.RDS
	createInput          *rds.CreateDBClusterInput
	deleteInput          *rds.DeleteDBClusterInput
	restoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput
	runtimeObj           *kubev1alpha1.DBCluster
	clusterID            string
}

func NewCluster(rdsClient *rds.RDS, createInput *rds.CreateDBClusterInput,
	deleteInput *rds.DeleteDBClusterInput,
	restoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput,
	cr *kubev1alpha1.DBCluster, client client.Client, clusterID string) RDS {

	return &cluster{
		rdsClient:            rdsClient,
		createInput:          createInput,
		restoreFromSnapInput: restoreFromSnapInput,
		deleteInput:          deleteInput,
		k8sClient:            client,
		runtimeObj:           cr,
		clusterID:            clusterID,
	}
}

// Create Cluster
func (dh cluster) Create() error {

	exists, _, _ := utils.DbClusterExists(utils.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
	if !exists {

		/*
			Create could be 2 things
			1. Controller never created any db cluster so it must install a fresh new one

			2. Controller did create the db cluster atleast once and it no longer exists, so it must try to restore from a snapshot
				- for this, we try to find the latest snapshot ID available  and put that in CR, so it becomes a restore ActionType
				- If we dont find a snapshotID, we create a fresh one again.
		 */

		 // First lets check if controller ever created it and if it did and no longer exists in AWS, we switch to restore option
		 if dh.runtimeObj.Status.Created {
		 	latestSnapIDAvial, err := dh.findLatestAvailableSnapshot()
		 	if err != nil && latestSnapIDAvial != "" {
		 		return err
			}

		 	// Only if we found a snapshot avail, we go to restore
		 	if latestSnapIDAvial != "" {
				dh.runtimeObj.ClusterSpec.SnapshotIdentifier = &latestSnapIDAvial
				return utils.UpdateCr(dh.k8sClient, dh.runtimeObj)
			}
		 	// else we continue with a fresh install
		 }


		if err := dh.addCredsToClusterInput(); err != nil {
			return err
		}

		if _, err := dh.rdsClient.CreateDBCluster(dh.createInput); err != nil {
			logrus.Errorf("Failed to create new DB Cluster, %v", err)
			return err
		}
	}

	return nil
}

// Delete Cluster
func (dh cluster) Delete() error {

	exists, _, _ := utils.DbClusterExists(utils.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
	if exists {

		dh.setTimestampInSnapshotName()
		logrus.Infof(*dh.deleteInput.FinalDBSnapshotIdentifier)
		if _, err := dh.rdsClient.DeleteDBCluster(dh.deleteInput); err != nil {

			// if cluster is not found, return error, else move on from the delete call
			if err.(awserr.Error).Code() != rds.ErrCodeDBClusterNotFoundFault {
				logrus.Errorf("Failed to delete DB cluster: %v", err)
				return err
			}

		}
		logrus.Warnf("Successfully Deleted DB Cluster: %v", *dh.deleteInput.DBClusterIdentifier)
		return utils.RemoveFinalizer(dh.runtimeObj, dh.k8sClient, utils.DBClusterFinalizer)
	}
	return nil
}

// Restore Cluster from snapshot
func (dh cluster) Restore() error {

	exists, _, _ := utils.DbClusterExists(utils.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
	if !exists {

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

func (dh *cluster) addCredsToClusterInput() error {
	// ALWAYS grab credentials from a secret
	// a secret WILL exist whether its the user creates it or gets created by the controller
	ns := dh.runtimeObj.Namespace
	secretName := dh.runtimeObj.Status.SecretName
	userKey := dh.runtimeObj.Status.UsernameKey
	passKey := dh.runtimeObj.Status.PasswordKey
	exists, secret := utils.SecretExists(ns, secretName, dh.k8sClient)
	// incase it does not exist
	if !exists {
		return apiErrors.NewNotFound(v1.Resource("secret"), secret.Name)
	}

	//  or is getting deleted
	if secret.DeletionTimestamp != nil {
		return apiErrors.NewForbidden(v1.Resource("secret"), secretName, errors.New("K8sSecretGettingDeleted"))
	}

	if dh.createInput.MasterUsername == nil && dh.createInput.MasterUserPassword == nil {
		username := string(secret.Data[userKey])
		password := string(secret.Data[passKey])

		dh.createInput.MasterUsername = &username
		dh.createInput.MasterUserPassword = &password
		logrus.Infof("addCredsToClusterInput got invoked")
	}

	return nil
}

func (dh *cluster) setTimestampInSnapshotName() {
	if dh.deleteInput.FinalDBSnapshotIdentifier != nil && !*dh.deleteInput.SkipFinalSnapshot {
		currentTime := time.Now().Format("2006-01-02:03-02-44")
		snashotName := fmt.Sprintf("%v-%v", *dh.deleteInput.DBClusterIdentifier, strings.Replace(currentTime, ":", "-", -1))
		dh.deleteInput.FinalDBSnapshotIdentifier = &snashotName
	}
}

func (dh *cluster) SyncAwsStatusWithCRStatus() (string, error) {

	exists, out, _ := utils.DbClusterExists(utils.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
	currentLocalPhase := dh.runtimeObj.Status.CurrentPhase

	if exists {
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in AWS: %v", dh.runtimeObj.Name, dh.runtimeObj.Namespace, *out.DBClusters[0].Status)
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in CR: %v", dh.runtimeObj.Name, dh.runtimeObj.Namespace, currentLocalPhase)

		if currentLocalPhase != strings.ToLower(*out.DBClusters[0].Status) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", dh.runtimeObj.Namespace)
			dh.runtimeObj.Status.CurrentPhase = strings.ToLower(*out.DBClusters[0].Status)
			if err := utils.UpdateCrStatus(dh.k8sClient, dh.runtimeObj); err != nil {
				return "", err
			}
		}
	}
	return dh.runtimeObj.Status.CurrentPhase, nil

}

func (dh *cluster) findLatestAvailableSnapshot() (string, error) {
	cmd := fmt.Sprintf("aws rds describe-db-cluster-snapshots  --query \"DBClusterSnapshots[?DBClusterIdentifier=='%v']\" --region %v | jq -r 'max_by(.SnapshotCreateTime).DBClusterSnapshotIdentifier'", *dh.createInput.DBClusterIdentifier, "us-east-1")
	snapID, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		logrus.Errorf("Failed to execute command to get latest available snapshot for %v: %s", *dh.createInput.DBClusterIdentifier, err)
		return "", err
	}

	logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Latest snapshot id available: %v", dh.runtimeObj.Namespace, *dh.createInput.DBClusterIdentifier, strings.TrimSpace(string(snapID)))

	return strings.TrimSpace(string(snapID)), err
}
