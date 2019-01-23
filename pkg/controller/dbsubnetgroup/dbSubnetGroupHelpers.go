package dbsubnetgroup

import (
	"context"

	"github.com/sirupsen/logrus"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/aws/aws-sdk-go/service/rds"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileDBSubnetGroup) subnetGroupExists(name string) (bool, *rds.DescribeDBSubnetGroupsOutput) {
	out, err := r.rdsClient.DescribeDBSubnetGroups(&rds.DescribeDBSubnetGroupsInput{DBSubnetGroupName: &name})
	if err != nil {
		return false, nil // does not exist with that name
	}
	return true, out
}

func (r *ReconcileDBSubnetGroup) createSubnetGroup(request reconcile.Request) error {
	instance := &agillv1alpha1.DBSubnetGroup{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return err
	}
	createInput := r.createDBSubnetGroupInput(instance)

	_, err = r.rdsClient.CreateDBSubnetGroup(createInput)
	if err != nil {
		logrus.Errorf("Something went wrong while creating db subnet group: %v", err)
		return err
	}

	// update status
	instance.Status.Created = true
	if err := r.client.Update(context.TODO(), instance); err != nil {
		logrus.Errorf("Something went wrong while update CR status for DBSubnetGroup: %v", err)
		return err
	}

	return err

}
