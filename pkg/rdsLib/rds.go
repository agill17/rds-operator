package rdsLib

import (
	"errors"

	"github.com/agill17/rds-operator/pkg/lib"
)

type RDSAction string

const (
	CREATE  RDSAction = "new"
	RESTORE RDSAction = "restoreFromSnapshot"
	DELETE  RDSAction = "delete"
	UNKNOWN RDSAction = "unknown"
)

type RDS interface {
	Create() error
	Delete() error
	Restore() error
	GetAWSStatus() (string, error)
}

// SyncAndReconcileIfNotReady is a generic func to update CR status with aws resource status
// also throws generic errors if a aws resource is not available/ready
// this is being used as a checkpoint to make sure aws resource is ready and available before performing more actions on it
// for example;
// - we do not want to deploy a k8s svc if db is not yet ready
// - we do not want to run a initDB job if db is not yet ready
func SyncAndReconcileIfNotReady(rds RDS) error {
	// always update first before checking ( so restore and delete can be handled )
	currentPhase, _ := rds.GetAWSStatus()

	// hack to set correct error messages
	var msgPrefix string
	switch rds.(type) {
	case *cluster:
		msgPrefix = "Cluster"
	case *instance:
		msgPrefix = "Instance"
	}

	switch currentPhase {
	case "available":
		return nil
	case "creating", "backing-up", "restoring", "modifying":
		return &lib.ErrorResourceCreatingInProgress{Message: msgPrefix + "CreatingInProgress"}
	case "deleting":
		return &lib.ErrorResourceDeletingInProgress{Message: msgPrefix + "DeletingInProgress"}
	case "":
		return errors.New(msgPrefix + "NotYetInitilaized")
	}

	return nil
}
