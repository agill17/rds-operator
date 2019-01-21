package dbinstance

import (
	"github.com/agill17/rds-operator/pkg/controller/lib"
	// h "cloud.google.com/go/bigquery/benchmarks"
	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBInstance) deleteDBInstance(cr *agillv1alpha1.DBInstance, dbID string) error {
	var err error
	var deleteInput *rds.DeleteDBInstanceInput
	if exists, _ := r.dbInstanceExists(dbID); exists {
		deleteInput = &rds.DeleteDBInstanceInput{
			DeleteAutomatedBackups: &cr.Spec.DeletePolicy.DeleteAutomatedBackups,
			DBInstanceIdentifier:   &dbID,
			SkipFinalSnapshot:      &cr.Spec.DeletePolicy.SkipFinalSnapshot,
		}
		logrus.Warnf("Namespace: %v | DB Identifier: %v | Msg: Starting to delete db rds instance", cr.Namespace, dbID)
		spew.Dump(cr.Spec.DeletePolicy)
		_, err = r.rdsClient.DeleteDBInstance(deleteInput)
		if err != nil && err.(awserr.Error).Code() != rds.ErrCodeInvalidDBInstanceStateFault {
			logrus.Errorf("ERROR while deleting db instance: %v", err)
			return err
		} else if err := lib.WaitForExistence("notAvailable", dbID, cr.Namespace, r.rdsClient); err == nil {
			logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Successfully deleted DB from RDS", cr.Namespace, dbID)
		}
	} else if !exists {
		logrus.Infof("DB instance does not exist in AWS, skipping delete of RDS: %v", dbID)
	}
	return err
}
