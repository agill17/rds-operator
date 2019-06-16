package dbcluster

import (
	"context"

	"github.com/agill17/rds-operator/pkg/rdsLib"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/agill17/rds-operator/pkg/lib"

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

var log = logf.Log.WithName("controller_dbcluster")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDBCluster{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dbcluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &kubev1alpha1.DBCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDBCluster{}

// ReconcileDBCluster reconciles a DBCluster object
type ReconcileDBCluster struct {
	client    client.Client
	scheme    *runtime.Scheme
	rdsClient *rds.RDS
}

func (r *ReconcileDBCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// set up rds client
	if r.rdsClient == nil {
		r.rdsClient = lib.GetRDSClient()
	}

	// Fetch the DBCluster cr
	cr := &kubev1alpha1.DBCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if err := validateRequiredInput(cr); err != nil {
		return reconcile.Result{}, err
	}

	// get action type ( ALWAYS )
	// actionType changes between create/restore and delete
	// actionType will be delete when its time to delete rds cluster
	actionType := getActionType(cr)

	if err := r.setUpDefaultsIfNeeded(cr, actionType); err != nil {
		return reconcile.Result{}, err
	}

	// set up finalizers
	currentFinalizers := cr.GetFinalizers()
	zeroFinalizers := len(currentFinalizers) == 0
	deletionTimeExists := cr.DeletionTimestamp != nil

	// add finalizers
	if !deletionTimeExists && zeroFinalizers {
		currentFinalizers = append(currentFinalizers, lib.DBClusterFinalizer)
		cr.SetFinalizers(currentFinalizers)
		if err := lib.UpdateCr(r.client, cr); err != nil {
			return reconcile.Result{}, err
		}
	}

	// returns cluster struct which is also part of rds interface
	// so we can call all funcs that are part of the interface as long as cluster satifies the interface
	// cluster obj implements and satifies RDS interface by implementing all methods of that interface
	clusterObj := rdsLib.NewCluster(r.rdsClient, cr.Spec.CreateClusterSpec,
		cr.Spec.DeleteSpec, cr.Spec.CreateClusterFromSnapshot)

	if err := r.crud(cr, clusterObj, actionType); err != nil {
		switch err.(type) {
		case *lib.ErrorResourceCreatingInProgress:
			logrus.Errorf("Namespace: %v | CR: %v | Msg: Cluster still in creating phase. Reconciling to check again.", cr.Namespace, cr.Name)
			return reconcile.Result{Requeue: true}, nil
		}
		return reconcile.Result{}, err
	}

	// create/update secret
	if err := r.createSecret(cr, actionType); err != nil && !errors.IsForbidden(err) {
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: true}, nil
}
