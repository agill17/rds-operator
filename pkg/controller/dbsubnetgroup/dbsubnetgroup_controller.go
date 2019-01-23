package dbsubnetgroup

import (
	"context"

	"github.com/agill17/rds-operator/pkg/controller/lib"

	"github.com/aws/aws-sdk-go/service/rds"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	corev1 "k8s.io/api/core/v1"
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

var log = logf.Log.WithName("controller_dbsubnetgroup")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new DBSubnetGroup Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDBSubnetGroup{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dbsubnetgroup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource DBSubnetGroup
	err = c.Watch(&source.Kind{Type: &agillv1alpha1.DBSubnetGroup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner DBSubnetGroup
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &agillv1alpha1.DBSubnetGroup{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDBSubnetGroup{}

// ReconcileDBSubnetGroup reconciles a DBSubnetGroup object
type ReconcileDBSubnetGroup struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	scheme    *runtime.Scheme
	rdsClient *rds.RDS
}

// Reconcile reads that state of the cluster for a DBSubnetGroup object and makes changes based on the state read
// and what is in the DBSubnetGroup.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDBSubnetGroup) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	r.rdsClient = lib.GetRdsClient()

	// Fetch the DBSubnetGroup instance
	instance := &agillv1alpha1.DBSubnetGroup{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// get finalizers
	deletionTimeStampExists := instance.GetDeletionTimestamp() != nil
	currentFinalizers := instance.GetFinalizers()
	anyFinalizersExists := len(currentFinalizers) > 0

	// set finalizers if needed
	if !deletionTimeStampExists && !anyFinalizersExists {
		currentFinalizers = append(currentFinalizers, "")
		instance.SetFinalizers(currentFinalizers)
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err // try again
		}
	}

	// delete -- empty out finalizers -- do not reque
	if deletionTimeStampExists && anyFinalizersExists {
		err := r.deleteSubnetGroup(instance.Name)
		if err != nil {
			return reconcile.Result{}, err
		}
		instance.SetFinalizers([]string{})
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	existsInAws, _ := r.subnetGroupExists(instance.Name)
	statusMarkedDeployed := instance.Status.Created
	// if status is NOT deplyed && does NOT exists in AWS
	//   OR
	// if does not exists in AWS && status is marked as deployed
	// create
	if (!statusMarkedDeployed && !existsInAws) || (!existsInAws && statusMarkedDeployed) {
		err = r.createSubnetGroup(request)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
