package dbcluster

import (
	"errors"
	"strings"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) updateLocalStatusWithAwsStatus(cr *kubev1alpha1.DBCluster, clusterID string) (string, error) {

	exists, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: clusterID})
	currentLocalPhase := cr.Status.CurrentPhase

	if exists {
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in AWS: %v", cr.Name, cr.Namespace, *out.DBClusters[0].Status)
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in CR: %v", cr.Name, cr.Namespace, currentLocalPhase)

		if currentLocalPhase != strings.ToLower(*out.DBClusters[0].Status) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", cr.Namespace)
			cr.Status.CurrentPhase = strings.ToLower(*out.DBClusters[0].Status)
			if err := lib.UpdateCrStatus(r.client, cr); err != nil {
				return "", err
			}
		}
	}
	return cr.Status.CurrentPhase, nil

}

func (r *ReconcileDBCluster) handlePhases(cr *kubev1alpha1.DBCluster, clusterID string) error {

	// always update first before checking ( so restore and delete can be handled )
	currentPhase, _ := r.updateLocalStatusWithAwsStatus(cr, clusterID)

	switch currentPhase {
	case "available":
		return nil
	case "creating", "backing-up", "restoring":
		return &lib.ErrorResourceCreatingInProgress{Message: "ClusterCreatingInProgress"}
	case "deleting":
		return &lib.ErrorResourceDeletingInProgress{Message: "ClusterDeletingInProgress"}
	case "":
		return errors.New("ClusterNotYetInitilaized")
	}
	return nil
}
