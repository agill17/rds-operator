package rdsLib

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
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
}

func NewCluster(rdsClient *rds.RDS, createInput *rds.CreateDBClusterInput,
	deleteInput *rds.DeleteDBClusterInput,
	restoreFromSnapInput *rds.RestoreDBClusterFromSnapshotInput,
	cr *kubev1alpha1.DBCluster, client client.Client) RDS {

	return &cluster{
		rdsClient:            rdsClient,
		createInput:          createInput,
		restoreFromSnapInput: restoreFromSnapInput,
		deleteInput:          deleteInput,
		k8sClient:            client,
		runtimeObj:           cr,
	}
}

// Create Cluster
func (dh *cluster) Create() error {
	if exists := dh.clusterExists(); !exists {

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
func (dh *cluster) Delete() error {

	if exists := dh.clusterExists(); exists {
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
	if exists := dh.clusterExists(); !exists {

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

// return bool ( exist / not exist ) and a remote status of the resource
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

func (dh *cluster) addCredsToClusterInput() error {
	// ALWAYS grab credentials from a secret
	// a secret WILL exist whether its the user creates it or gets created by the controller
	ns := dh.runtimeObj.Namespace
	secretName := dh.runtimeObj.Status.SecretName
	userKey := dh.runtimeObj.Status.UsernameKey
	passKey := dh.runtimeObj.Status.PasswordKey
	exists, secret := lib.SecretExists(ns, secretName, dh.k8sClient)
	// incase it does not exist
	if !exists {
		return &lib.ErrorKubernetesSecretDoesNotExist{Message: "K8S secret does not exist: " + secretName}
	}

	//  or is getting deleted
	if secret.DeletionTimestamp != nil {
		return &lib.ErrorKubernetesSecretGettingDeleted{
			Message: "K8S secret is getting deleted: " + secretName,
		}
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
		snashotName := fmt.Sprintf("%v-%v", dh.deleteInput.DBClusterIdentifier, strings.Replace(currentTime, ":", "-", -1))
		dh.deleteInput.FinalDBSnapshotIdentifier = &snashotName
	}
}
