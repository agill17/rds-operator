package dbcluster

import (
	"context"

	"github.com/agill17/rds-operator/pkg/rdsLib"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	defaultUsername = "admin"
	defaultPassword = "password"
	defaultUserKey  = "DB_USER"
	defaultPassKey  = "DB_PASS"
)

/*
	- Create Secret
		- when useCredentialsFrom is not defined
			OR
		- when username, password defined is within CreateClusterSpec
			OR
		- when username, password not defined at all ( assume defaults )
	- Do not Secret
		- when useCredentialsFrom is defined

*/

func (r *ReconcileDBCluster) reconcileSecret(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) error {
	var err error
	createSecret, secretName, userKey, passKey := shouldCreateSecret(cr)
	if createSecret {
		_, err = r.createUpdateSecret(cr, rdsAction)
		if err != nil {
			return err
		}
	}

	return setSecretStatusInCR(cr, r.client, secretName, userKey, passKey)
}

func (r *ReconcileDBCluster) createUpdateSecret(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) (*v1.Secret, error) {

	secretObj := getSecretObj(cr, rdsAction)
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, secretObj, func(runtime.Object) error {
		controllerutil.SetControllerReference(cr, secretObj, r.scheme)

		return nil
	})
	if err != nil {
		logrus.Errorf("Failed to create/update secret for namespace: %v -- %v", cr.Namespace, err)
		return nil, err
	}
	return secretObj, nil
}

// only used when createSecret is true
func getSecretObj(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) *v1.Secret {
	data := getSecretData(cr, rdsAction)
	secretName := getClusterSecretName(cr.Name)

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
		Data: data,
	}

	return s

}

// only used when createSecret is true
func getSecretData(cr *kubev1alpha1.DBCluster, rdsAction rdsLib.RDSAction) map[string][]byte {
	var user, pass string = defaultUsername, defaultPassword

	if rdsAction == rdsLib.CREATE && cr.Spec.CreateClusterSpec.MasterUsername != nil {
		user = *cr.Spec.CreateClusterSpec.MasterUsername
	}

	if rdsAction == rdsLib.CREATE && cr.Spec.CreateClusterSpec.MasterUserPassword != nil {
		pass = *cr.Spec.CreateClusterSpec.MasterUserPassword
	}

	in := map[string][]byte{
		defaultUserKey: []byte(user),
		defaultPassKey: []byte(pass),
	}

	return in
}

// returns shouldCreate(bool), secretName, userKey, passKey
func shouldCreateSecret(cr *kubev1alpha1.DBCluster) (bool, string, string, string) {

	// when user provides a secret to get credentials from
	// dont create secret
	if useCredentialsFrom(cr) {
		return false, cr.Spec.CredentialsFrom.SecretName.Name, cr.Spec.CredentialsFrom.UsernameKey,
			cr.Spec.CredentialsFrom.PasswordKey
	}

	// or
	// when user provides username & password using cr.Spec.CreateClusterSpec,
	// create secret
	// or
	// when using defaultusername and defaultPassword
	return true, getClusterSecretName(cr.Name), defaultUserKey, defaultPassKey
}

func setSecretStatusInCR(cr *kubev1alpha1.DBCluster, client client.Client, secretName, userKey, passKey string) error {

	if err := setSecretNameInStatus(secretName, cr, client); err != nil {
		return err
	}
	if err := setUserKeyInStatus(userKey, cr, client); err != nil {
		return err
	}

	return setPassKeyInStatus(passKey, cr, client)
}

func setUserKeyInStatus(userKey string, cr *kubev1alpha1.DBCluster, client client.Client) error {
	if cr.Status.UsernameKey != userKey {
		cr.Status.UsernameKey = userKey
		return lib.UpdateCrStatus(client, cr)
	}
	return nil
}

func setPassKeyInStatus(passKey string, cr *kubev1alpha1.DBCluster, client client.Client) error {
	if cr.Status.PasswordKey != passKey {
		cr.Status.PasswordKey = passKey
		return lib.UpdateCrStatus(client, cr)
	}
	return nil
}

func setSecretNameInStatus(secretName string, cr *kubev1alpha1.DBCluster, client client.Client) error {
	if cr.Status.SecretName != secretName {
		cr.Status.SecretName = secretName
		return lib.UpdateCrStatus(client, cr)
	}
	return nil
}
