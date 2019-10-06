package rdsLib

import (
	"errors"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/agill17/rds-operator/pkg/utils"
)

type RDSAction string

const (
	CREATE  RDSAction = "new"
	RESTORE RDSAction = "restoreFromSnapshot"
	DELETE  RDSAction = "delete"
	UNKNOWN RDSAction = "unknown"
	RECOVER RDSAction = "recoverWhenDeleted"
)

type RDS interface {
	Create() error
	Delete() error
	Restore() error
	SyncAwsStatusWithCRStatus() (string, error)
}

// AWSPhaseHandler is a generic func to update CR status with aws resource status
// also throws errors if a aws resource is not available/ready
// this is being used as a checkpoint to make sure aws resource is ready and available before performing more actions on it
// for example;
// - we do not want to deploy a k8s svc if db is not yet ready
// - we do not want to run a initDB job if db is not yet ready
func AWSPhaseHandler(rds RDS) error {
	// always update first before checking ( so restore and delete can be handled )
	crStatus, _ := rds.SyncAwsStatusWithCRStatus()

	// hack to set correct error messages
	var msgPrefix string
	switch rds.(type) {
	case *cluster:
		msgPrefix = "Cluster"
	case *instance:
		msgPrefix = "Instance"
	case *subnetGroup:
		msgPrefix = "SubnetGroup"
	}

	switch crStatus {
	case "available":
		return nil
	case "creating", "backing-up", "restoring", "modifying":
		return utils.ErrorResourceCreatingInProgress{Message: msgPrefix + "CreatingInProgress"}
	case "deleting":
		return utils.ErrorResourceDeletingInProgress{Message: msgPrefix + "DeletingInProgress"}
	case "":
		return errors.New(msgPrefix + "NotYetInitilaized")
	}

	return nil
}

func Crud(rdsObject RDS, actionType RDSAction, crStatusCreated bool, client client.Client) error {

	switch actionType {

	// fresh install
	case CREATE:

		if err := rdsObject.Create(); err != nil {
			return err
		}

		// delete event
	case DELETE:
		err := rdsObject.Delete()
		if err != nil {
			return err
		}

		// restore ( means different for each object )
	case RESTORE:
		if err := rdsObject.Restore(); err != nil {
			return err
		}

	}

	// return err if not ready in AWS yet
	if err := AWSPhaseHandler(rdsObject); err != nil {
		return err
	}


	return nil
}
