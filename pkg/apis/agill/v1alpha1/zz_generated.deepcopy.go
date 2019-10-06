// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSpec) DeepCopyInto(out *ClusterSpec) {
	*out = *in
	in.CredentialsFrom.DeepCopyInto(&out.CredentialsFrom)
	if in.SnapshotIdentifier != nil {
		in, out := &in.SnapshotIdentifier, &out.SnapshotIdentifier
		*out = new(string)
		**out = **in
	}
	if in.AvailabilityZones != nil {
		in, out := &in.AvailabilityZones, &out.AvailabilityZones
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	if in.BacktrackWindow != nil {
		in, out := &in.BacktrackWindow, &out.BacktrackWindow
		*out = new(int64)
		**out = **in
	}
	if in.BackupRetentionPeriod != nil {
		in, out := &in.BackupRetentionPeriod, &out.BackupRetentionPeriod
		*out = new(int64)
		**out = **in
	}
	if in.CharacterSetName != nil {
		in, out := &in.CharacterSetName, &out.CharacterSetName
		*out = new(string)
		**out = **in
	}
	if in.DBClusterIdentifier != nil {
		in, out := &in.DBClusterIdentifier, &out.DBClusterIdentifier
		*out = new(string)
		**out = **in
	}
	if in.DBClusterParameterGroupName != nil {
		in, out := &in.DBClusterParameterGroupName, &out.DBClusterParameterGroupName
		*out = new(string)
		**out = **in
	}
	if in.DBSubnetGroupName != nil {
		in, out := &in.DBSubnetGroupName, &out.DBSubnetGroupName
		*out = new(string)
		**out = **in
	}
	if in.DatabaseName != nil {
		in, out := &in.DatabaseName, &out.DatabaseName
		*out = new(string)
		**out = **in
	}
	if in.DeletionProtection != nil {
		in, out := &in.DeletionProtection, &out.DeletionProtection
		*out = new(bool)
		**out = **in
	}
	if in.DestinationRegion != nil {
		in, out := &in.DestinationRegion, &out.DestinationRegion
		*out = new(string)
		**out = **in
	}
	if in.EnableCloudwatchLogsExports != nil {
		in, out := &in.EnableCloudwatchLogsExports, &out.EnableCloudwatchLogsExports
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	if in.EnableIAMDatabaseAuthentication != nil {
		in, out := &in.EnableIAMDatabaseAuthentication, &out.EnableIAMDatabaseAuthentication
		*out = new(bool)
		**out = **in
	}
	if in.Engine != nil {
		in, out := &in.Engine, &out.Engine
		*out = new(string)
		**out = **in
	}
	if in.EngineMode != nil {
		in, out := &in.EngineMode, &out.EngineMode
		*out = new(string)
		**out = **in
	}
	if in.EngineVersion != nil {
		in, out := &in.EngineVersion, &out.EngineVersion
		*out = new(string)
		**out = **in
	}
	if in.KmsKeyId != nil {
		in, out := &in.KmsKeyId, &out.KmsKeyId
		*out = new(string)
		**out = **in
	}
	if in.MasterUserPassword != nil {
		in, out := &in.MasterUserPassword, &out.MasterUserPassword
		*out = new(string)
		**out = **in
	}
	if in.MasterUsername != nil {
		in, out := &in.MasterUsername, &out.MasterUsername
		*out = new(string)
		**out = **in
	}
	if in.OptionGroupName != nil {
		in, out := &in.OptionGroupName, &out.OptionGroupName
		*out = new(string)
		**out = **in
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int64)
		**out = **in
	}
	if in.PreSignedUrl != nil {
		in, out := &in.PreSignedUrl, &out.PreSignedUrl
		*out = new(string)
		**out = **in
	}
	if in.PreferredBackupWindow != nil {
		in, out := &in.PreferredBackupWindow, &out.PreferredBackupWindow
		*out = new(string)
		**out = **in
	}
	if in.PreferredMaintenanceWindow != nil {
		in, out := &in.PreferredMaintenanceWindow, &out.PreferredMaintenanceWindow
		*out = new(string)
		**out = **in
	}
	if in.ReplicationSourceIdentifier != nil {
		in, out := &in.ReplicationSourceIdentifier, &out.ReplicationSourceIdentifier
		*out = new(string)
		**out = **in
	}
	if in.SourceRegion != nil {
		in, out := &in.SourceRegion, &out.SourceRegion
		*out = new(string)
		**out = **in
	}
	if in.StorageEncrypted != nil {
		in, out := &in.StorageEncrypted, &out.StorageEncrypted
		*out = new(bool)
		**out = **in
	}
	if in.VpcSecurityGroupIds != nil {
		in, out := &in.VpcSecurityGroupIds, &out.VpcSecurityGroupIds
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSpec.
func (in *ClusterSpec) DeepCopy() *ClusterSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialsFrom) DeepCopyInto(out *CredentialsFrom) {
	*out = *in
	if in.SecretName != nil {
		in, out := &in.SecretName, &out.SecretName
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialsFrom.
func (in *CredentialsFrom) DeepCopy() *CredentialsFrom {
	if in == nil {
		return nil
	}
	out := new(CredentialsFrom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBCluster) DeepCopyInto(out *DBCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.ClusterSpec.DeepCopyInto(&out.ClusterSpec)
	in.InitClusterDB.DeepCopyInto(&out.InitClusterDB)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBCluster.
func (in *DBCluster) DeepCopy() *DBCluster {
	if in == nil {
		return nil
	}
	out := new(DBCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBClusterList) DeepCopyInto(out *DBClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DBCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBClusterList.
func (in *DBClusterList) DeepCopy() *DBClusterList {
	if in == nil {
		return nil
	}
	out := new(DBClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBClusterStatus) DeepCopyInto(out *DBClusterStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBClusterStatus.
func (in *DBClusterStatus) DeepCopy() *DBClusterStatus {
	if in == nil {
		return nil
	}
	out := new(DBClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBInstance.
func (in *DBInstance) DeepCopy() *DBInstance {
	if in == nil {
		return nil
	}
	out := new(DBInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBInstanceList) DeepCopyInto(out *DBInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DBInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBInstanceList.
func (in *DBInstanceList) DeepCopy() *DBInstanceList {
	if in == nil {
		return nil
	}
	out := new(DBInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBInstanceSpec) DeepCopyInto(out *DBInstanceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBInstanceSpec.
func (in *DBInstanceSpec) DeepCopy() *DBInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(DBInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBInstanceStatus) DeepCopyInto(out *DBInstanceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBInstanceStatus.
func (in *DBInstanceStatus) DeepCopy() *DBInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(DBInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBSubnetGroup.
func (in *DBSubnetGroup) DeepCopy() *DBSubnetGroup {
	if in == nil {
		return nil
	}
	out := new(DBSubnetGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBSubnetGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBSubnetGroupList) DeepCopyInto(out *DBSubnetGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DBSubnetGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBSubnetGroupList.
func (in *DBSubnetGroupList) DeepCopy() *DBSubnetGroupList {
	if in == nil {
		return nil
	}
	out := new(DBSubnetGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBSubnetGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBSubnetGroupSpec) DeepCopyInto(out *DBSubnetGroupSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBSubnetGroupSpec.
func (in *DBSubnetGroupSpec) DeepCopy() *DBSubnetGroupSpec {
	if in == nil {
		return nil
	}
	out := new(DBSubnetGroupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBSubnetGroupStatus) DeepCopyInto(out *DBSubnetGroupStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBSubnetGroupStatus.
func (in *DBSubnetGroupStatus) DeepCopy() *DBSubnetGroupStatus {
	if in == nil {
		return nil
	}
	out := new(DBSubnetGroupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitClusterDB) DeepCopyInto(out *InitClusterDB) {
	*out = *in
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitClusterDB.
func (in *InitClusterDB) DeepCopy() *InitClusterDB {
	if in == nil {
		return nil
	}
	out := new(InitClusterDB)
	in.DeepCopyInto(out)
	return out
}
