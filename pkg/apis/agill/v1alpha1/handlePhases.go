package v1alpha1

// type HandlePhase interface {
// 	SyncRemoteStateWithCR(client client.Client, rdsClient *rds.RDS) error
// }

// func (i *DBInstance) SyncRemoteStateWithCR(client client.Client, rdsClient *rds.RDS) error {

// 	exists, out := lib.DBInstanceExists(lib.RDSGenerics{RDSClient: rdsClient, InstanceID: *i.Spec.DBInstanceIdentifier})
// 	currentLocalPhase := i.Status.CurrentPhase

// 	if exists {
// 		logrus.Infof("DBInstance: Current phase in AWS: %v", *out.DBInstances[0].DBInstanceStatus)
// 		logrus.Infof("DBInstance: Current phase in CR: %v", currentLocalPhase)

// 		if currentLocalPhase != strings.ToLower(*out.DBInstances[0].DBInstanceStatus) {
// 			logrus.Warnf("Updating current phase in CR for namespace: %v", i.Namespace)
// 			i.Status.CurrentPhase = strings.ToLower(*out.DBInstances[0].DBInstanceStatus)
// 			if err := lib.UpdateCrStatus(client, i); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func (i *DBCluster) SyncRemoteStateWithCR(client client.Client, rdsClient *rds.RDS) error {

// 	exists, out := lib.DBInstanceExists(lib.RDSGenerics{RDSClient: rdsClient, InstanceID: *i.Spec.DBClusterIdentifier})
// 	currentLocalPhase := i.Status.CurrentPhase

// 	if exists {
// 		logrus.Infof("DBInstance: Current phase in AWS: %v", *out.DBInstances[0].DBInstanceStatus)
// 		logrus.Infof("DBInstance: Current phase in CR: %v", currentLocalPhase)

// 		if currentLocalPhase != strings.ToLower(*out.DBInstances[0].DBInstanceStatus) {
// 			logrus.Warnf("Updating current phase in CR for namespace: %v", i.Namespace)
// 			i.Status.CurrentPhase = strings.ToLower(*out.DBInstances[0].DBInstanceStatus)
// 			if err := lib.UpdateCrStatus(client, i); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
