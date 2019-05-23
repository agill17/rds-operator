package dbcluster

import (
	"context"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	kubeapierror "k8s.io/apimachinery/pkg/api/errors"
)

func (r *ReconcileDBCluster) createItAndUpdateState(cr *kubev1alpha1.DBCluster, installType lib.ClusterInstallType) error {
	var err error
	var clusterID string

	switch installType {

	case lib.CLUSTER_INSTALL_NEW:
		clusterID = *cr.Spec.DBClusterIdentifier
		if exists, _ := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: clusterID}); !exists {
			err = r.createCluster(cr)
		}

	case lib.CLUSTER_INSTALL_FROM_SNAPSHOT:
		clusterID = *cr.CreateFromSnapshot.DBClusterIdentifier
		if exists, _ := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: clusterID}); !exists {
			err = r.createFromSnapshot(cr)
		}

	}

	if err != nil {
		logrus.Errorf("Something went wrong while creating the db cluster: %v", err)
		return err
	}

	// check aws state and return error if not ready/available yet in AWS
	if err := r.handlePhases(cr, clusterID); err != nil {
		return err
	}

	// once cr phase is available, change the created to true and update status
	cr.Status.Created = true
	_, cr.Status.DescriberClusterOutput = lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.DBClusterIdentifier})
	if err := lib.UpdateCrStatus(r.client, cr); err != nil {
		return err
	}

	spew.Dump(cr.Status)

	return nil

}

func (r *ReconcileDBCluster) createFromSnapshot(cr *kubev1alpha1.DBCluster) error {

	if _, err := r.rdsClient.RestoreDBClusterFromSnapshot(cr.CreateFromSnapshot); err != nil {
		logrus.Errorf("Failed to create a new db cluster from snapshot ~~~> %v", err)
		return err
	}

	return nil
}

func (r *ReconcileDBCluster) createCluster(cr *kubev1alpha1.DBCluster) error {
	var err error

	if _, err = r.rdsClient.CreateDBCluster(cr.Spec); err != nil {
		logrus.Errorf("ERROR while creating DB Cluster %v:", err)
		spew.Dump(cr)
		return err
	}

	return nil
}

func (r *ReconcileDBCluster) createSecret(cr *kubev1alpha1.DBCluster, installType lib.ClusterInstallType) error {
	if installType == lib.CLUSTER_INSTALL_NEW {
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
