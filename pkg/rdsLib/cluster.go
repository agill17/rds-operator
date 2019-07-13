package rdsLib

import (
	"errors"
	"fmt"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
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
func (dh *cluster) Create() error {

	exists, _ := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
	if !exists {

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

	exists, _ := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
	if exists {
		if _, err := dh.rdsClient.DeleteDBCluster(dh.deleteInput); err != nil {
			logrus.Errorf("Failed to delete DB cluster: %v", err)
			return err
		}
		logrus.Warnf("Successfully Deleted DB Cluster: %v", *dh.deleteInput.DBClusterIdentifier)
		dh.runtimeObj.SetFinalizers([]string{})
		return lib.UpdateCr(dh.k8sClient, dh.runtimeObj)
	}
	return nil
}

// Restore Cluster
func (dh *cluster) Restore() error {

	exists, _ := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
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
	exists, secret := lib.SecretExists(ns, secretName, dh.k8sClient)
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
		snashotName := fmt.Sprintf("%v-%v", dh.deleteInput.DBClusterIdentifier, strings.Replace(currentTime, ":", "-", -1))
		dh.deleteInput.FinalDBSnapshotIdentifier = &snashotName
	}
}

func (dh *cluster) GetAWSStatus() (string, error) {

	exists, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: dh.rdsClient, ClusterID: dh.clusterID})
	currentLocalPhase := dh.runtimeObj.Status.CurrentPhase

	if exists {
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in AWS: %v", dh.runtimeObj.Name, dh.runtimeObj.Namespace, *out.DBClusters[0].Status)
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in CR: %v", dh.runtimeObj.Name, dh.runtimeObj.Namespace, currentLocalPhase)

		if currentLocalPhase != strings.ToLower(*out.DBClusters[0].Status) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", dh.runtimeObj.Namespace)
			dh.runtimeObj.Status.CurrentPhase = strings.ToLower(*out.DBClusters[0].Status)
			if err := lib.UpdateCrStatus(dh.k8sClient, dh.runtimeObj); err != nil {
				return "", err
			}
		}
	}
	return dh.runtimeObj.Status.CurrentPhase, nil

}
