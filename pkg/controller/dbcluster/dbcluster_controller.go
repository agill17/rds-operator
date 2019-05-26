package dbcluster

import (
	"context"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"

	"github.com/agill17/rds-operator/pkg/lib"
	"github.com/agill17/rds-operator/pkg/lib/dbHelpers"

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

	// get install type ( ALWAYS )
	installType := getInstallType(cr)

	if err := r.setUpDefaultsIfNeeded(cr, installType); err != nil {
		return reconcile.Result{}, err
	}

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

	clusterObj := dbHelpers.Cluster{
		RDSClient:            r.rdsClient,
		CreateInput:          cr.Spec.CreateClusterSpec,
		DeleteInput:          cr.Spec.DeleteSpec,
		RestoreFromSnapInput: cr.Spec.CreateClusterFromSnapshot,
	}

	if installType == dbHelpers.DELETE {
		// delete
		if deletionTimeExists && anyFinalizersExists {
			err := dbHelpers.InstallRestoreDelete(&clusterObj, installType)
			if err != nil {
				return reconcile.Result{}, err
			}
			cr.SetFinalizers([]string{})
			if err := lib.UpdateCr(r.client, cr); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
	} else if installType == dbHelpers.CREATE {
		// create cluster
		if !cr.Status.Created {
			err := r.createItAndUpdateState(cr, &clusterObj)
			if err != nil {
				switch err.(type) {
				case *lib.ErrorResourceCreatingInProgress:
					logrus.Warnf("Namespace: %v | CR: %v | Msg: Cluster still in creating phase. Reconciling to check again.", cr.Namespace, cr.Name)
					return reconcile.Result{Requeue: true}, nil
				default:
					return reconcile.Result{}, err
				}
			}
		}
	} else if installType == dbHelpers.RESTORE {
		// create from snapshot
		logrus.Infof("Recreate cluster requested for namespace: %v", cr.Namespace)
		if err := r.restoreAndUpdateState(cr, &clusterObj); err != nil {
			return reconcile.Result{}, err
		}
	}

	// create secret
	if err := r.createSecret(cr); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: true}, nil
}
