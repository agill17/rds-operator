package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DBClusterStatus struct {
	Created                  bool                          `json:"created"`
	SecretName               string                        `json:"secretName"`
	UsernameKey              string                        `json:"usernameKey"`
	PasswordKey              string                        `json:"passwordKey"`
	PrimaryInstanceID        string                        `json:"primaryInstanceID"`
	DescriberClusterOutput   *rds.DescribeDBClustersOutput `json:"describeClusterOutput,omitempty"`
	RestoredFromSnapshotName string                        `json:"restoredFromSnapshotName"`
	CurrentPhase             string                        `json:"currentPhase"`
}

type ClusterSpec struct {
	// +optional
	CredentialsFrom CredentialsFrom `json:"credentialsFrom,omitempty"`

	// The identifier for the DB snapshot or DB cluster snapshot to restore from.
	//
	// You can use either the name or the Amazon Resource Name (ARN) to specify
	// a DB cluster snapshot. However, you can use only the ARN to specify a DB
	// snapshot.
	//
	// Constraints:
	//
	//    * Must match the identifier of an existing Snapshot.
	//
	// SnapshotIdentifier is a required field if restoring from an existing rds snapshot
	SnapshotIdentifier *string `json:"snapshotIdentifier,omitempty"`

	// A list of EC2 Availability Zones that instances in the DB cluster can be
	// created in. For information on AWS Regions and Availability Zones, see Choosing
	// the Regions and Availability Zones (http://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Concepts.RegionsAndAvailabilityZones.html)
	// in the Amazon Aurora User Guide.
	AvailabilityZones *[]string `json:"availabilityZones,omitempty"`

	// The target backtrack window, in seconds. To disable backtracking, set this
	// value to 0.
	//
	// Default: 0
	//
	// Constraints:
	//
	//    * If specified, this value must be set to a number from 0 to 259,200 (72
	//    hours).
	BacktrackWindow *int64 `json:"backtrackWindow,omitempty"`

	// The number of days for which automated backups are retained. You must specify
	// a minimum value of 1.
	//
	// Default: 1
	//
	// Constraints:
	//
	//    * Must be a value from 1 to 35
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=35
	BackupRetentionPeriod *int64 `json:"backupRetentionPeriod,omitempty"`

	// A value that indicates that the DB cluster should be associated with the
	// specified CharacterSet.
	CharacterSetName *string `json:"characterSetName,omitempty"`

	// The DB cluster identifier. This parameter is stored as a lowercase string.
	//
	// Constraints:
	//
	//    * Must contain from 1 to 63 letters, numbers, or hyphens.
	//
	//    * First character must be a letter.
	//
	//    * Can't end with a hyphen or contain two consecutive hyphens.
	//
	// Example: my-cluster1
	//
	// DBClusterIdentifier is a required field
	// +kubebuilder:validation:Maximum=63
	// +kubebuilder:validation:Minimum=1
	DBClusterIdentifier *string `json:"dbClusterIdentifier,required"`

	// The name of the DB cluster parameter group to associate with this DB cluster.
	// If this argument is omitted, default.aurora5.6 is used.
	//
	// Constraints:
	//
	//    * If supplied, must match the name of an existing DB cluster parameter
	//    group.
	DBClusterParameterGroupName *string `json:"dbBClusterParameterGroupName,omitempty"`

	// A DB subnet group to associate with this DB cluster.
	//
	// Constraints: Must match the name of an existing DBSubnetGroup. Must not be
	// default.
	//
	// Example: mySubnetgroup
	// +optional
	DBSubnetGroupName *string `json:"dbSubnetGroupName,omitempty"`

	// The name for your database of up to 64 alpha-numeric characters. If you do
	// not provide a name, Amazon RDS will not create a database in the DB cluster
	// you are creating.
	DatabaseName *string `json:"databaseName,required"`

	// Indicates if the DB cluster should have deletion protection enabled. The
	// database can't be deleted when this value is set to true. The default is
	// false.
	DeletionProtection *bool `json:"deletionProtection,required"`

	// DestinationRegion is used for presigning the request to a given region.
	DestinationRegion *string `json:"destinationRegion,omitempty"`

	// The list of log types that need to be enabled for exporting to CloudWatch
	// Logs. The values in the list depend on the DB engine being used. For more
	// information, see Publishing Database Logs to Amazon CloudWatch Logs (http://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_LogAccess.html#USER_LogAccess.Procedural.UploadtoCloudWatch)
	// in the Amazon Aurora User Guide.
	EnableCloudwatchLogsExports *[]string `json:"enableCloudwatchLogsExports,omitempty"`

	// True to enable mapping of AWS Identity and Access Management (IAM) accounts
	// to database accounts, and otherwise false.
	//
	// Default: false
	EnableIAMDatabaseAuthentication *bool `json:"enableIAMDatabaseAuthentication,omitempty"`

	// The name of the database engine to be used for this DB cluster.
	//
	// Valid Values: aurora (for MySQL 5.6-compatible Aurora), aurora-mysql (for
	// MySQL 5.7-compatible Aurora), and aurora-postgresql
	//
	// Engine is a required field
	// +kubebuilder:validation:Enum=aurora,aurora-mysql,aurora-postgresql
	Engine *string `json:"engine,required"`

	// The DB engine mode of the DB cluster, either provisioned, serverless, or
	// parallelquery.
	// +kubebuilder:validation:Enum=provisioned,serverless,parallelquery
	EngineMode *string `json:"engineMode,omitempty"`

	// The version number of the database engine to use.
	//
	// Aurora MySQL
	//
	// Example: 5.6.10a, 5.7.12
	//
	// Aurora PostgreSQL
	//
	// Example: 9.6.3
	EngineVersion *string `json:"engineVersion,required"`

	// The AWS KMS key identifier for an encrypted DB cluster.
	//
	// The KMS key identifier is the Amazon Resource Name (ARN) for the KMS encryption
	// key. If you are creating a DB cluster with the same AWS account that owns
	// the KMS encryption key used to encrypt the new DB cluster, then you can use
	// the KMS key alias instead of the ARN for the KMS encryption key.
	//
	// If an encryption key is not specified in KmsKeyId:
	//
	//    * If ReplicationSourceIdentifier identifies an encrypted source, then
	//    Amazon RDS will use the encryption key used to encrypt the source. Otherwise,
	//    Amazon RDS will use your default encryption key.
	//
	//    * If the StorageEncrypted parameter is true and ReplicationSourceIdentifier
	//    is not specified, then Amazon RDS will use your default encryption key.
	//
	// AWS KMS creates the default encryption key for your AWS account. Your AWS
	// account has a different default encryption key for each AWS Region.
	//
	// If you create a Read Replica of an encrypted DB cluster in another AWS Region,
	// you must set KmsKeyId to a KMS key ID that is valid in the destination AWS
	// Region. This key is used to encrypt the Read Replica in that AWS Region.
	KmsKeyId *string `json:"kmsKeyId,omitempty"`

	// The password for the master database user. This password can contain any
	// printable ASCII character except "/", """, or "@".
	//
	// Constraints: Must contain from 8 to 41 characters.
	// +kubebuilder:validation:Maximum=41
	// +kubebuilder:validation:Minimum=8
	MasterUserPassword *string `json:"masterUserPassword,omitempty"`

	// The name of the master user for the DB cluster.
	//
	// Constraints:
	//
	//    * Must be 1 to 16 letters or numbers.
	//
	//    * First character must be a letter.
	//
	//    * Can't be a reserved word for the chosen database engine.
	// +kubebuilder:validation:Maximum=17
	// +kubebuilder:validation:Minimum=2
	MasterUsername *string `json:"masterUsername,omitempty"`

	// A value that indicates that the DB cluster should be associated with the
	// specified option group.
	//
	// Permanent options can't be removed from an option group. The option group
	// can't be removed from a DB cluster once it is associated with a DB cluster.
	OptionGroupName *string `json:"optionGroup,omitempty"`

	// The port number on which the instances in the DB cluster accept connections.
	//
	// Default: 3306 if engine is set as aurora or 5432 if set to aurora-postgresql.
	Port *int64 `json:"port,omitempty"`

	// A URL that contains a Signature Version 4 signed request for the CreateDBCluster
	// action to be called in the source AWS Region where the DB cluster is replicated
	// from. You only need to specify PreSignedUrl when you are performing cross-region
	// replication from an encrypted DB cluster.
	//
	// The pre-signed URL must be a valid request for the CreateDBCluster API action
	// that can be executed in the source AWS Region that contains the encrypted
	// DB cluster to be copied.
	//
	// The pre-signed URL request must contain the following parameter values:
	//
	//    * KmsKeyId - The AWS KMS key identifier for the key to use to encrypt
	//    the copy of the DB cluster in the destination AWS Region. This should
	//    refer to the same KMS key for both the CreateDBCluster action that is
	//    called in the destination AWS Region, and the action contained in the
	//    pre-signed URL.
	//
	//    * DestinationRegion - The name of the AWS Region that Aurora Read Replica
	//    will be created in.
	//
	//    * ReplicationSourceIdentifier - The DB cluster identifier for the encrypted
	//    DB cluster to be copied. This identifier must be in the Amazon Resource
	//    Name (ARN) format for the source AWS Region. For example, if you are copying
	//    an encrypted DB cluster from the us-west-2 AWS Region, then your ReplicationSourceIdentifier
	//    would look like Example: arn:aws:rds:us-west-2:123456789012:cluster:aurora-cluster1.
	//
	// To learn how to generate a Signature Version 4 signed request, see  Authenticating
	// Requests: Using Query Parameters (AWS Signature Version 4) (http://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html)
	// and  Signature Version 4 Signing Process (http://docs.aws.amazon.com/general/latest/gr/signature-version-4.html).
	PreSignedUrl *string `json:"preSignedUrl,omitempty"`

	// The daily time range during which automated backups are created if automated
	// backups are enabled using the BackupRetentionPeriod parameter.
	//
	// The default is a 30-minute window selected at random from an 8-hour block
	// of time for each AWS Region. To see the time blocks available, see  Adjusting
	// the Preferred DB Cluster Maintenance Window (http://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_UpgradeDBInstance.Maintenance.html#AdjustingTheMaintenanceWindow.Aurora)
	// in the Amazon Aurora User Guide.
	//
	// Constraints:
	//
	//    * Must be in the format hh24:mi-hh24:mi.
	//
	//    * Must be in Universal Coordinated Time (UTC).
	//
	//    * Must not conflict with the preferred maintenance window.
	//
	//    * Must be at least 30 minutes.
	// +kubebuilder:validation:Pattern=hh24:mi-hh24:mi
	PreferredBackupWindow *string `json:"preferredBackupWindow,omitempty"`

	// The weekly time range during which system maintenance can occur, in Universal
	// Coordinated Time (UTC).
	//
	// Format: ddd:hh24:mi-ddd:hh24:mi
	//
	// The default is a 30-minute window selected at random from an 8-hour block
	// of time for each AWS Region, occurring on a random day of the week. To see
	// the time blocks available, see  Adjusting the Preferred DB Cluster Maintenance
	// Window (http://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/USER_UpgradeDBInstance.Maintenance.html#AdjustingTheMaintenanceWindow.Aurora)
	// in the Amazon Aurora User Guide.
	//
	// Valid Days: Mon, Tue, Wed, Thu, Fri, Sat, Sun.
	//
	// Constraints: Minimum 30-minute window.
	// +kubebuilder:validation:Pattern=ddd:hh24:mi-ddd:hh24:mi
	PreferredMaintenanceWindow *string `json:"preferredMaintenanceWindow,omitempty"`

	// The Amazon Resource Name (ARN) of the source DB instance or DB cluster if
	// this DB cluster is created as a Read Replica.
	ReplicationSourceIdentifier *string `json:"replicationSourceIdentifier,omitempty"`

	// SourceRegion is the source region where the resource exists. This is not
	// sent over the wire and is only used for presigning. This value should always
	// have the same region as the source ARN.
	SourceRegion *string `json:"sourceRegion,omitempty"`

	// Specifies whether the DB cluster is encrypted.
	StorageEncrypted *bool `json:"storageEncrypted,omitempty"`

	// A list of EC2 VPC security groups to associate with this DB cluster.
	VpcSecurityGroupIds *[]string `json:"vpcSecurityGroupIds,omitempty"`
}

type InitClusterDB struct {
	Spec *v1.JobSpec `json:"jobSpec,omitempty"`
}

type CredentialsFrom struct {
	UsernameKey string                       `json:"usernameKey,required"`
	PasswordKey string                       `json:"passwordKey,required"`
	SecretName  *corev1.LocalObjectReference `json:"secret,required"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DBCluster is the Schema for the dbclusters API
// +k8s:openapi-gen=true
type DBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	ClusterSpec       ClusterSpec `json:"clusterSpec,required"`
	// k8s serviceName to use to point to rds cluster endpoint
	// +kubebuilder:validation:MaxLength=15
	ServiceName string `json:"serviceName,required"`
	// k8s secret name to use when controller is deploying credentials
	// +kubebuilder:validation:MaxLength=5
	// +optional
	NewSecretName string          `json:"newSecretName,omitempty"`
	InitClusterDB InitClusterDB   `json:"initClusterDB,omitempty"`
	Status        DBClusterStatus `json:"status,omitempty"`
	Region        string          `json:"region"`
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

//
//// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
//func (in *DBCluster) DeepCopyInto(out *DBCluster) {
//	*out = *in
//	out.TypeMeta = in.TypeMeta
//	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
//	// in.Status.DeepCopyInto(&out.Status)
//	copier.Copy(&in.ClusterSpec, &out.ClusterSpec)
//	return
//}
