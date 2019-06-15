package rdsLib

import (
	"strings"
)

type RDSAction string
type RDS_RESOURCE_STATE string

const (
	RDS_AVAILABLE  RDS_RESOURCE_STATE = "available"
	RDS_CREATING   RDS_RESOURCE_STATE = "creating"
	RDS_DELETING   RDS_RESOURCE_STATE = "deleting"
	RDS_RESTORING  RDS_RESOURCE_STATE = "restoring"
	RDS_BACKING_UP RDS_RESOURCE_STATE = "backing-up"
	RDS_UNKNOWN    RDS_RESOURCE_STATE = "unknown"
	CREATE         RDSAction          = "new"
	RESTORE        RDSAction          = "restoreFromSnapshot"
	DELETE         RDSAction          = "delete"
	UNKNOWN        RDSAction          = "unknown"
)

type RDS interface {
	Create() error
	Delete() error
	Restore() error
	GetAWSStatus() RDS_RESOURCE_STATE
}

// take aws remote status ( returned in string ) and turn it into our type so we can do type checking
func parseRemoteStatus(rs string) RDS_RESOURCE_STATE {
	switch strings.ToLower(rs) {
	case "available":
		return RDS_AVAILABLE
	case "deleting":
		return RDS_DELETING
	case "creating":
		return RDS_CREATING
	case "backing-up":
		return RDS_BACKING_UP
	default:
		return RDS_UNKNOWN
	}
}
