package lib

import (
	"context"
	"math/rand"

	"k8s.io/apimachinery/pkg/api/meta"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
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

func SecretExists(namespace, secretName string, client client.Client) (bool, *v1.Secret) {
	secretFound := &v1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{
		Name: secretName, Namespace: namespace},
		secretFound)
	if err != nil && errors.IsNotFound(err) {
		return false, nil
	}

	return true, secretFound
}

func AddFinalizer(runtimeObj runtime.Object, client client.Client, finalizer string) error {
	// get the runtime obj interface so I can add finalizers in metadata
	// note: Accessor returns meta.Object which is an interface with funcs to muck around with
	// k8s object metadata fields ONLY
	accessor, err := meta.Accessor(runtimeObj)
	if err != nil {
		return nil
	}

	currentFinalizers := accessor.GetFinalizers()

	if !finalizerExists(currentFinalizers, finalizer) {
		currentFinalizers = append(currentFinalizers, finalizer)
		accessor.SetFinalizers(currentFinalizers)
		return UpdateCr(client, runtimeObj)
	}

	return nil

}

func finalizerExists(currentList []string, lookupFinalizer string) bool {
	for _, e := range currentList {
		if e == lookupFinalizer {
			return true
		}
	}
	return false
}
