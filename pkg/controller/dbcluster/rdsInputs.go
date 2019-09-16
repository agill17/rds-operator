package dbcluster

import (
	"fmt"
	"github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jinzhu/copier"
)

func deleteClusterInput(cr *v1alpha1.DBCluster) *rds.DeleteDBClusterInput {
	finalSnapId := fmt.Sprintf("%v-%v", cr.Name, cr.Namespace)
	return &rds.DeleteDBClusterInput{
		DBClusterIdentifier: cr.ClusterSpec.DBClusterIdentifier,
		FinalDBSnapshotIdentifier: aws.String(finalSnapId),
		SkipFinalSnapshot: aws.Bool(false),
	}
}

func restoreFromSnapshotInput(cr *v1alpha1.DBCluster) (*rds.RestoreDBClusterFromSnapshotInput, error) {
	var restoreIn rds.RestoreDBClusterFromSnapshotInput
	if err := copier.Copy(&restoreIn, cr.ClusterSpec); err != nil {
	return nil, err
	}
	return &restoreIn, nil
}


func createClusterInput(cr *v1alpha1.DBCluster) (*rds.CreateDBClusterInput, error) {
	var createIn rds.CreateDBClusterInput
	if err := copier.Copy(&createIn, cr.ClusterSpec); err != nil {
		return nil, err
	}
	return &createIn, nil
}