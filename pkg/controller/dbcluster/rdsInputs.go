package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/service/rds"
)

func GetRestoreClusterDBFromSnapInput(cr *kubev1alpha1.DBCluster, clusterDBId string, snapID string) *rds.RestoreDBClusterFromSnapshotInput {

	restoreClusterIn := &rds.RestoreDBClusterFromSnapshotInput{
		AvailabilityZones:   cr.Spec.AvailabilityZones,
		DBClusterIdentifier: &clusterDBId,
		DBSubnetGroupName:   cr.Spec.DBSubnetGroupName,
		DatabaseName:        cr.Spec.DatabaseName,
		DeletionProtection:  cr.Spec.DeletionProtection,
		Engine:              cr.Spec.Engine,
		EngineMode:          cr.Spec.EngineMode,
		EngineVersion:       cr.Spec.EngineVersion,
		SnapshotIdentifier:  &snapID,
		VpcSecurityGroupIds: cr.Spec.VpcSecurityGroupIds,
	}
	return restoreClusterIn
}
