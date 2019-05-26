package dbcluster

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) updateCrStats(cr *kubev1alpha1.DBCluster) error {
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		logrus.Errorf("Failed while updating status of DBCluster in namespace: %v -- %v", cr.Namespace, err)
		return err
	}
	return err
}

func (r *ReconcileDBCluster) getCurrentStatusFromAWS(dbClusterID string) string {
	exists, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: dbClusterID})
	if exists {
		return *out.DBClusters[0].Status
	}
	return ""
}

// use for cluster and instance specs
func (r *ReconcileDBCluster) setDBID(ns, crName string) string {
	return ns + "-" + crName
}

func getLatestClusterSnapID(clusterDBID, ns, region string) (string, error) {
	cmd := fmt.Sprintf("aws rds describe-db-cluster-snapshots  --query \"DBClusterSnapshots[?DBClusterIdentifier=='%v']\" --region %v | jq -r 'max_by(.SnapshotCreateTime).DBClusterSnapshotIdentifier'", clusterDBID, region)
	snapID, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		logrus.Errorf("Failed to execute aws-cli command: %s", err)
		return "", err
	}

	logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Latest snapshot id available: %v", ns, clusterDBID, strings.TrimSpace(string(snapID)))

	return strings.TrimSpace(string(snapID)), err
}

func (r *ReconcileDBCluster) updateLocalStatusWithAwsStatus(cr *kubev1alpha1.DBCluster) (string, error) {

	exists, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.CreateClusterSpec.DBClusterIdentifier})
	currentLocalPhase := cr.Status.CurrentPhase

	if exists {
		logrus.Infof("DBCluster: Current phase in AWS: %v", *out.DBClusters[0].Status)
		logrus.Infof("DBCluster: Current phase in CR: %v", currentLocalPhase)

		if currentLocalPhase != strings.ToLower(*out.DBClusters[0].Status) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", cr.Namespace)
			cr.Status.CurrentPhase = strings.ToLower(*out.DBClusters[0].Status)
			if err := r.updateCrStats(cr); err != nil {
				return "", err
			}
		}
	}
	return cr.Status.CurrentPhase, nil

}

func (r *ReconcileDBCluster) handlePhases(cr *kubev1alpha1.DBCluster) error {

	// always update first before checking ( so restore and delete can be handled )
	currentPhase, _ := r.updateLocalStatusWithAwsStatus(cr)

	switch currentPhase {
	case "available":
		return nil
	case "creating", "backing-up":
		return lib.ErrorResourceCreatingInProgress{Message: "ClusterCreatingInProgress"}
	case "deleting":
		return lib.ErrorResourceDeletingInProgress{Message: "ClusterDeletingInProgress"}
	case "":
		return errors.New("ClusterNotYetInitilaized")
	}
	return nil
}

func getClusterSecretName(cr *kubev1alpha1.DBCluster) string {
	name := cr.ClusterSecretName
	if len(name) == 0 {
		name = cr.Name + "-secret"
	}
	return name
}

func validateRequiredInput(cr *kubev1alpha1.DBCluster) error {
	if cr.Spec.CreateClusterSpec == nil {
		return errors.New("CreateClusterSpecEmptyError")
	}

	if cr.Region == "" {
		return errors.New("regionCannotBeEmptyError")
	}

	if cr.Spec.DeleteSpec == nil {
		return errors.New("deleteClusterSpecCannotBeEmptyError")
	}
	return nil
}
