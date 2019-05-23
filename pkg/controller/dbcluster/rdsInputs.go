package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/service/rds"
)

func GetRestoreClusterDBFromSnapInput(cr *kubev1alpha1.DBCluster, clusterDBId string, snapID string) *rds.RestoreDBClusterFromSnapshotInput {

	restoreClusterIn := &rds.RestoreDBClusterFromSnapshotInput{
		AvailabilityZones:   cr.Spec.CreateClusterSpec.AvailabilityZones,
		DBClusterIdentifier: &clusterDBId,
		DBSubnetGroupName:   cr.Spec.CreateClusterSpec.DBSubnetGroupName,
		DatabaseName:        cr.Spec.CreateClusterSpec.DatabaseName,
		DeletionProtection:  cr.Spec.CreateClusterSpec.DeletionProtection,
		Engine:              cr.Spec.CreateClusterSpec.Engine,
		EngineMode:          cr.Spec.CreateClusterSpec.EngineMode,
		EngineVersion:       cr.Spec.CreateClusterSpec.EngineVersion,
		SnapshotIdentifier:  &snapID,
		VpcSecurityGroupIds: cr.Spec.CreateClusterSpec.VpcSecurityGroupIds,
	}
	return restoreClusterIn
}
