package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/rds"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DBClusterSpec defines the desired state of DBCluster
type DBClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Azs                         []string     `json:"azs"`
	BacktrackWindow             int64        `json:"backtrackWindow"`
	BackupRetentionPeriod       int64        `json:"backupRetentionPeriod"`
	DBClusterParameterGroupName string       `json:"dBClusterParameterGroupName,omitempty"`
	DBSubnetGroupName           string       `json:"dBSubnetGroupName"`
	DatabaseName                string       `json:"dbName"`
	DeletionProtection          bool         `json:"deletionProtection"`
	Engine                      string       `json:"engine"`
	EngineMode                  string       `json:"engineMode"`
	EngineVersion               string       `json:"engineVersion"`
	MasterUsername              string       `json:"masterUsername"`
	MasterPassword              string       `json:"masterPassword"`
	StorageEncrypted            bool         `json:"storageEncrypted"`
	VpcSecurityGroupIds         []string     `json:"vpcSecurityGroupIds"`
	DeletePolicy                DeletePolicy `json:"deletePolicy"`
}

// DBClusterStatus defines the observed state of DBCluster
type DBClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Created          bool                          `json:"created"`
	RDSClusterStatus *rds.DescribeDBClustersOutput `json:"rdsClusterStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBCluster is the Schema for the dbclusters API
// +k8s:openapi-gen=true
type DBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBClusterSpec   `json:"spec,omitempty"`
	Status DBClusterStatus `json:"status"`
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
