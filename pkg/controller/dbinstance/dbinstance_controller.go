package dbinstance

import (
	"context"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_dbinstance")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDBInstance{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dbinstance-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource DBInstance
	err = c.Watch(&source.Kind{Type: &kubev1alpha1.DBInstance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDBInstance{}

// ReconcileDBInstance reconciles a DBInstance object
type ReconcileDBInstance struct {
	client    client.Client
	scheme    *runtime.Scheme
	rdsClient *rds.RDS
}

// Reconcile reads that state of the cluster for a DBInstance object and makes changes based on the state read
// and what is in the DBInstance.Spec
func (r *ReconcileDBInstance) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Fetch the DBInstance instance
	cr := &kubev1alpha1.DBInstance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// get actionType ( AWLAYS )
	// actionType could switch between create/restore AND delete
	actionType := getActionType(cr)

	// validate create, delete and restore specs are not nil
	if err := validateSpecBasedOnType(cr, actionType); err != nil {
		return reconcile.Result{}, err
	}

	// set up cr fields
	if err := r.setCRDefaultsIfNeeded(cr, actionType); err != nil {
		return reconcile.Result{}, err
	}

	// set up rds client
	r.rdsClient = lib.GetRDSClient()

	// set up finalizers
	currentFinalizers := cr.GetFinalizers()
	anyFinalizersExists := len(currentFinalizers) > 0
	deletionTimeExists := cr.DeletionTimestamp != nil

	// add finalizers
	if !deletionTimeExists && !anyFinalizersExists {
		currentFinalizers = append(currentFinalizers, lib.DBInstanceFinalizer)
		cr.SetFinalizers(currentFinalizers)
		if err := lib.UpdateCr(r.client, cr); err != nil {
			return reconcile.Result{}, err
		}
	}

	// get new instanceObj
	insObj := rdsLib.NewInstance(
		r.rdsClient,
		cr.Spec.CreateInstanceSpec,
		cr.Spec.DeleteInstanceSpec,
		cr.Spec.RestoreInstanceFromSnap,
	)

	err = r.crud(cr, insObj, actionType)
	if err != nil {
		switch err.(type) {
		case *lib.ErrorResourceCreatingInProgress:
			logrus.Warnf("DBInstance not up yet, Reconciling to check again")
			return reconcile.Result{Requeue: true}, nil
		default:
			logrus.Errorf("Namespace: %v | DB Instance ID: %v | Msg: Something went wrong when creating db instance: %v", cr.Namespace, *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier, err)
		}
		return reconcile.Result{}, err
	}

	// create a k8s service with rds endpoint as ExternalName
	if err := r.createExternalNameSvc(cr); err != nil {
		return reconcile.Result{}, err
	}

	// create a k8s secret with DB secrets like username, password, endpoint, etc
	if err := r.createSecret(cr); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
