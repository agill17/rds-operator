package controller

import (
	"github.com/agill17/rds-operator/pkg/controller/dbsubnetgroup"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, dbsubnetgroup.Add)
}
