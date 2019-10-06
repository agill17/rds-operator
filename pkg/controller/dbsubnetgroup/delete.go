package dbsubnetgroup

import (
	"github.com/agill17/rds-operator/pkg/utils"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

func (r *ReconcileDBSubnetGroup) deleteSubnetGroup(name string) error {
	var err error
	if exists, _ := utils.DBSubnetGroupExists(utils.RDSGenerics{RDSClient: r.rdsClient, SubnetGroupName: name}); exists {
		logrus.Warnf("DBSubentGroup %v exists, going to delete now.", name)
		_, err := r.rdsClient.DeleteDBSubnetGroup(&rds.DeleteDBSubnetGroupInput{DBSubnetGroupName: &name})
		if err != nil {
			logrus.Errorf("Something went wrong while deleting db subnet group name: %v", err)
			return err
		}
	}
	return err
}
