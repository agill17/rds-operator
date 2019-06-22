package dbcluster

import (
	"fmt"
	"strings"
	"time"

	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBCluster) setUpDefaultsIfNeeded(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) error {
	if err := r.setCRDeleteClusterID(cr, rdsAction); err != nil {
		return err
	}

	if err := r.setCRDeleteSpecSnapName(cr, rdsAction); err != nil {
		return err
	}

	return r.setSvcName(cr)
}

// used when deleteSpec.DBClusterID is not set, than use the one provided within cr.createClusterSpec
// this is so we can make the clusterID optional in deleteSpec
func (r *ReconcileDBCluster) setCRDeleteClusterID(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) error {
	var id string
	if cr.Spec.DeleteSpec.DBClusterIdentifier == nil || *cr.Spec.DeleteSpec.DBClusterIdentifier == "" {
		id = getDBClusterID(cr, rdsAction)
		logrus.Warnf("Setting spec.DeleteClusterSpec.DBClusterIdentifier: %v", id)
		cr.Spec.DeleteSpec.DBClusterIdentifier = &id
		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR while setting up DeleteSpec.DBClusterIdentifier: %v", err)
			return err
		}
	}

	return nil
}

func (r *ReconcileDBCluster) setCRDeleteSpecSnapName(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) error {
	var clusterID string

	if !(*cr.Spec.DeleteSpec.SkipFinalSnapshot) && cr.Spec.DeleteSpec.FinalDBSnapshotIdentifier == nil {
		clusterID = getDBClusterID(cr, rdsAction)
		currentTime := time.Now().Format("2006-01-02:03-02-44")
		snashotName := fmt.Sprintf("%v-%v", clusterID, strings.Replace(currentTime, ":", "-", -1))
		cr.Spec.DeleteSpec.FinalDBSnapshotIdentifier = aws.String(snashotName)

		return lib.UpdateCr(r.client, cr)
	}
	return nil
}

func (r *ReconcileDBCluster) setSvcName(cr *kubev1alpha1.DBCluster) error {
	if cr.Spec.ServiceName == "" {
		cr.Spec.ServiceName = fmt.Sprintf("%v-%v-rds-cluster", cr.Name, cr.Namespace)
		return lib.UpdateCr(r.client, cr)
	}
	return nil
}
