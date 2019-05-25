package dbinstance

import (
	"fmt"
	"os/exec"
	"strings"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBInstance) restore(cr *kubev1alpha1.DBInstance) error {

	/*
		Two types of restore
		1. Restore DBInstance thats part of DBCluster: then perform plain new dbInstance deployment
		2. Restore DBInstance thats NOT part of DBCluster: then get the latest snapID and restore using that
	*/

	// is it a DBInstance on its own?
	// true if yes, false if associated to DBCluster
	isStandAloneInstance := isStandAlone(cr)

	if !isStandAloneInstance {
		cr.Status.DeployedInitially = false // let operator perform a brand new create
		if err := r.updateResourceStatus(cr); err != nil {
			return err
		}
	} else if isStandAloneInstance {
		// meaning instance on its own, then restore from latest available snapshot
		if err := r.restoreStandAloneInstance(cr); err != nil {
			return err
		}
	}

	return nil
}

func getInstanceSnapID(cr *kubev1alpha1.DBInstance) (string, error) {
	dbInsID := *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier
	cmd := fmt.Sprintf("aws rds describe-db-snapshots --query \"DBSnapshots[?DBInstanceIdentifier=='%v']\" --region %v | jq -r 'max_by(.SnapshotCreateTime).DBSnapshotIdentifier'", dbInsID, cr.Spec.Region)
	snapID, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		logrus.Errorf("Failed to execute command: %s", err)
		return "", err
	}

	logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Latest snapshot id available: %v", cr.Namespace, dbInsID, strings.TrimSpace(string(snapID)))

	return strings.TrimSpace(string(snapID)), err
}

// this will restore instances from snapshot that were not prt of a cluster
func (r *ReconcileDBInstance) restoreStandAloneInstance(cr *kubev1alpha1.DBInstance) error {
	dbInsID := *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier
	if exists, _ := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: r.rdsClient, InstanceID: dbInsID}); !exists {
		if _, err := r.rdsClient.RestoreDBInstanceFromDBSnapshot(restoreInstanceInput(cr)); err != nil {
			logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: ERROR While restoring db from snapshot: %v", cr.Namespace, dbInsID, err)
			return err
		}
	}
	if err := r.handlePhases(cr); err != nil {
		return err
	}

	if err := r.updateK8sFiles(cr); err != nil {
		return err
	}
	return nil
}

func restoreInstanceInput(cr *kubev1alpha1.DBInstance) *rds.RestoreDBInstanceFromDBSnapshotInput {
	snapID, _ := getInstanceSnapID(cr)
	restoreDBInput := &rds.RestoreDBInstanceFromDBSnapshotInput{
		AutoMinorVersionUpgrade: cr.Spec.CreateInstanceSpec.AutoMinorVersionUpgrade,
		AvailabilityZone:        cr.Spec.CreateInstanceSpec.AvailabilityZone,
		CopyTagsToSnapshot:      aws.Bool(true),
		DBInstanceClass:         cr.Spec.CreateInstanceSpec.DBInstanceClass,
		DBInstanceIdentifier:    cr.Spec.CreateInstanceSpec.DBInstanceIdentifier,
		DBSubnetGroupName:       cr.Spec.CreateInstanceSpec.DBSubnetGroupName,
		DeletionProtection:      cr.Spec.CreateInstanceSpec.DeletionProtection,
		Engine:                  cr.Spec.CreateInstanceSpec.Engine,
		DBSnapshotIdentifier:    &snapID,
	}
	return restoreDBInput

}
