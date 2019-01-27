package lib

import (
	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/utils"
)

const (
	DefaultMasterUsername        = "admin"
	DefaultMasterPassword        = "admin"
	DefaultAz                    = "us-east-1a"
	DefaultStorageType           = "gp2"
	DefaultinstanceClass         = "db.t2.micro"
	DefaultStorageSize     int64 = 10
	DefaultSubnetGroupName       = "default"
	DefaultTimeoutForJob         = 300
	DBClusterFinalizer           = "agill.apps.dbCluster"
	DBInstanceFinalizer          = "agill.apps.dbInstance"
	DBSubnetGroupFinalizer       = "agill.apps.dbSubnetGroup"
)

// use for cluster and instance specs
func SetDBID(ns, crName string) string {
	return ns + "-" + crName
}

func GetUsernamePassword(cr *agillv1alpha1.DBInstance) (string, string) {
	username := utils.RandStringBytes(8)
	password := utils.RandStringBytes(8)

	if cr.Spec.MasterUsername != "" {
		username = cr.Spec.MasterUsername
	}

	if cr.Spec.MasterPassword != "" {
		password = cr.Spec.MasterPassword
	}
	return username, password
}
