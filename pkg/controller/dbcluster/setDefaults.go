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

	if err := r.setUpCredentialsIfNeeded(cr, rdsAction); err != nil {
		return err
	}

	if err := r.setCRDeleteClusterID(cr, rdsAction); err != nil {
		return err
	}

	if err := r.setCRDeleteSpecSnapName(cr, rdsAction); err != nil {
		return err
	}

	return nil
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

func (r *ReconcileDBCluster) setCRUsername(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) error {

	if rdsAction == rdsLib.CREATE && cr.Spec.CreateClusterSpec.MasterUsername == nil {
		if r.useCredentialsFrom(cr) {
			if exists, secretObj := lib.SecretExists(
				cr.Namespace,
				cr.Spec.CredentialsFrom.SecretName.Name,
				r.client); exists {
				su := string(secretObj.Data[cr.Spec.CredentialsFrom.UsernameKey])
				cr.Spec.CreateClusterSpec.MasterUsername = &su
			}
		} else {
			u := lib.RandStringBytes(9)
			cr.Spec.CreateClusterSpec.MasterUsername = &u
		}

		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR while setting up username: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBCluster) setCRPassword(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) error {
	if rdsAction == rdsLib.CREATE && cr.Spec.CreateClusterSpec.MasterUserPassword == nil {
		if r.useCredentialsFrom(cr) {
			if exists, secretObj := lib.SecretExists(
				cr.Namespace,
				cr.Spec.CredentialsFrom.SecretName.Name,
				r.client); exists {
				sp := string(secretObj.Data[cr.Spec.CredentialsFrom.PasswordKey])
				cr.Spec.CreateClusterSpec.MasterUserPassword = &sp
			}
		} else {
			p := lib.RandStringBytes(9)
			cr.Spec.CreateClusterSpec.MasterUserPassword = &p
		}

		if err := lib.UpdateCr(r.client, cr); err != nil {
			logrus.Errorf("Failed to update DBCluster CR while setting up credentials: %v", err)
			return err
		}
	}
	return nil
}

func (r *ReconcileDBCluster) useCredentialsFrom(cr *kubev1alpha1.DBCluster) bool {
	if cr.Spec.CredentialsFrom.UsernameKey != "" && cr.Spec.CredentialsFrom.PasswordKey != "" {
		logrus.Infof("Using spec.credentialsFrom for username and password")
		return true
	}
	logrus.Infof("Checking credentials, if not passed random string will be generated and dumped into clusterSecret")
	return false
}

func (r *ReconcileDBCluster) setUpCredentialsIfNeeded(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) error {
	if err := r.setCRUsername(cr, rdsAction); err != nil {
		return err
	}

	if err := r.setCRPassword(cr, rdsAction); err != nil {
		return err
	}

	return nil
}
