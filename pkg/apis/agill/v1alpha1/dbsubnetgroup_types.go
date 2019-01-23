package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DBSubnetGroupSpec defines the desired state of DBSubnetGroup
type DBSubnetGroupSpec struct {
	SubnetIds []string `json:"subnetIds"`
}

// DBSubnetGroupStatus defines the observed state of DBSubnetGroup
type DBSubnetGroupStatus struct {
	Created bool `json:"created"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBSubnetGroup is the Schema for the dbsubnetgroups API
// +k8s:openapi-gen=true
type DBSubnetGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBSubnetGroupSpec   `json:"spec,omitempty"`
	Status DBSubnetGroupStatus `json:"status,omitempty"`
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
