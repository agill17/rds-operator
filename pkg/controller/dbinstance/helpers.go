package dbinstance

import (
	"errors"
	"strings"

	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
)

// throws ErrorResourceCreatingInProgress when dbCluster in AWS is not marked available
func (r *ReconcileDBInstance) dbClusterReady(clusterID string) error {
	var err error
	exists, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: clusterID})
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
		logrus.Infof("DBInstance CR: %v | Namespace: %v | Current phase in AWS: %v", cr.Name, cr.Namespace, *out.DBInstances[0].DBInstanceStatus)
		logrus.Infof("DBInstance CR: %v | Namespace: %v | Current phase in CR: %v", cr.Name, cr.Namespace, currentLocalPhase)

		if currentLocalPhase != strings.ToLower(*out.DBInstances[0].DBInstanceStatus) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", cr.Namespace)
			cr.Status.CurrentPhase = strings.ToLower(*out.DBInstances[0].DBInstanceStatus)
			if err := lib.UpdateCrStatus(r.client, cr); err != nil {
				return "", err
			}
		}
	}
	return cr.Status.CurrentPhase, nil

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
		sName = cr.Name + "-instance-service"
	}
	return sName
}

func getActionType(cr *kubev1alpha1.DBInstance) rdsLib.RDSAction {
	if cr.GetDeletionTimestamp() != nil && len(cr.GetFinalizers()) > 0 {
		return rdsLib.DELETE
	} else if cr.Spec.CreateInstanceSpec != nil {
		return rdsLib.CREATE
	} else if cr.Spec.RestoreInstanceFromSnap != nil {
		return rdsLib.RESTORE
	}

	return rdsLib.UNKNOWN
}

func isPartOfCluster(cr *kubev1alpha1.DBInstance) bool {

	if cr.Spec.CreateInstanceSpec.DBClusterIdentifier != nil {
		return true
	}
	return false
}
