package state

import "github.com/aws/aws-sdk-go/service/rds"

type State struct {
	DBClusterState  DBClusterState  `json:"dbClusterState"`
	DBInstanceState DBInstanceState `json:"dbInstanceState"`
}

type DBClusterState struct {
	ID               string                        `json:"id"`
	Created          bool                          `json:"created"`
	RecreateIt       bool                          `json:"recreateIt"`
	RDSClusterStatus *rds.DescribeDBClustersOutput `json:"rdsClusterStatus"`
}

type DBInstanceState struct {
	ID                string                         `json:"id"`
	Created           bool                           `json:"created"`
	RecreateIt        bool                           `json:"recreateIt"`
	RDSInstanceStatus *rds.DescribeDBInstancesOutput `json:"rdsInstanceStatus"`
}
