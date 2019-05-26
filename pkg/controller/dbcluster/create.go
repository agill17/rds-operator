package dbcluster

import (
	"github.com/agill17/rds-operator/pkg/lib/dbHelpers"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/davecgh/go-spew/spew"
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
