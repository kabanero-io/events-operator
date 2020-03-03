package eventmediator

import (
	"context"

    // routev1 "github.com/openshift/api/route/v1"
	eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
	corev1 "k8s.io/api/core/v1"
    appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
    logf "sigs.k8s.io/controller-runtime/pkg/log"
    "github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
    "os"
    //"strconv"
)

const (
    MEDIATOR_NAME = "MEDIATOR_NAME" // environment variable. If not set, we're running as operator.
)

/* This controller can serve as either an operator, or a regular controller.
  As an operator, it manages Deployments for each EventMEdiator crd instance
  As a controller, the environment variable MEDIATOR_NAME is set to the name of the mediator, and it is responsible for
 updates on the EventMediator CRD instance.
*/
var MediatorName string
var isOperator bool = false

func init () {
    MediatorName :=  os.Getenv(MEDIATOR_NAME)
    isOperator = (MediatorName == "")
}

var log = logf.Log.WithName("controller_eventmediator")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new EventMediator Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileEventMediator{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("eventmediator-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource EventMediator
	err = c.Watch(&source.Kind{Type: &eventsv1alpha1.EventMediator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

    // Watch for deployments
    if isOperator {
        err = c.Watch(
        &source.Kind{Type: &appsv1.Deployment{}},
        &handler.EnqueueRequestForOwner{
            IsController: true,
            OwnerType:    &eventsv1alpha1.EventMediator{}},
        )
	    if err != nil {
		    return err
        }

        err = c.Watch(
        &source.Kind{Type: &corev1.Service{}},
        &handler.EnqueueRequestForOwner{
            IsController: true,
            OwnerType:    &eventsv1alpha1.EventMediator{}},
        )
	    if err != nil {
		    return err
        }
    } 

	return nil
}

// blank assignment to verify that ReconcileEventMediator implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileEventMediator{}

// ReconcileEventMediator reconciles a EventMediator object
type ReconcileEventMediator struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a EventMediator object and makes changes based on the state read
// and what is in the EventMediator.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileEventMediator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling EventMediator")

	// Fetch the EventMediator instance
	instance := &eventsv1alpha1.EventMediator{}
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

    if isOperator {
        result, err := r.reconcileDeployment(request, instance, reqLogger)
        if err != nil  || result.Requeue{
            return result, err
        }

        result, err = r.reconcileService(request, instance, reqLogger)
        return result, err
    } else {
        /* plain controller for one mediator */
        if instance.ObjectMeta.Name ==  MEDIATOR_NAME {
            /* TODO: We should handle this */
        }
    }

	return reconcile.Result{}, nil
}

/* Reconcile deployment for an operator */
func (r *ReconcileEventMediator) reconcileOperator(request reconcile.Request, mediator *eventsv1alpha1.EventMediator, reqLogger logr.Logger) (reconcile.Result, error) {
    result, err := r.reconcileDeployment(request, mediator, reqLogger)
    if err != nil {
        return result, err
    }
    result, err = r.reconcileService(request, mediator, reqLogger)
    if err != nil {
        return result, err
    }
	return reconcile.Result{}, nil
}

/* reconcile Operator */
func (r *ReconcileEventMediator) reconcileDeployment(request reconcile.Request, instance *eventsv1alpha1.EventMediator,  reqLogger logr.Logger) (reconcile.Result, error) {
    /* Check if the deployment already exists, if not create a new one */
    deployment := &appsv1.Deployment{}
    err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
    if err != nil && errors.IsNotFound(err) {
        // Define a new deployment
        dep := r.deploymentForEventMediator(instance)
        reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
        err = r.client.Create(context.TODO(), dep)
        if err != nil {
            reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
            return reconcile.Result{}, err
        }
        // Deployment created successfully - return and requeue
        return reconcile.Result{Requeue: true}, nil
    } else if err != nil {
        reqLogger.Error(err, "Failed to get Deployment")
        return reconcile.Result{}, err
    }

    if portChangedForDeployment(deployment, instance)  {
        deployment.Spec.Template.Spec.Containers[0].Ports = generateDeploymentPorts(instance)
        err = r.client.Update(context.TODO(), deployment)
        if err != nil {
           reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
            return reconcile.Result{}, err
         }
        // Spec updated - return and requeue
        return reconcile.Result{Requeue: true}, nil
    }

    return reconcile.Result{}, nil
}

/* reconcile Operator */
func (r *ReconcileEventMediator) reconcileService(request reconcile.Request, instance *eventsv1alpha1.EventMediator, reqLogger logr.Logger) (reconcile.Result, error) {
    service := &corev1.Service{}
    err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, service)
    if err != nil && errors.IsNotFound(err) {
        // Define a new service
        serv := r.serviceForEventMediator(instance)
        reqLogger.Info("Creating a new Service", "Service.Namespace", serv.Namespace, "Service.Name", serv.Name)
        err = r.client.Create(context.TODO(), service)
        if err != nil {
            reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", serv.Namespace, "Service.Name", serv.Name)
            return reconcile.Result{}, err
        }
        // Servicecreated successfully - return and requeue
        return reconcile.Result{Requeue: true}, nil
    } else if err != nil {
        reqLogger.Error(err, "Failed to get Service")
        return reconcile.Result{}, err
    }

    if portChangedForService(service, instance)  {
        service.Spec.Ports = generateServicePorts(instance)
        err = r.client.Update(context.TODO(), service)
        if err != nil {
           reqLogger.Error(err, "Failed to update Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
            return reconcile.Result{}, err
         }
        // Spec updated - return and requeue
        return reconcile.Result{Requeue: true}, nil
    }
    return reconcile.Result{}, nil
}

/* Return true if the ports in a Deployment have changed */
func portChangedForDeployment(deployment *appsv1.Deployment, mediator *eventsv1alpha1.EventMediator) bool {

    ports := deployment.Spec.Template.Spec.Containers[0].Ports
    listener := mediator.Spec.Listener

    check := make(map[int32] int32)
    for _, portInfo := range ports {
        check[portInfo.ContainerPort] = portInfo.ContainerPort
    }

    numMediatorPorts := 0
    if listener != nil {
       if listener.HttpsPort != 0 {
           numMediatorPorts++
           if   _, exists:= check[int32(listener.HttpsPort)]; ! exists {
               return true
           }
        }
    }
    if len(ports) != numMediatorPorts {
         return true
    }

    return false
}

func generateDeploymentPorts(mediator *eventsv1alpha1.EventMediator) []corev1.ContainerPort {
    var ports []corev1.ContainerPort = make([]corev1.ContainerPort, 0);
    if mediator.Spec.Listener != nil {
        listener := mediator.Spec.Listener
        var port int32
        if listener.HttpsPort != 0 {
             port = int32(listener.HttpsPort)
             ports = append(ports, corev1.ContainerPort {
                    ContainerPort:  port,
                    Name:          "httpsPort",
               } )
        }
    }
    return ports
}

// Return a deployment object
func (r *ReconcileEventMediator) deploymentForEventMediator(mediator *eventsv1alpha1.EventMediator) *appsv1.Deployment {
    ls := labelsForEventMediator(mediator.Name)
    var replicas int32 = 1
    env  := []corev1.EnvVar {
        {
             Name: MEDIATOR_NAME,
             Value: mediator.Name,
        },
    }
    ports := generateDeploymentPorts(mediator)

    dep := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      mediator.Name,
            Namespace: mediator.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &replicas,
            Selector: &metav1.LabelSelector{
                MatchLabels: ls,
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: ls,
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{{
                        Image:   "mchengdocker/kabanero-events:0.2",
                        Name:    "evnetmediator",
                        Command: []string{"entrypoint"},
                        Ports: ports,
                        Env: env,
                      }},
                },
            },
        },
    }

    // Set owner and controller
    controllerutil.SetControllerReference(mediator, dep, r.scheme)
    return dep
}

func labelsForEventMediator(name string) map[string]string {
    return map[string]string{"app": name, "eventmediator_cr": name}
}

func generateServicePorts(mediator *eventsv1alpha1.EventMediator) []corev1.ServicePort {
    ports := make([]corev1.ServicePort, 0)
    if mediator.Spec.Listener != nil {
        listener := mediator.Spec.Listener 
        var port int32
        if listener.HttpsPort != 0 {
             port = int32(listener.HttpsPort)
             ports = append(ports, corev1.ServicePort {
                    Port:  port,
               } )
        }
    }
    return ports
}

// Return a Service object
func (r *ReconcileEventMediator) serviceForEventMediator(mediator *eventsv1alpha1.EventMediator) *corev1.Service {
    ls := labelsForEventMediator(mediator.Name)
    servicePorts := generateServicePorts(mediator)

    service := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      mediator.Name,
            Namespace: mediator.Namespace,
        },
        Spec: corev1.ServiceSpec {
            Ports: servicePorts,
            Selector: ls,
            Type: corev1. ServiceTypeClusterIP,
        },
    }

    // Set owner and controller
    controllerutil.SetControllerReference(mediator, service, r.scheme)
    return service
}

/* Return true if the ports in a Service have changed */
func portChangedForService(service *corev1.Service, mediator *eventsv1alpha1.EventMediator) bool {

    ports := service.Spec.Ports
    listener := mediator.Spec.Listener

    check := make(map[int32] int32)
    for _, portInfo := range ports {
        check[portInfo.Port] = portInfo.Port
    }

    numMediatorPorts := 0
    if listener != nil {
       if listener.HttpsPort != 0 {
           numMediatorPorts++
           if   _, exists:= check[int32(listener.HttpsPort)]; ! exists {
               return true
           }
        }
    }
    if len(ports) != numMediatorPorts {
         return true
    }
    return false
}
