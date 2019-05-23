package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DBClusterStatus struct {
	Created                  bool                          `json:"created"`
	RestoreNeeded            bool                          `json:"restoreNeeded"`
	DescriberClusterOutput   *rds.DescribeDBClustersOutput `json:"describeClusterOutput"`
	RestoredFromSnapshotName string                        `json:"restoredFromSnapshotName"`
	SecretUpdateNeeded       bool                          `json:"secretUpdateNeeded"`
	Username                 string                        `json:"username"`
	Password                 string                        `json:"password"`
	CurrentPhase             string                        `json:"currentPhase"`
}

type DBClusterSpec struct {
	CreateClusterSpec *rds.CreateDBClusterInput `json:"createClusterSpec"`
	DeleteSpec        *rds.DeleteDBClusterInput `json:"deleteClusterSpec,required"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBCluster is the Schema for the dbclusters API
// +k8s:openapi-gen=true
type DBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DBClusterSpec   `json:"spec"`
	Status            DBClusterStatus `json:"status,omitempty"`
	ClusterSecretName string          `json:"clusterSecretName"`
	Region            string          `json:"region"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBClusterList contains a list of DBCluster
type DBClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBCluster{}, &DBClusterList{})
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBCluster) DeepCopyInto(out *DBCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	// in.Status.DeepCopyInto(&out.Status)
	copier.Copy(&in.Spec, &out.Spec)
	return
}
