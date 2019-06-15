package rdsLib

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
}
