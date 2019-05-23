package lib

type ClusterInstallType string

const (
	DBInstanceFinalizer                              = "agill.apps.dbInstance"
	DBClusterFinalizer                               = "agill.apps.dbCluster"
	DBSubnetGroupFinalizer                           = "agill.apps.dbSubnetGroupFinalizer"
	LetterBytes                                      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultRegion                                    = "us-east-1"
	CLUSTER_INSTALL_NEW           ClusterInstallType = "newInstall"
	CLUSTER_INSTALL_FROM_SNAPSHOT ClusterInstallType = "newInstallFromSnapshop"
)
