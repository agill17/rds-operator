package lib

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

type RDSGenerics struct {
	InstanceID, ClusterID, SubnetGroupName string
	RDSClient                              *rds.RDS
}

// Returns bool, DescribeDBClustersOutput
func DbClusterExists(r RDSGenerics) (bool, *rds.DescribeDBClustersOutput, error) {
	exists := true
	output, err := r.RDSClient.DescribeDBClusters(&rds.DescribeDBClustersInput{
		DBClusterIdentifier: &r.ClusterID,
	})

	if err != nil {
		return false, nil, err
	}
	return exists, output, nil
}

func DBInstanceExists(r RDSGenerics) (bool, *rds.DescribeDBInstancesOutput) {
	exists := true

	output, err := r.RDSClient.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &r.InstanceID,
	})

	if err != nil {
		logrus.Errorf("Failed to check if DBInstance %v exists: %v", r.InstanceID, err)
		exists = false
	}
	return exists, output
}

func DBSubnetGroupExists(r RDSGenerics) (bool, *rds.DescribeDBSubnetGroupsOutput) {
	out, err := r.RDSClient.DescribeDBSubnetGroups(&rds.DescribeDBSubnetGroupsInput{DBSubnetGroupName: &r.SubnetGroupName})
	if err != nil{
		return false, nil // does not exist with that name
	}
	return true, out
}
