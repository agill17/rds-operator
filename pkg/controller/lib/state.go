package lib

import (
	"github.com/aws/aws-sdk-go/service/rds"
)

// deployed, creating, failed, notCreated

type CentralState struct {
	InstanceCreated bool
	ClusterCreated  bool
	ClusterID       string
	DBID            string
	InstanceOutput  *rds.DescribeDBInstancesOutput
	ClusterOutput   *rds.DescribeDBClustersOutput
}
