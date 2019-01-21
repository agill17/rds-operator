package dbinstance

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	agillv1alpha1 "github.com/agill17/rds-operator/pkg/apis/agill/v1alpha1"
	"github.com/agill17/rds-operator/pkg/controller/lib"
	"github.com/aws/aws-sdk-go/service/rds"
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

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new DBInstance Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
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
	err = c.Watch(&source.Kind{Type: &agillv1alpha1.DBInstance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner DBInstance

	return nil
}

var _ reconcile.Reconciler = &ReconcileDBInstance{}

// ReconcileDBInstance reconciles a DBInstance object
type ReconcileDBInstance struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	scheme    *runtime.Scheme
	rdsClient *rds.RDS
}

// Reconcile reads that state of the cluster for a DBInstance object and makes changes based on the state read
// and what is in the DBInstance.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDBInstance) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	// reqLogger.Info("Reconciling DBInstance")

	// Fetch the DBInstance instance
	instance := &agillv1alpha1.DBInstance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set some vars
	r.rdsClient = lib.GetRdsClient()
	ns := instance.Namespace
	crName := instance.Name
	dbID := lib.SetDBID(ns, crName)
	user, pass := lib.GetUsernamePassword(instance)

	// get finalizers
	deletionTimeExists := instance.GetDeletionTimestamp() != nil
	currentFinalizers := instance.GetFinalizers()
	anyFinalizersExists := len(currentFinalizers) > 0

	// set finalizers ( so controller can cleanup during delete events )
	if !anyFinalizersExists && !deletionTimeExists {
		finalizersToAdd := append(currentFinalizers, lib.DBInstanceFinalizer)
		instance.SetFinalizers(finalizersToAdd)
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// delete dbInstance -- if successful -- empty out the finalizers -- do not reque
	if deletionTimeExists && anyFinalizersExists {
		logrus.Warnf("reconcile: time to delete %s", instance.Name)
		err := r.deleteDBInstance(instance, dbID)
		if err != nil {
			return reconcile.Result{}, err
		}
		instance.SetFinalizers([]string{})
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		logrus.Infof("Current finalizers after cleanup: %v", instance.GetFinalizers())
		return reconcile.Result{}, nil
	}

	// create
	if err := r.createInstanceAndUpdateState(instance, dbID, ns, request); err != nil {
		logrus.Errorf("Namespace: %v | DB Instance ID: %v | Msg: Something went wrong when creating db instance: %v", instance.Namespace, dbID, err)
		return reconcile.Result{}, err
	} else if instance.Status.DeployedInitially {
		if err := r.createInitDBJob(instance, request); err != nil {
			logrus.Errorf("Namespace: %v | DB Instance ID: %v | Msg: Something went wrong when creating init-db job: %v", instance.Namespace, dbID, err)
			return reconcile.Result{}, err
		}
		r.createExternalNameSvc(instance, dbID, request)
		r.createSecret(instance, dbID, instance.Spec.DBName, user, pass, request)
	}

	return reconcile.Result{Requeue: true, RequeueAfter: 1 * time.Second}, nil
}

func (r *ReconcileDBInstance) getCrInstance(request reconcile.Request) (*agillv1alpha1.DBInstance, error) {
	cr := &agillv1alpha1.DBInstance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, err
		}
	}
	return cr, err

}

/*
	Cases to handle:
	1. DB exists in AWS AND CR has no status about it -- dont do anything but just log that msg
	2. DB does not exist in AWS AND CR has no status about any deployment -- create a new fresh DB
	3. DB does not exist in AWS AND CR has status that it was deployed atleast once -- reheal from snapshot if defined to do so, notify
	4. DB does exists in AWS AND CR has status about it -- dont do anything but just log that msg
*/
func (r *ReconcileDBInstance) createInstanceAndUpdateState(cr *agillv1alpha1.DBInstance, dbID, ns string, request reconcile.Request) error {
	var err error

	dbInput := r.createDBInstanceInput(cr, dbID)
	dbExistsInAws, _ := r.dbInstanceExists(dbID)
	crHasStatus := crHasDBStatus(cr)

	if dbExistsInAws && !crHasStatus {
		// 1.
		logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: dbID already exists in AWS! Please create a new DB by updating the name under metadata of the cr", cr.Namespace, dbID)
	} else if !dbExistsInAws && !crHasStatus {
		// 2.
		if _, err := r.createNewDBInstance(cr, dbID, dbInput, request); err != nil {
			return err
		}
	} else if !dbExistsInAws && crHasStatus {
		// 3.
		logrus.Errorf("Namespace: %v | DB Identifier: %v | Msg: CR is marked as db is deployed but does not exist in AWS!!!", cr.Namespace, dbID)
		if err := r.restoreFromSnapshot(cr, dbID, request); err != nil {
			return err
		}
	} else if dbExistsInAws && crHasStatus {
		// 4.
		// logrus.Infof("Namespace: %v | DB Identifier: %v | Msg: CR has DB marked as deployed and also exists in aws", cr.Namespace, dbID)
	}
	return err
}

// Contains verifies if a list of strings contains a given string
func Contains(l []string, s string) bool {
	for _, elem := range l {
		if elem == s {
			return true
		}
	}
	return false
}
