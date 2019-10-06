package dbcluster

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"

	"github.com/aws/aws-sdk-go/service/rds"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/agill17/rds-operator/pkg/utils"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

// used to set env keys in initDBJob
const (
	dbClusterEndpointEnvVar = "DB_CLUSTER_ENDPOINT"
	dbClusterUsername       = "DB_CLUSTER_USERNAME"
	dbClusterPassword       = "DB_CLUSTER_PASSWORD"
)

/*
 TODO: Add notes about reconcileInitDBJob
*/

func reconcileInitDBJob(cr *kubev1alpha1.DBCluster, client client.Client, rdsClient *rds.RDS) error {
	// if defined then proceed
	if cr.InitClusterDB.Image != "" {

		if err := setPrimaryInstanceID(cr, rdsClient, client); err != nil {
			return err
		}

		instanceAvail, err := isPrimaryInstanceAvailable(cr, rdsClient)
		if err != nil {
			return err
		}

		if instanceAvail {

			// set up a spec
			jobSpec := &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "init-db",
					Namespace: cr.Namespace,
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							RestartPolicy: v1.RestartPolicyNever,
							NodeSelector:     cr.InitClusterDB.NodeSelector,
							ImagePullSecrets: cr.InitClusterDB.ImagePullSecrets,
							Containers: []v1.Container{{
								Name:  "init-db",
								Image: cr.InitClusterDB.Image,
								Env: populateEnvVarsInCr(cr),
								Command: cr.InitClusterDB.Command,
							}},
						},
					},
				},
			}

			if err := client.Create(context.TODO(), jobSpec); err != nil && !apiErrors.IsAlreadyExists(err) {
				logrus.Errorf("Error while creating initDB job in namespace: %v -- %v", cr.Namespace, err)
				return err
			}

			// check if completed
			// get the job deployed inside cluster
			job := &batchv1.Job{}
			err = client.Get(context.TODO(), types.NamespacedName{Name: "init-db", Namespace: cr.Namespace}, job)
			if err != nil {
				return err
			}

			// ensure job is no longer active before moving on
			if job.Status.Active >= 1 {
				return errors.New("InitDBJobStillRunning")
			}

		}
	}
	return nil
}

// This should be run ONLY after cluster is in available state AND has 1 instance in available state
func populateEnvVarsInCr(cr *kubev1alpha1.DBCluster) []v1.EnvVar {
	secretName := cr.Status.SecretName
	userKey := cr.Status.UsernameKey
	passKey := cr.Status.PasswordKey
	clusterEndpoint := *cr.Status.DescriberClusterOutput.DBClusters[0].Endpoint



	return []v1.EnvVar{
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
}

/*
	if initDBJobDefined,
	then we need to run the job AFTER the primary DBInstance
	is marked available. We cannot run initDBJob until then
	So repeat the DescribeDBCluster func until we have 1
	primary instance available and set it inside the DBCluster
	CR so we dont have to do this check over and over again
*/

func setPrimaryInstanceID(cr *kubev1alpha1.DBCluster, rdsClient *rds.RDS, client client.Client) error {
	clusterAvail := cr.Status.CurrentPhase == "available"
	if cr.Status.PrimaryInstanceID == "" && clusterAvail {

		// describe cluster to get instance members
		_, out,_ := utils.DbClusterExists(utils.RDSGenerics{RDSClient: rdsClient, ClusterID: getDBClusterID(cr)})
		if len(out.DBClusters[0].DBClusterMembers) == 0 {
			return utils.ErrorNoDBInstanceAttachedToClusterYet{Message: "There are no db instances attached to " + *cr.ClusterSpec.DBClusterIdentifier + " cluster yet",}
		}
		for _, eachMember := range out.DBClusters[0].DBClusterMembers {
			if *eachMember.IsClusterWriter {
				logrus.Infof("Namespace: %v | CR: %v | Found instance that is part of cluster to run initDBJob", cr.Namespace, cr.Name)
				cr.Status.PrimaryInstanceID = *eachMember.DBInstanceIdentifier
				return utils.UpdateCrStatus(client, cr)
			}
		}
	}
	return nil
}

// get instanceID from status, if set, check instance status, if available return true
func isPrimaryInstanceAvailable(cr *kubev1alpha1.DBCluster, rdsClient *rds.RDS) (bool, error) {

	if cr.Status.PrimaryInstanceID == "" {
		return false, utils.ErrorResourceCreatingInProgress{Message: "DBCluster is not aware of a instance ID yet. Please attach a DBInstance if you want to run initDBJob"}
	}

	dbInstanceExists, instanceOut := utils.DBInstanceExists(utils.RDSGenerics{RDSClient: rdsClient, InstanceID: cr.Status.PrimaryInstanceID})
	if dbInstanceExists {
		if *instanceOut.DBInstances[0].DBInstanceStatus == "available" {
			return true, nil
		}
	}

	return false, utils.ErrorResourceCreatingInProgress{Message: "DBCluster does not have a primary dbInstance in available state yet"}
}
