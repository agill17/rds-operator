package dbcluster

import (
	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileDBCluster) getSecretObj(cr *kubev1alpha1.DBCluster, installType rdsLib.RDSAction) *v1.Secret {

	secretName := getClusterSecretName(cr)
	s := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: cr.Namespace,
			Labels:    cr.GetLabels(),
		},
		Type: v1.SecretType("Opaque"),
		Data: getSecretData(cr, installType),
	}
	// setup ownerReference
	controllerutil.SetControllerReference(cr, s, r.scheme)

	return s

}

func getSecretData(cr *kubev1alpha1.DBCluster, installType rdsLib.RDSAction) map[string][]byte {
	in := map[string][]byte{
		"DATABASE_ID":             []byte(*cr.Status.DescriberClusterOutput.DBClusters[0].DBClusterIdentifier),
		"DATABASE_NAME":           []byte(*cr.Status.DescriberClusterOutput.DBClusters[0].DatabaseName),
		"ClUSTER_ENDPOINT":        []byte(*cr.Status.DescriberClusterOutput.DBClusters[0].Endpoint),
		"CLUSTER_READER_ENDPOINT": []byte(*cr.Status.DescriberClusterOutput.DBClusters[0].ReaderEndpoint),
	}

	if installType == rdsLib.CREATE {
		in["DATABASE_USERNAME"] = []byte(*cr.Spec.CreateClusterSpec.MasterUsername)
		in["DATABASE_PASSWORD"] = []byte(*cr.Spec.CreateClusterSpec.MasterUserPassword)
	}
	// else if installType == rdsLib.RESTORE {
	// 	// get data from another field if defined I guess?
	// }

	return in
}
