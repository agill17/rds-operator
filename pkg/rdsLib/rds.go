package rdsLib

type RDSAction string

const (
	CREATE  RDSAction = "new"
	RESTORE RDSAction = "fromSnapshot"
	DELETE  RDSAction = "delete"
	UNKNOWN RDSAction = "unknown"
)

type RDS interface {
	Create() error
	Delete() error
	Restore() error
}

func InstallRestoreDelete(dbInput RDS, action RDSAction) error {
	switch action {
	case CREATE:
		return dbInput.Create()
	case DELETE:
		return dbInput.Delete()
	case RESTORE:
		return dbInput.Restore()
	}

	return nil
}
