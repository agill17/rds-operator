package rdsLib

import (
	"strings"

	"github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type subnetGroup struct {
	createIn   *rds.CreateDBSubnetGroupInput
	deleteIn   *rds.DeleteDBSubnetGroupInput
	runtimeObj *v1alpha1.DBSubnetGroup
	rdsClient  *rds.RDS
	k8sClient  client.Client
}

func NewSubnetGroup(createIn *rds.CreateDBSubnetGroupInput, deleteIn *rds.DeleteDBSubnetGroupInput,
	runtimeObj *v1alpha1.DBSubnetGroup, rdsClient *rds.RDS, k8sClient client.Client) RDS {
	return &subnetGroup{
		createIn:   createIn,
		deleteIn:   deleteIn,
		runtimeObj: runtimeObj,
		rdsClient:  rdsClient,
		k8sClient:  k8sClient,
	}
}

func (s *subnetGroup) Create() error {
	exists, _ := lib.DBSubnetGroupExists(lib.RDSGenerics{SubnetGroupName: *s.createIn.DBSubnetGroupName, RDSClient: s.rdsClient})
	if !exists {
		logrus.Infof("Creating DBSubnetGroup for namespace: %v", s.runtimeObj.Namespace)
		if _, err := s.rdsClient.CreateDBSubnetGroup(s.createIn); err != nil {
			logrus.Errorf("Something went wrong while creating db subnet group: %v", err)
			return err
		}
	}
	return nil
}

func (s *subnetGroup) Delete() error {
	exists, _ := lib.DBSubnetGroupExists(lib.RDSGenerics{SubnetGroupName: *s.createIn.DBSubnetGroupName, RDSClient: s.rdsClient})
	if exists {
		logrus.Infof("Deleting DBSubnetGroup for namespace: %v", s.runtimeObj.Namespace)
		if _, err := s.rdsClient.DeleteDBSubnetGroup(s.deleteIn); err != nil {
			logrus.Errorf("Something went wrong while deleting db subnet group name: %v", err)
			return err
		}
		return lib.RemoveFinalizer(s.runtimeObj, s.k8sClient, lib.DBSubnetGroupFinalizer)
	}
	return nil
}

func (s *subnetGroup) Restore() error {
	return s.Create()
}

func (s *subnetGroup) SyncAwsStatusWithCRStatus() (string, error) {
	exists, out := lib.DBSubnetGroupExists(lib.RDSGenerics{RDSClient: s.rdsClient, SubnetGroupName: *s.createIn.DBSubnetGroupName})
	currentLocalPhase := s.runtimeObj.Status.CurrentPhase
	if exists {
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in AWS: %v", s.runtimeObj.Name, s.runtimeObj.Namespace, *out.DBSubnetGroups[0].SubnetGroupStatus)
		logrus.Infof("DBCluster CR: %v | Namespace: %v | Current phase in CR: %v", s.runtimeObj.Name, s.runtimeObj.Namespace, currentLocalPhase)
		if currentLocalPhase != strings.ToLower(*out.DBSubnetGroups[0].SubnetGroupStatus) {
			logrus.Warnf("Updating current phase in CR for namespace: %v", s.runtimeObj.Namespace)
			s.runtimeObj.Status.CurrentPhase = strings.ToLower(*out.DBSubnetGroups[0].SubnetGroupStatus)
			if err := lib.UpdateCrStatus(s.k8sClient, s.runtimeObj); err != nil {
				return "", err
			}
		}
	}
	return s.runtimeObj.Status.CurrentPhase, nil
}
