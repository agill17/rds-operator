package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jinzhu/copier"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DBInstanceSpec defines the desired state of DBInstance
type DBInstanceSpec struct {
	Region             string                     `json:"region"`
	InstanceSecretName string                     `json:"instanceSecretName,omitempty"`
	ServiceName        string                     `json:"serviceName,omitempty"`
	InitDB             InitDB                     `json:"initDB,omitempty"`
	CreateInstanceSpec *rds.CreateDBInstanceInput `json:"createInstanceSpec,required"`
	DeleteInstanceSpec DeleteInstanceSpec         `json:"deleteInstanceSpec,required"`
}

// DBInstanceStatus defines the observed state of DBInstance
type DBInstanceStatus struct {
	DBClusterMarkedAvail     bool                           `json:"dbClusterMarkedAvail"`
	RDSInstanceStatus        *rds.DescribeDBInstancesOutput `json:"rdsInstanceStatus"`
	DeployedInitially        bool                           `json:"deployedInitially"`
	RestoredFromSnapshotName string                         `json:"restoredFromSnapshotName"`
	UpdateKubeFiles          bool                           `json:"updateKubeFiles"`
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

type DeleteInstanceSpec struct {
	DeleteAutomatedBackups    bool   `json:"DeleteAutomatedBackups"`
	SkipFinalSnapshot         bool   `json:"SkipFinalSnapshot"`
	FinalDBSnapshotIdentifier string `json:"FinalDBSnapshotIdentifier,omitempty"`
}

type InitDB struct {
	Image             string            `json:"image"`
	WaitTillCompleted bool              `json:"waitTillCompleted"`
	Timeout           int               `json:"timeout"`
	BackOffLimit      int               `json:"backOffLimit"`
	ImagePullSecret   string            `json:"imagePullSecret,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	NodeSelector      map[string]string `json:"nodeSelector,omitempty"`
	SQLFile           string            `json:"sqlFile"`
	Command           []string          `json:"command"`
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
