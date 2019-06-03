package dbinstance

import (
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/service/rds"
)

func (r *ReconcileDBInstance) createNewDBInstance(cr *kubev1alpha1.DBInstance, insObj rdsLib.RDS) (*rds.CreateDBInstanceOutput, error) {

	var createOut *rds.CreateDBInstanceOutput
	var err error
	dbInsID := *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier

	if err := r.waitForClusterIfNeeded(cr); err != nil {
		return nil, err
	}

	// checks if it exists first
	err = insObj.Create()
	if err != nil {
		return nil, err
	}

	if err := r.handlePhases(cr); err != nil {
		return nil, err
	}

	cr.Status.Created = true
	_, rdsInstanceStatus := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: r.rdsClient, InstanceID: dbInsID})
	cr.Status.RDSInstanceStatus = rdsInstanceStatus

	// update status
	if err = lib.UpdateCrStatus(r.client, cr); err != nil {
		logrus.Errorf("Failed to update cr status for DBInstance: %v", err)
		return nil, err
	}
	logrus.Infof("CreateDBInstance was successful. Updated status")

	return createOut, err
}

// assuming instance is part of cluster, then use this func to wait until cluster is ready
// only valid when on fresh installsm restore from snapshot is done on clusters
func (r *ReconcileDBInstance) waitForClusterIfNeeded(cr *kubev1alpha1.DBInstance) error {
	var err error
	dbInsID := *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier
	// when cluster is still not available, this will throw ErrorClusterCreatingInProgress
	// only run this when this DBInstance is part of a DBCluster
	if cr.Spec.CreateInstanceSpec.DBClusterIdentifier != nil && !cr.Status.DBClusterMarkedAvail {
		dbClsID := *cr.Spec.CreateInstanceSpec.DBClusterIdentifier
		logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Part of cluster: %v -- checking if its available first", cr.Namespace, dbInsID, dbClsID)
		err = r.dbClusterReady(*cr.Spec.CreateInstanceSpec.DBClusterIdentifier)
		if err != nil {
			return err
		}
		cr.Status.DBClusterMarkedAvail = true
		return lib.UpdateCrStatus(r.client, cr)
	}
	return err
}
