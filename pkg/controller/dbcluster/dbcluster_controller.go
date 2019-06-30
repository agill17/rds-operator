package dbcluster

import (
	"context"

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

	// create/update secret
	// this also updates CR status with what secret to use and what userKey and passKey to use
	// whether that key is coming from a user provided secret to use or from internally
	if err := r.reconcileSecret(cr, actionType); err != nil && !errors.IsForbidden(err) {
		return reconcile.Result{}, err
	}

	if err := r.crud(cr, actionType); err != nil {
		switch err.(type) {
		case *lib.ErrorResourceCreatingInProgress, *lib.ErrorKubernetesSecretDoesNotExist:
			logrus.Warnf("Namespace: %v | CR: %v | Msg: %v", cr.Namespace, cr.Name, err)
			return reconcile.Result{Requeue: true}, nil
		case *lib.ErrorKubernetesSecretGettingDeleted:
			logrus.Warnf("Namespace: %v | CR: %v | Msg: K8S objects are in deleting phase. Stopping reconcile..", cr.Namespace, cr.Name)
			return reconcile.Result{Requeue: false}, nil
		}
		return reconcile.Result{}, err
	}

	if err := reconcileInitDBJob(cr, r.client, r.rdsClient); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.createExternalSvc(cr); err != nil && !errors.IsForbidden(err) {
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: true}, nil
}
