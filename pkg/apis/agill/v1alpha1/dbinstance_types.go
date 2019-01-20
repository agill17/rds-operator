package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/rds"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DBInstanceSpec defines the desired state of DBInstance
type DBInstanceSpec struct {
	Engine                   string       `json:"engine"`
	EngineVersion            string       `json:"engineVersion"`
	PubliclyAccessible       bool         `json:"publiclyAccessible"`
	Az                       string       `json:"az"`
	AutoMinorVersionUpgrade  bool         `json:"autoMinorVersionUpgrade"`
	AllocatedStorage         int64        `json:"allocatedStorage,omitempty"`
	BackupRetentionPeriod    int64        `json:"backupRetentionPeriod,omitempty"`
	DBInstanceClass          string       `json:"dbInstanceClass"`
	DBName                   string       `json:"dbName"`
	DBClusterIdentifier      string       `json:"dbClusterId,omitempty"`
	DBParameterGroupName     string       `json:"dBParameterGroupName,omitempty"`
	DBSecurityGroups         []string     `json:"dBSecurityGroups"`
	DBSubnetGroupName        string       `json:"dBSubnetGroupName"`
	DeletionProtection       bool         `json:"deletionProtection"`
	MasterUsername           string       `json:"masterUsername"`
	MasterPassword           string       `json:"masterPassword"`
	StorageEncrypted         bool         `json:"storageEncrypted,omitempty"`
	VpcSecurityGroupIds      []string     `json:"vpcSecurityGroupIds,omitempty"`
	RehealFromLatestSnapshot bool         `json:"rehealFromLatestSnapshot"`
	DBSecretName             string       `json:"dbSecretName"`
	InitDBJob                InitDBJob    `json:"initDBJob,omitempty"`
	ExternalSvcName          string       `json:"externalSvcName"`
	DeletePolicy             DeletePolicy `json:"deletePolicy"`
}

type DeletePolicy struct {
	DeleteAutomatedBackups bool `json:"deleteAutomatedBackups"`
	SkipFinalSnapshot      bool `json:"skipFinalSnapshot"`
}

// DBInstanceStatus defines the observed state of DBInstance
type DBInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	RDSInstanceStatus        *rds.DescribeDBInstancesOutput `json:"rdsInstanceStatus"`
	DeployedInitially        bool                           `json:"deployedInitially"`
	RestoredFromSnapshotName string                         `json:"restoredFromSnapshotName"`
	UpdateKubeFiles          bool                           `json:"updateKubeFiles"`
	InitJobTimedOut          bool                           `json:"initJobTimedOut"`
	InitJobSuccessfull       bool                           `json:"initJobSuccessfull"`
}

//
type InitDBJob struct {
	Image             string            `json:"image"`
	WaitTillCompleted bool              `json:"waitTillCompleted"`
	Timeout           int               `json:"timeout"`
	BackOffLimit      int               `json:"backOffLimit"`
	ImagePullSecret   string            `json:"imagePullSecret,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	NodeSelector      map[string]string `json:"nodeSelector,omitempty"`
	SQLFile           string            `json:"sqlFile"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBInstance is the Schema for the dbinstances API
// +k8s:openapi-gen=true
type DBInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBInstanceSpec   `json:"spec,omitempty"`
	Status DBInstanceStatus `json:"status,omitempty"`
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
