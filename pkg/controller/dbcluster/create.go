package dbcluster

import (
	"context"
	"strings"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	kubeapierror "k8s.io/apimachinery/pkg/api/errors"
)

func (r *ReconcileDBCluster) createItAndUpdateState(cr *kubev1alpha1.DBCluster) error {

	// initial create
	username, password := cr.Status.Username, cr.Status.Password
	_, err := r.createCluster(cr, username, password)
	if err != nil {
		logrus.Errorf("Something went wrong while creating the db cluster: %v", err)
		return err
	}

	if err := r.handlePhases(cr); err != nil {
		return err
	}

	// once cr phase is available, change the created to true and update status
	cr.Status.Created = true
	_, cr.Status.DescriberClusterOutput = lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.DBClusterIdentifier})
	if err := r.updateCrStats(cr); err != nil {
		return err
	}

	spew.Dump(cr.Status)
	return nil

}

func (r *ReconcileDBCluster) createCluster(cr *kubev1alpha1.DBCluster, username, password string) (*rds.DescribeDBClustersOutput, error) {
	var err error
	if exists, _ := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.DBClusterIdentifier}); !exists {
		logrus.Infof("Creating db cluster first")
		_, err = r.rdsClient.CreateDBCluster(cr.Spec)
		if err != nil {
			logrus.Errorf("ERROR while creating DB Cluster%v:", err)
			spew.Dump(cr)
			return nil, err
		}
	}

	_, dbClusterOutput := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.DBClusterIdentifier})

	cr.Status.CurrentPhase = strings.ToLower(*dbClusterOutput.DBClusters[0].Status)
	cr.Status.Username = username
	cr.Status.Password = password
	if err := r.updateCrStats(cr); err != nil {
		return nil, err
	}

	return dbClusterOutput, err
}

func (r *ReconcileDBCluster) createSecret(cr *kubev1alpha1.DBCluster, installType string) error {
	if installType == "newInstall" {
		secretObj := r.getSecretObj(cr)
		if err := r.client.Create(context.TODO(), secretObj); err != nil && !kubeapierror.IsAlreadyExists(err) {
			logrus.Errorf("Error while creating secret object: %v", err)
			return err
		} else if kubeapierror.IsAlreadyExists(err) && cr.Status.SecretUpdateNeeded {
			logrus.Warnf("Updating cluster secret in namespace: %v", cr.Namespace)
			r.client.Update(context.TODO(), secretObj)
			cr.Status.SecretUpdateNeeded = false
			return r.updateCrStats(cr)
		}
	} else {
		logrus.Warnf("Not creating any secret as this deployment was creation from from existing snapshot in ns: %v", cr.Namespace)
	}
	return nil
}
