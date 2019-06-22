package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// only use this after DBCluster is available
func getClusterSvc(cr *kubev1alpha1.DBCluster) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.ServiceName,
			Namespace: cr.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type:         v1.ServiceTypeExternalName,
			ExternalName: *cr.Status.DescriberClusterOutput.DBClusters[0].Endpoint,
		},
	}
}
