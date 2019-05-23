package dbcluster

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"

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

	installType := getCreationType(cr)

	// set up finalizers
	currentFinalizers := cr.GetFinalizers()
	anyFinalizersExists := len(currentFinalizers) > 0
	deletionTimeExists := cr.DeletionTimestamp != nil

	// add finalizers
	if !deletionTimeExists && !anyFinalizersExists {
		currentFinalizers = append(currentFinalizers, lib.DBClusterFinalizer)
		cr.SetFinalizers(currentFinalizers)
		if err := lib.UpdateCr(r.client, cr); err != nil {
			return reconcile.Result{}, err
		}
	}

	// delete
	if err := r.handleDelete(cr); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.setUpDefaultsIfNeeded(cr, installType); err != nil {
		return reconcile.Result{}, err
	}

	// create cluster
	if !cr.Status.Created && installType == "newInstall" {
		err := r.createItAndUpdateState(cr)
		if err != nil {
			switch err.(type) {
			case *lib.ErrorResourceCreatingInProgress:
				logrus.Warnf("Namespace: %v | ClusterID: %v | Msg: Cluster still in creating phase. Reconciling to check again.", cr.Namespace, *cr.Spec.DBClusterIdentifier)
				return reconcile.Result{Requeue: true}, nil
			default:
				return reconcile.Result{}, err
			}
		}
	}

	// create secret
	if err := r.createSecret(cr, installType); err != nil {
		return reconcile.Result{}, err
	}

	// restore
	if cr.Status.RestoreNeeded {
		logrus.Infof("Recreate cluster requested for namespace: %v", cr.Namespace)
		// TODO: do type check here and surpress custom errors
		if err := r.restoreClusterFromSnap(cr, installType); err != nil {
			return reconcile.Result{}, err
		}
	}

	// at the end -- keep restore status up to date if RehealFromLatestSnapshot is turned on
	// or if a cr is asking to create a fresh db from an existing snapID
	exists, _ := lib.DbClusterExists(&lib.RDSGenerics{RDSClient: r.rdsClient, ClusterID: *cr.Spec.DBClusterIdentifier})
	if (!exists && cr.Status.Created) || (installType == "newInstallFromSnapshot") && cr.Status.RestoredFromSnapshotName == "" {
		cr.Status.RestoreNeeded = true
		if err := r.updateCrStats(cr); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
}
