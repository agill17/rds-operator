package dbcluster

import (
	"fmt"
	"strings"
	"time"

	"github.com/agill17/rds-operator/pkg/lib"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func (r *ReconcileDBCluster) deleteCluster(cr *kubev1alpha1.DBCluster) error {
	var err error
	deleteClusterInput := cr.DeleteSpec

	if !(*deleteClusterInput.SkipFinalSnapshot) {
		currentTime := time.Now().Format("2006-01-02:03-02-44")
		snashotName := fmt.Sprintf("%v-%v", *cr.Spec.DBClusterIdentifier, strings.Replace(currentTime, ":", "-", -1))
		if cr.DeleteSpec.FinalDBSnapshotIdentifier == nil {
			deleteClusterInput.FinalDBSnapshotIdentifier = aws.String(snashotName)
		}
	}

	if exists, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.DBClusterIdentifier}); exists {

		// already in deleting state?
		if *out.DBClusters[0].Status == "deleting" {
			logrus.Warnf("DBCluster is already in deleting state: %v: %v", deleteClusterInput.DBClusterIdentifier, *out.DBClusters[0].Status)
			return nil
		}

		if _, err = r.rdsClient.DeleteDBCluster(deleteClusterInput); err != nil {
			logrus.Errorf("ERROR while deleting cluster: %v", err)
			return err
		}
	}
	return err
}

func (r *ReconcileDBCluster) handleDelete(cr *kubev1alpha1.DBCluster) error {
	deletionTimeExists := cr.GetDeletionTimestamp() != nil
	zeroFinalizers := len(cr.GetFinalizers()) == 0

	if deletionTimeExists && !zeroFinalizers {
		if err := r.deleteCluster(cr); err != nil {
			return err
		}

		logrus.Infof("Successfully deleted DBCluster from namespace: %v", cr.Namespace)

		cr.SetFinalizers([]string{})
		if err := lib.UpdateCr(r.client, cr); err != nil {
			return err
		}
	}
	return nil
}
