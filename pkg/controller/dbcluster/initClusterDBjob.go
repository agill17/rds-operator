package dbcluster

import (
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/rds"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/agill17/rds-operator/pkg/lib"

	v1 "k8s.io/api/core/v1"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
)

func reconcileInitDBJob(cr *kubev1alpha1.DBCluster, client client.Client, rdsClient *rds.RDS) error {
	if err := setPrimaryInstanceID(cr, client); err != nil {
		return err
	}

	instanceAvail, err := isPrimaryInstanceAvailable(cr, rdsClient)
	if err != nil {
		return err
	}

	if instanceAvail {
		logrus.Infof("Applying initDBJob Spec..")
	}

	return nil
}

func setPrimaryInstanceID(cr *kubev1alpha1.DBCluster, client client.Client) error {
	clusterMembers := cr.Status.DescriberClusterOutput.DBClusters[0].DBClusterMembers
	for _, memeber := range clusterMembers {
		if *memeber.IsClusterWriter {
			cr.Status.PrimaryInstanceID = *memeber.DBInstanceIdentifier
			return lib.UpdateCrStatus(client, cr)
		}
	}
	return nil
}

func isPrimaryInstanceAvailable(cr *kubev1alpha1.DBCluster, rdsClient *rds.RDS) (bool, error) {

	if cr.Status.PrimaryInstanceID == "" {
		return false, &lib.ErrorResourceCreatingInProgress{Message: "DBCluster is not aware of a instance ID yet. Please attach a DBInstance if you want to run initDBJob"}
	}

	dbInstanceExists, instanceOut := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: rdsClient, InstanceID: cr.Status.PrimaryInstanceID})
	if dbInstanceExists {
		if *instanceOut.DBInstances[0].DBInstanceStatus == "available" {
			return true, nil
		} else {
			return false, &lib.ErrorResourceCreatingInProgress{Message: "DBCluster does not have a dbInstance in available state yet"}
		}
	}

	return false, nil
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
