package dbsubnetgroup

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBSubnetGroup) deleteSubnetGroup(name string) error {
	var err error
	if exists, _ := r.subnetGroupExists(name); exists {
		_, err := r.rdsClient.DeleteDBSubnetGroup(&rds.DeleteDBSubnetGroupInput{DBSubnetGroupName: &name})
		if err != nil {
			logrus.Errorf("Something went wrong while deleting db subnet group name: %v", err)
			return err
		}
		logrus.Infof("Successfully deleted db subnetgroup: %v", name)
	}
	return err
}
