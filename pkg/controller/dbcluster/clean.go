package dbcluster

import (
	"fmt"
	"time"

	// h "cloud.google.com/go/bigquery/benchmarks"
	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) deleteCluster(cr *agillv1alpha1.DBCluster, dbID string) error {
	var err error
	deleteClusterInput := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: &dbID,
		SkipFinalSnapshot:   &cr.Spec.DeletePolicy.SkipFinalSnapshot,
	}

	if *deleteClusterInput.SkipFinalSnapshot == false {
		currentTime := time.Now().Format("2006-01-02")
		snashotName := fmt.Sprintf("%v-%v", dbID, currentTime)
		deleteClusterInput.FinalDBSnapshotIdentifier = aws.String(snashotName)
	}

	if exists, _ := r.dbClusterExists(dbID); exists {
		if _, err = r.rdsClient.DeleteDBCluster(deleteClusterInput); err != nil {
			logrus.Errorf("ERROR while deleting cluster: %v", err)
			return err
		}
	}
	return err
}
