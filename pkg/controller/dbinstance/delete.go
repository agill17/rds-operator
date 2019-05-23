package dbinstance

import (
	"fmt"
	"strings"
	"time"

	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"

	// h "cloud.google.com/go/bigquery/benchmarks"
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
)

func (r *ReconcileDBInstance) deleteDBInstance(cr *kubev1alpha1.DBInstance, dbID string) error {
	var err error

	deleteInput := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier:   &dbID,
		SkipFinalSnapshot:      &cr.DeleteInstance.SkipFinalSnapshot,
		DeleteAutomatedBackups: &cr.DeleteInstance.DeleteAutomatedBackups,
	}

	if isStandAlone(cr) {
		currentTime := time.Now().Format("2006-01-02:03-02-44")
		deleteSnapID := fmt.Sprintf("%v-%v", *cr.Spec.DBInstanceIdentifier, strings.Replace(currentTime, ":", "-", -1))
		deleteInput.FinalDBSnapshotIdentifier = &deleteSnapID
	}

	if exists, stat := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: r.rdsClient, InstanceID: dbID}); exists {

		logrus.Infof("Does DBInstanceID :%v exists in AWS: %v", dbID, exists)
		// is it already in deleting state?
		if *stat.DBInstances[0].DBInstanceStatus == "deleting" {
			logrus.Warnf("DBInstance is already in deleting state: %v", dbID)
			return nil
		}

		logrus.Warnf("Namespace: %v | DB Identifier: %v | Msg: Starting to delete db rds instance", cr.Namespace, dbID)
		_, err = r.rdsClient.DeleteDBInstance(deleteInput)
		if err != nil && err.(awserr.Error).Code() != rds.ErrCodeInvalidDBInstanceStateFault {
			logrus.Errorf("ERROR while deleting db instance: %v", err)
			return err
		}
	} else if !exists {
		logrus.Infof("DB instance does not exist in AWS, skipping delete of RDS: %v", dbID)
	}
	logrus.Infof("Successfully deleted DBInstance for Namespace: %v", cr.Namespace)
	return err
}
