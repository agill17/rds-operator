package dbinstance

import (
	"context"
	"errors"
	"strings"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
)

func crHasDBStatus(cr *kubev1alpha1.DBInstance) bool {
	var crHasSTatus bool
	if cr.Status.DeployedInitially {
		// cr has db status -- meaning db had already been deployed for this cr.
		crHasSTatus = true
	}
	return crHasSTatus
}

// throws ErrorResourceCreatingInProgress when dbCluster in AWS is not marked available
func (r *ReconcileDBInstance) dbClusterReady(cr *kubev1alpha1.DBInstance) error {
	var err error
	dbClsID := *cr.Spec.CreateInstanceSpec.DBClusterIdentifier
	exists, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: dbClsID})
	if exists {
		if strings.ToLower(*out.DBClusters[0].Status) != "available" {
			return &lib.ErrorResourceCreatingInProgress{Message: "ClusterCreatingInProgress"}
		}
	}

	return err
}
func (r *ReconcileDBInstance) updateLocalStatusWithAwsStatus(cr *kubev1alpha1.DBInstance) (string, error) {
	dbInsID := *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier
	exists, out := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: r.rdsClient, InstanceID: dbInsID})
	currentLocalPhase := cr.Status.CurrentPhase

	if exists {
		logrus.Infof("DBInstance: Current phase in AWS: %v", *out.DBInstances[0].DBInstanceStatus)
		logrus.Infof("DBInstance: Current phase in CR: %v", currentLocalPhase)

		if currentLocalPhase != strings.ToLower(*out.DBInstances[0].DBInstanceStatus) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", cr.Namespace)
			cr.Status.CurrentPhase = strings.ToLower(*out.DBInstances[0].DBInstanceStatus)
			if err := r.updateResourceStatus(cr); err != nil {
				return "", err
			}
		}
	}
	return cr.Status.CurrentPhase, nil

}
func (r *ReconcileDBInstance) updateResourceStatus(resource *kubev1alpha1.DBInstance) error {
	var err error
	err = r.client.Status().Update(context.TODO(), resource)
	if err != nil {
		logrus.Errorf("Failed to update status in DBInstance Controller: %v", err)
		return err
	}
	return err
}

func (r *ReconcileDBInstance) handlePhases(cr *kubev1alpha1.DBInstance) error {

	// always update first before checking ( so restore and delete can be handled )
	currentPhase, _ := r.updateLocalStatusWithAwsStatus(cr)

	switch currentPhase {
	case "available":
		return nil
	case "creating", "backing-up", "restoring":
		return &lib.ErrorResourceCreatingInProgress{Message: "InstanceCreatingInProgress"}
	case "deleting":
		return &lib.ErrorResourceDeletingInProgress{Message: "InstanceDeletingInProgress"}
	case "":
		return errors.New("InstanceNotYetInitilaized")
	}
	return nil
}

func validateRequiredParams(cr *kubev1alpha1.DBInstance) error {
	if cr.Spec.CreateInstanceSpec == nil {
		return errors.New("createInstanceSpecEmptyError")
	}
	return nil
}

func getSecretName(cr *kubev1alpha1.DBInstance) string {
	sName := cr.Spec.InstanceSecretName
	if sName == "" {
		sName = cr.Name + "-secret"
	}
	return sName
}

func getSvcName(cr *kubev1alpha1.DBInstance) string {
	sName := cr.Spec.ServiceName
	if sName == "" {
		sName = cr.Name + "-svc"
	}
	return sName
}

func isStandAlone(cr *kubev1alpha1.DBInstance) bool {

	// true: when not associated to DBCluster
	// false: when associated to DBCluster

	if cr.Spec.CreateInstanceSpec.DBClusterIdentifier == nil {
		return true
	}
	return false
}
