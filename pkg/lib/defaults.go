package lib

type ResourcePhase string

const (
	DBInstanceFinalizer    = "agill.apps.dbInstance"
	DBClusterFinalizer     = "agill.apps.dbCluster"
	DBSubnetGroupFinalizer = "agill.apps.dbSubnetGroupFinalizer"
	LetterBytes            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultRegion          = "us-east-1"
)
