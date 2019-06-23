package dbinstance

import (
	"strings"

	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
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

func getInstanceID(cr *kubev1alpha1.DBInstance) string {
	if cr.Spec.CreateInstanceSpec != nil {
		return *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier
	} else if cr.Spec.RestoreInstanceFromSnap != nil {
		return *cr.Spec.RestoreInstanceFromSnap.DBInstanceIdentifier
	}

	return ""
}

// assuming instance is part of cluster, then use this func to wait until cluster is ready
// only valid when on fresh installs
func (r *ReconcileDBInstance) waitForClusterIfNeeded(cr *kubev1alpha1.DBInstance) error {
	var err error
	dbInsID := getInstanceID(cr)
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
