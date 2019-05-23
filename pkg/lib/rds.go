package lib

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSGenerics struct {
	InstanceID, ClusterID, SubnetGroupName string
	RDSClient                              *rds.RDS
}

// Returns bool, DescribeDBClustersOutput
func DbClusterExists(r *RDSGenerics) (bool, *rds.DescribeDBClustersOutput) {
	exists := true
	output, err := r.RDSClient.DescribeDBClusters(&rds.DescribeDBClustersInput{
		DBClusterIdentifier: &r.ClusterID,
	})

	if err != nil && err.(awserr.Error).Code() == rds.ErrCodeDBClusterNotFoundFault {
		exists = false
	}
	return exists, output
}

func DBInstanceExists(r *RDSGenerics) (bool, *rds.DescribeDBInstancesOutput) {
	exists := true

	output, err := r.RDSClient.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &r.InstanceID,
	})

	if err != nil && err.(awserr.Error).Code() == rds.ErrCodeDBInstanceNotFoundFault {
		exists = false
	}
	return exists, output
}

func DBSubnetGroupExists(r *RDSGenerics) (bool, *rds.DescribeDBSubnetGroupsOutput) {
	out, err := r.RDSClient.DescribeDBSubnetGroups(&rds.DescribeDBSubnetGroupsInput{DBSubnetGroupName: &r.SubnetGroupName})
	if err != nil && err.(awserr.Error).Code() == rds.ErrCodeDBSubnetGroupNotFoundFault {
		return false, nil // does not exist with that name
	}
	return true, out
}
