package dbinstance

import (
	"strings"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

func (r *ReconcileDBInstance) createDBInstanceInput(cr *agillv1alpha1.DBInstance, dbID string) *rds.CreateDBInstanceInput {
	instanceInput := &rds.CreateDBInstanceInput{
		AvailabilityZone:        &cr.Spec.Az,
		DBInstanceIdentifier:    &dbID,
		AutoMinorVersionUpgrade: &cr.Spec.AutoMinorVersionUpgrade,

		CopyTagsToSnapshot:  aws.Bool(true),
		DBClusterIdentifier: &cr.Spec.DBClusterIdentifier,
		DBInstanceClass:     &cr.Spec.DBInstanceClass,
		DBSecurityGroups:    aws.StringSlice(cr.Spec.DBSecurityGroups),
		DBSubnetGroupName:   &cr.Spec.DBSubnetGroupName,
		Engine:              &cr.Spec.Engine,
		EngineVersion:       &cr.Spec.EngineVersion,
		PubliclyAccessible:  &cr.Spec.PubliclyAccessible,
	}

	if strings.ToLower(cr.Spec.Engine) != "aurora-mysql" {
		instanceInput.DBName = &cr.Spec.DBName
		instanceInput.DeletionProtection = &cr.Spec.DeletionProtection
		instanceInput.AllocatedStorage = &cr.Spec.AllocatedStorage
		instanceInput.BackupRetentionPeriod = &cr.Spec.BackupRetentionPeriod
		instanceInput.MasterUsername = &cr.Spec.MasterUsername
		instanceInput.MasterUserPassword = &cr.Spec.MasterPassword
		instanceInput.StorageEncrypted = &cr.Spec.StorageEncrypted
		instanceInput.VpcSecurityGroupIds = aws.StringSlice(cr.Spec.VpcSecurityGroupIds)
	}

	if cr.Spec.DBParameterGroupName != "" {
		instanceInput.DBParameterGroupName = &cr.Spec.DBParameterGroupName
	}

	if cr.Spec.DBClusterIdentifier != "" {
		instanceInput.DBClusterIdentifier = &cr.Spec.DBClusterIdentifier
	}

	return instanceInput
}

func (r *ReconcileDBInstance) restoreFromSnapInput(cr *agillv1alpha1.DBInstance, dbID string, snapID string) *rds.RestoreDBInstanceFromDBSnapshotInput {
	restoreDBInput := &rds.RestoreDBInstanceFromDBSnapshotInput{
		AutoMinorVersionUpgrade: aws.Bool(cr.Spec.AutoMinorVersionUpgrade),
		AvailabilityZone:        aws.String(cr.Spec.Az),
		CopyTagsToSnapshot:      aws.Bool(true),
		DBInstanceClass:         aws.String(cr.Spec.DBInstanceClass),
		DBInstanceIdentifier:    aws.String(dbID),
		DBSubnetGroupName:       aws.String(cr.Spec.DBSubnetGroupName),
		DeletionProtection:      aws.Bool(cr.Spec.DeletionProtection),
		Engine:                  aws.String(cr.Spec.Engine),
		DBSnapshotIdentifier:    aws.String(snapID),
	}
	return restoreDBInput
}
