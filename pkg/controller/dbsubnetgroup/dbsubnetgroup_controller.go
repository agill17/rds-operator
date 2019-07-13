package dbsubnetgroup

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/agill17/rds-operator/pkg/lib"

	"github.com/aws/aws-sdk-go/service/rds"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
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

	return nil
}

var _ reconcile.Reconciler = &ReconcileDBSubnetGroup{}

// ReconcileDBSubnetGroup reconciles a DBSubnetGroup object
type ReconcileDBSubnetGroup struct {
	rdsClient *rds.RDS
	client    client.Client
	scheme    *runtime.Scheme
}

func (r *ReconcileDBSubnetGroup) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	if r.rdsClient == nil {
		r.rdsClient = lib.GetRDSClient()
	}

	// Fetch the DBSubnetGroup cr
	cr := &agillv1alpha1.DBSubnetGroup{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if cr.Spec == nil {
		logrus.Errorf("ERROR createDBSubnetGroupSpec cannot be empty in namespace: %v", cr.Namespace)
		return reconcile.Result{Requeue: true}, nil
	}

	// update Missing cr fields
	if err := r.setCRDefaultsIfNeeded(cr); err != nil {
		return reconcile.Result{}, err
	}

	// get finalizers
	deletionTimeStampExists := cr.GetDeletionTimestamp() != nil
	currentFinalizers := cr.GetFinalizers()
	anyFinalizersExists := len(currentFinalizers) > 0

	// set finalizers if needed
	if !deletionTimeStampExists && !anyFinalizersExists {
		if err := lib.AddFinalizer(cr, r.client, lib.DBSubnetGroupFinalizer); err != nil {
			return reconcile.Result{}, err
		}
	}

	// delete -- empty out finalizers -- do not reque
	if deletionTimeStampExists && anyFinalizersExists {
		err := r.deleteSubnetGroup(*cr.Spec.DBSubnetGroupName)
		if err != nil {
			return reconcile.Result{}, err
		}
		cr.SetFinalizers([]string{})
		err = lib.UpdateCr(r.client, cr)
		if err != nil {
			return reconcile.Result{}, err
		}
		logrus.Infof("Successfully deleted DBSubnetGroup %v for namespace: %v", cr.Name, cr.Namespace)
		return reconcile.Result{}, nil
	}

	existsInAws, _ := lib.DBSubnetGroupExists(&lib.RDSGenerics{RDSClient: r.rdsClient, SubnetGroupName: *cr.Spec.DBSubnetGroupName})
	statusMarkedDeployed := cr.Status.Created
	if (!statusMarkedDeployed && !existsInAws) || (!existsInAws && statusMarkedDeployed) {
		logrus.Infof("Creating DBSubnetGroup %v for namespace: %v", cr.Name, cr.Namespace)
		err = r.createSubnetGroup(cr)
		if err != nil {
			return reconcile.Result{}, err
		}
		logrus.Infof("Successfully created DBSubnetGroup %v for namespace: %v", cr.Name, cr.Namespace)
	}

	return reconcile.Result{}, nil
}
