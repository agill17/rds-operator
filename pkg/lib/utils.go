package lib

import (
	"context"
	"math/rand"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LetterBytes[rand.Intn(len(LetterBytes))]
	}
	return string(b)
}

func WaitForExistence(waitType string, dbID string, ns string, rdsClient *rds.RDS) error {
	var err error
	dbInput := &rds.DescribeDBInstancesInput{DBInstanceIdentifier: &dbID}
	logrus.Warningf("Namespace: %v | DB Identifier: %v | Msg: Waiting for DB instance to become %v", ns, dbID, waitType)
	switch waitType {
	case "available":
		err = rdsClient.WaitUntilDBInstanceAvailableWithContext(context.Background(), dbInput)
		break
	case "notAvailable":
		err = rdsClient.WaitUntilDBInstanceDeletedWithContext(context.Background(), dbInput)
		break
	}

	if err != nil {
		logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: ERROR while waiting for rds db instance to %v", ns, dbID, waitType)
	}
	return err
}

// UpdateCr is used to update CR spec ( not status ), things like adding/removing finalizers, spec updates
func UpdateCr(client client.Client, object runtime.Object) error {
	if err := client.Update(context.TODO(), object); err != nil {
		logrus.Errorf("Failed to update CR obejct: ~~> %v", err)
		return err
	}
	return nil
}

// UpdateCrStatus is only used for updating status subresource in a CR object
func UpdateCrStatus(client client.Client, object runtime.Object) error {
	if err := client.Status().Update(context.TODO(), object); err != nil {
		logrus.Errorf("Failed to update status for CR obejct: ~~> %v", err)
		return err
	}
	return nil
}
