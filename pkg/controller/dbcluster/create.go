package dbcluster

import (
	"context"

	"github.com/agill17/rds-operator/pkg/lib/dbHelpers"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	kubeapierror "k8s.io/apimachinery/pkg/api/errors"
)

func (r *ReconcileDBCluster) createItAndUpdateState(cr *kubev1alpha1.DBCluster, cluster *dbHelpers.Cluster) error {
	var err error

	err = dbHelpers.InstallRestoreDelete(cluster, dbHelpers.CREATE)
	if err != nil {
		return err
	}

	// check aws state and return error if not ready/available yet in AWS
	if err := r.handlePhases(cr); err != nil {
		return err
	}

	// once cr phase is available, change the created to true and update status
	cr.Status.Created = true
	_, cr.Status.DescriberClusterOutput = lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.CreateClusterSpec.DBClusterIdentifier})
	if err := lib.UpdateCrStatus(r.client, cr); err != nil {
		return err
	}

	spew.Dump(cr.Status)

	return nil

}

func (r *ReconcileDBCluster) createSecret(cr *kubev1alpha1.DBCluster) error {
	secretObj := r.getSecretObj(cr)
	if err := r.client.Create(context.TODO(), secretObj); err != nil && !kubeapierror.IsAlreadyExists(err) {
		logrus.Errorf("Error while creating secret object: %v", err)
		return err
	} else if kubeapierror.IsAlreadyExists(err) && cr.Status.SecretUpdateNeeded {
		logrus.Warnf("Updating cluster secret in namespace: %v", cr.Namespace)
		r.client.Update(context.TODO(), secretObj)
		cr.Status.SecretUpdateNeeded = false
		return lib.UpdateCrStatus(r.client, cr)
	}

	return nil
}
