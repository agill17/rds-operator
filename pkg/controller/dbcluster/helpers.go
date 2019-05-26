package dbcluster

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/agill17/rds-operator/pkg/lib/dbHelpers"

	"context"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
	kubeapierror "k8s.io/apimachinery/pkg/api/errors"
)

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
			if err := lib.UpdateCrStatus(r.client, cr); err != nil {
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
		return &lib.ErrorResourceCreatingInProgress{Message: "ClusterCreatingInProgress"}
	case "deleting":
		return &lib.ErrorResourceDeletingInProgress{Message: "ClusterDeletingInProgress"}
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
	if cr.Spec.CreateClusterSpec == nil && cr.Spec.CreateClusterFromSnapshot == nil {
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

func getInstallType(cr *kubev1alpha1.DBCluster) dbHelpers.DBInstallType {
	if cr.GetDeletionTimestamp() != nil && len(cr.GetFinalizers()) > 0 {
		return dbHelpers.DELETE
	} else if cr.Spec.CreateClusterFromSnapshot != nil {
		return dbHelpers.RESTORE
	} else if cr.Spec.CreateClusterSpec != nil {
		return dbHelpers.CREATE
	}

	return dbHelpers.UNKNOWN
}

func (r *ReconcileDBCluster) createSecret(cr *kubev1alpha1.DBCluster) error {
	secretObj := r.getSecretObj(cr)
	if err := r.client.Create(context.TODO(), secretObj); err != nil && !kubeapierror.IsAlreadyExists(err) {
		logrus.Errorf("Error while creating secret object: %v", err)
		return err
	} else if kubeapierror.IsAlreadyExists(err) && cr.Status.SecretUpdateNeeded {
		logrus.Warnf("Updating cluster secret in namespace: %v", cr.Namespace)
		r.client.Update(context.TODO(), secretObj)
		cr.Status.SecretUpdateNeeded = false
		return lib.UpdateCrStatus(r.client, cr)
	}

	return nil
}
