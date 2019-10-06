package v1alpha1

import (
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DBSubnetGroupSpec defines the desired state of DBSubnetGroup
type DBSubnetGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// DBSubnetGroupStatus defines the observed state of DBSubnetGroup
type DBSubnetGroupStatus struct {
	CurrentPhase   string `json:"currentPhase"`
	Created        bool   `json:"created"`
	RecreateNeeded bool   `json:"recreateNeeded"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBSubnetGroup is the Schema for the dbsubnetgroups API
// +k8s:openapi-gen=true
type DBSubnetGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//Spec   *rds.CreateDBSubnetGroupInput `json:"createDBSubnetGroupSpec,omitempty"`
	Status DBSubnetGroupStatus           `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBSubnetGroupList contains a list of DBSubnetGroup
type DBSubnetGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBSubnetGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBSubnetGroup{}, &DBSubnetGroupList{})
}

func (in *DBSubnetGroup) DeepCopyInto(out *DBSubnetGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
	copier.Copy(&in, &out)

	return
}
