package lib

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
)

func GetLatestClusterSnapID(clusterDBID, ns, region string) (string, error) {
	cmd := fmt.Sprintf("aws rds describe-db-cluster-snapshots  --query \"DBClusterSnapshots[?DBClusterIdentifier=='%v']\" --region %v | jq -r 'max_by(.SnapshotCreateTime).DBClusterSnapshotIdentifier'", clusterDBID, region)
	snapID, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		logrus.Errorf("Failed to execute command: %s", err)
		return "", err
	}

	logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Latest snapshot id available: %v", ns, clusterDBID, strings.TrimSpace(string(snapID)))

	return strings.TrimSpace(string(snapID)), err
}

func GetLatestSnapID(dbID, ns string) (string, error) {
	cmd := fmt.Sprintf("aws rds describe-db-snapshots --query \"DBSnapshots[?DBInstanceIdentifier=='%v']\" --region us-east-1 | jq -r 'max_by(.SnapshotCreateTime).DBSnapshotIdentifier'", dbID)
	snapID, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		logrus.Errorf("Failed to execute command: %s", err)
		return "", err
	}

	logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: Latest snapshot id available: %v", ns, dbID, strings.TrimSpace(string(snapID)))

	return strings.TrimSpace(string(snapID)), err
}

func WaitForExistence(waitType string, dbID string, ns string, rdsClient *rds.RDS) error {
	var err error
	dbInput := &rds.DescribeDBInstancesInput{DBInstanceIdentifier: &dbID}
	logrus.Warningf("Namespace: %v | DB Identifier: %v | Msg: Waiting for DB instance to become %v", ns, dbID, waitType)
	switch waitType {
	case "available":
		err = rdsClient.WaitUntilDBInstanceAvailable(dbInput)
		break
	case "notAvailable":
		err = rdsClient.WaitUntilDBInstanceDeleted(dbInput)
		break
	}

	if err != nil {
		logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: ERROR while waiting for rds db instance to %v", ns, dbID, waitType)
	}
	return err
}

func GetImportJobCmd(dbEngine, dbName, username, password, endpoint, sqlFile string) []string {
	var actualCmd string
	if strings.ToLower(dbEngine) == "mysql" || strings.ToLower(dbEngine) == "aurora-mysql" {
		actualCmd = fmt.Sprintf("mysql -u %v -p%v -h %v %v < %v", username, password, endpoint, dbName, sqlFile)
	}
	// switch strings.ToLower(dbEngine) {
	// case "mysql":
	// 	actualCmd = fmt.Sprintf("mysql -u %v -p%v -h %v %v < %v", username, password, endpoint, dbName, sqlFile)
	// case "postgresql":
	// 	/** TODO **/
	// }
	return []string{"sh", "-c", actualCmd}
}

func GetTags(metaLabels map[string]string) []*rds.Tag {
	var tags []*rds.Tag

	for k, v := range metaLabels {
		tags = append(tags, &rds.Tag{Key: &k, Value: &v})
	}

	return tags
}
