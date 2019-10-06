package v1alpha1

import (
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DBInstanceSpec defines the desired state of DBInstance
type DBInstanceSpec struct {
	Region                  string                                    `json:"region"`
	InstanceSecretName      string                                    `json:"instanceSecretName,omitempty"`
	ServiceName             string                                    `json:"serviceName,omitempty"`
	//CreateInstanceSpec      *rds.CreateDBInstanceInput                `json:"createInstanceSpec,omitempty"`
	//RestoreInstanceFromSnap *rds.RestoreDBInstanceFromDBSnapshotInput `json:"createInstanceFromSnapshot,omitempty"`
	//DeleteInstanceSpec      *rds.DeleteDBInstanceInput                `json:"deleteInstanceSpec,required"`
}

/*
	When using dbInstance with dbCluster why do we need a deleteSpec here?
	What about snapshotting dbInstance when attached to DBCluster

	Read here;
	https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_DeleteInstance.html
	"
	If you try to delete the cluster's last DB instance,
	the behavior depends on the method you use.
	You can't delete the last DB instance through the AWS Management Console,
	because doing so also deletes the cluster.

	You can delete the last DB instance through
	the AWS CLI or API even if the DB cluster has deletion protection enabled.
	In this case, the DB cluster itself still exists and your data is preserved.
	You can access the data by attaching new DB instances to the cluster
	"

	So why deleteSpec in dbInstance?
	- To support standalone rds instances ( not aurora-* )
*/

// DBInstanceStatus defines the observed state of DBInstance
type DBInstanceStatus struct {
	DBClusterMarkedAvail     bool                           `json:"dbClusterMarkedAvail"`
	//RDSInstanceStatus        *rds.DescribeDBInstancesOutput `json:"rdsInstanceStatus"`
	Created                  bool                           `json:"created"`
	RestoredFromSnapshotName string                         `json:"restoredFromSnapshotName"`
	InitJobTimedOut          bool                           `json:"initJobTimedOut"`
	InitJobSuccessfull       bool                           `json:"initJobSuccessfull"`
	Username                 string                         `json:"username"`
	Password                 string                         `json:"password"`
	CurrentPhase             string                         `json:"currentPhase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBInstance is the Schema for the dbinstances API
// +k8s:openapi-gen=true
type DBInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DBInstanceSpec   `json:"spec,required"`
	Status            DBInstanceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBInstanceList contains a list of DBInstance
type DBInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBInstance{}, &DBInstanceList{})
}

func (in *DBInstance) DeepCopyInto(out *DBInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
	copier.Copy(&in.Spec, &out.Spec)

	return
}
