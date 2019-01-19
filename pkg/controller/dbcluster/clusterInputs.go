package dbcluster

import (
	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

// currently using this for aurora engines only
func getCreateDBClusterInput(cr *agillv1alpha1.DBCluster, dbID string) *rds.CreateDBClusterInput {

	clusterInput := &rds.CreateDBClusterInput{
		AvailabilityZones:     aws.StringSlice(cr.Spec.Azs),
		BackupRetentionPeriod: &cr.Spec.BackupRetentionPeriod,
		DBClusterIdentifier:   &dbID,
		DBSubnetGroupName:     &cr.Spec.DBSubnetGroupName,
		DatabaseName:          &cr.Spec.DatabaseName,
		DeletionProtection:    &cr.Spec.DeletionProtection,
		Engine:                &cr.Spec.Engine,
		EngineMode:            &cr.Spec.EngineMode,
		EngineVersion:         &cr.Spec.EngineVersion,
		MasterUsername:        &cr.Spec.MasterUsername,
		MasterUserPassword:    &cr.Spec.MasterPassword,
		StorageEncrypted:      &cr.Spec.StorageEncrypted,
		VpcSecurityGroupIds:   aws.StringSlice(cr.Spec.VpcSecurityGroupIds),
	}
	if cr.Spec.DBClusterParameterGroupName != "" {
		clusterInput.DBClusterParameterGroupName = &cr.Spec.DBClusterParameterGroupName
	}
	return clusterInput
}

func GetRestoreClusterDBFromSnapInput(clusterCr *agillv1alpha1.DBCluster, clusterDBId string, snapID string) *rds.RestoreDBClusterFromSnapshotInput {

	restoreClusterIn := &rds.RestoreDBClusterFromSnapshotInput{
		AvailabilityZones:   aws.StringSlice(clusterCr.Spec.Azs),
		DBClusterIdentifier: &clusterDBId,
		DBSubnetGroupName:   &clusterCr.Spec.DBSubnetGroupName,
		DatabaseName:        &clusterCr.Spec.DatabaseName,
		DeletionProtection:  &clusterCr.Spec.DeletionProtection,
		Engine:              &clusterCr.Spec.Engine,
		EngineMode:          &clusterCr.Spec.EngineMode,
		EngineVersion:       &clusterCr.Spec.EngineVersion,
		SnapshotIdentifier:  &snapID,
		VpcSecurityGroupIds: aws.StringSlice(clusterCr.Spec.VpcSecurityGroupIds),
	}
	return restoreClusterIn
}
