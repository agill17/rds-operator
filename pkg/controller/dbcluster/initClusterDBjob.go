package dbcluster

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/rds"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/agill17/rds-operator/pkg/lib"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

// used to set env keys in initDBJob
const (
	dbClusterEndpointEnvVar = "DB_CLUSTER_ENDPOINT"
	dbClusterUsername       = "DB_CLUSTER_USERNAME"
	dbClusterPassword       = "DB_CLUSTER_PASSWORD"
)

func reconcileInitDBJob(cr *kubev1alpha1.DBCluster, client client.Client, rdsClient *rds.RDS) error {

	// if defined then proceed
	if cr.Spec.InitClusterDB.Spec != nil {

		if err := setPrimaryInstanceID(cr, rdsClient, client); err != nil {
			return err
		}

		instanceAvail, err := isPrimaryInstanceAvailable(cr, rdsClient)
		if err != nil {
			return err
		}

		if instanceAvail {
			jobSpec := &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%v-initdb", cr.Name),
					Namespace: cr.Namespace,
				},
				Spec: *cr.Spec.InitClusterDB.Spec,
			}

			if err := client.Create(context.TODO(), jobSpec); err != nil && !apiErrors.IsAlreadyExists(err) {
				logrus.Errorf("Error while creating initDB job in namespace: %v -- %v", cr.Namespace, err)
				return err
			}

		}
	}
	return nil
}

// get instanceID from status, if set, check instance status, if available return true
func isPrimaryInstanceAvailable(cr *kubev1alpha1.DBCluster, rdsClient *rds.RDS) (bool, error) {

	if cr.Status.PrimaryInstanceID == "" {
		return false, lib.ErrorResourceCreatingInProgress{Message: "DBCluster is not aware of a instance ID yet. Please attach a DBInstance if you want to run initDBJob"}
	}

	dbInstanceExists, instanceOut := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: rdsClient, InstanceID: cr.Status.PrimaryInstanceID})
	if dbInstanceExists {
		if *instanceOut.DBInstances[0].DBInstanceStatus == "available" {
			return true, nil
		}
	}

	return false, lib.ErrorResourceCreatingInProgress{Message: "DBCluster does not have a primary dbInstance in available state yet"}
}

// This should be run ONLY after cluster is in available state AND has 1 instance in available state
func (r *ReconcileDBCluster) populateEnvVarsInCr(cr *kubev1alpha1.DBCluster) error {
	secretName := cr.Status.SecretName
	userKey := cr.Status.UsernameKey
	passKey := cr.Status.PasswordKey
	clusterEndpoint := *cr.Status.DescriberClusterOutput.DBClusters[0].Endpoint
	jobSpec := cr.Spec.InitClusterDB

	if len(jobSpec.Spec.Template.Spec.Containers) == 0 {
		return errors.New("InitJobContainerEmptyError")
	}

	currentEnvs := jobSpec.Spec.Template.Spec.Containers[0].Env

	envs := []v1.EnvVar{
		{
			Name:  dbClusterEndpointEnvVar,
			Value: clusterEndpoint,
		},
		{
			Name: dbClusterUsername,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secretName,
					},
					Key: userKey,
				},
			},
		},
		{
			Name: dbClusterPassword,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secretName,
					},
					Key: passKey,
				},
			},
		},
	}

	if !reflect.DeepEqual(currentEnvs, envs) {
		jobSpec.Spec.Template.Spec.Containers[0].Env = envs
		return lib.UpdateCr(r.client, cr)
	}

	return nil
}

/*
	if initDBJobDefined,
	then we need to run that AFTER the primary DBInstance
	is marked available. We cannot run initDBJob until then
	So repeat the DescribeDBCluster func until we have 1
	primary instance available and set it inside the DBCluster
	CR so we dont have to do this check over and over again
*/

func setPrimaryInstanceID(cr *kubev1alpha1.DBCluster, rdsClient *rds.RDS, client client.Client) error {
	clusterAvail := cr.Status.CurrentPhase == "available"
	if cr.Status.PrimaryInstanceID == "" && clusterAvail {
		// describe cluster to get instance members
		_, out := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: rdsClient, ClusterID: getDBClusterID(cr)})
		if len(out.DBClusters[0].DBClusterMembers) == 0 {
			return errors.New("NoDBInstancesAttahcedToCluster")
		}
		for _, eachMember := range out.DBClusters[0].DBClusterMembers {
			if *eachMember.IsClusterWriter {
				logrus.Infof("Namespace: %v | CR: %v | Found instance that is part of cluster to run initDBJob", cr.Namespace, cr.Name)
				cr.Status.PrimaryInstanceID = *eachMember.DBInstanceIdentifier
				return lib.UpdateCrStatus(client, cr)
			}
		}
	}
	return nil
}
