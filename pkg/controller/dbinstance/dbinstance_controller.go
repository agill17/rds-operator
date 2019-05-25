package dbinstance

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/agill17/rds-operator/pkg/lib"

	"github.com/aws/aws-sdk-go/service/rds"

	goerror "errors"

	kubev1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
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

	// if spec is nil, return err until not nil
	// this is to avoid null pointer dereference ( as i am directly using aws objects )
	if cr.Spec.CreateInstanceSpec == nil {
		logrus.Errorf("createInstanceSpec cannot be nil. Please provide a spec and try again in namespace: %v", cr.Namespace)
		return reconcile.Result{}, goerror.New("EmptyDBInstanceSpecError")
	}

	r.rdsClient = lib.GetRDSClient()

	// set up cr fields when not associated to cluster
	if err := r.setCRDefaultsIfNeeded(cr); err != nil {
		return reconcile.Result{}, err
	}

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

	// delete
	if err := r.handleDeleteEvents(cr, *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier); err != nil {
		return reconcile.Result{}, err
	}

	// create
	if !cr.Status.DeployedInitially {
		if _, err := r.createNewDBInstance(cr); err != nil {
			switch err.(type) {
			case *lib.ErrorResourceCreatingInProgress:
				logrus.Warnf("DBInstance not up yet, Reconciling to check again")
				return reconcile.Result{Requeue: true}, nil
			default:
				logrus.Errorf("Namespace: %v | DB Instance ID: %v | Msg: Something went wrong when creating db instance: %v", cr.Namespace, *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier, err)
				return reconcile.Result{}, err
			}
		}
	}

	// create a k8s service with rds endpoint as ExternalName
	if err := r.createExternalNameSvc(cr); err != nil {
		return reconcile.Result{}, err
	}

	// create a k8s secret with DB secrets like username, password, endpoint, etc
	if err := r.createSecret(cr); err != nil {
		return reconcile.Result{}, err
	}

	// restore
	instanceExists, _ := lib.DBInstanceExists(&lib.RDSGenerics{RDSClient: r.rdsClient, InstanceID: *cr.Spec.CreateInstanceSpec.DBInstanceIdentifier})
	if cr.Status.DeployedInitially && !instanceExists {
		if err := r.restore(cr); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
