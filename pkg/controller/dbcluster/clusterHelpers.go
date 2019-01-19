package dbcluster

import (
	"context"

	"github.com/davecgh/go-spew/spew"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileDBCluster) dbClusterExists(dbClusterID string) (bool, *rds.DescribeDBClustersOutput) {
	exists := true
	output, err := r.rdsClient.DescribeDBClusters(&rds.DescribeDBClustersInput{
		DBClusterIdentifier: &dbClusterID,
	})

	if err != nil && err.(awserr.Error).Code() == rds.ErrCodeDBClusterNotFoundFault {
		exists = false
	}
	return exists, output
}

func (r *ReconcileDBCluster) deleteCluster(cr agillv1alpha1.DBCluster, dbID string) error {
	var err error

	if exists, _ := r.dbClusterExists(dbID); exists {
		if _, err = r.rdsClient.DeleteDBCluster(&rds.DeleteDBClusterInput{
			DBClusterIdentifier: &dbID,
			SkipFinalSnapshot:   &cr.Spec.DeletePolicy.SkipFinalSnapshot,
		}); err != nil {
			logrus.Errorf("ERROR deleting cluster: %v", err)
			return err
		}

	}

	return err
}

func (r *ReconcileDBCluster) createCluster(cr *agillv1alpha1.DBCluster, dbID string, request reconcile.Request) error {
	if exists, _ := r.dbClusterExists(dbID); !exists {
		logrus.Infof("Creating db cluster first")
		input := getCreateDBClusterInput(cr, dbID)
		_, err := r.rdsClient.CreateDBCluster(input)
		if err != nil {
			logrus.Errorf("ERROR while creating DB Cluster%v:", err)
			spew.Dump(cr)
			return err
		}

		instance := &agillv1alpha1.DBCluster{}
		err = r.client.Get(context.TODO(), request.NamespacedName, instance)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}
		_, output := r.dbClusterExists(dbID)
		instance.Status.RDSClusterStatus = output
		instance.Status.Created = true
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return err
		}

	}
	return nil
}
