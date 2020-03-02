package eventmediator

import (
//	"context"
	"sync"

	eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
//	corev1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/types"
//	"sigs.k8s.io/controller-runtime/pkg/client"
//	"sigs.k8s.io/controller-runtime/pkg/controller"
	// "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
//	"sigs.k8s.io/controller-runtime/pkg/handler"
//	logf "sigs.k8s.io/controller-runtime/pkg/log"
//	"sigs.k8s.io/controller-runtime/pkg/manager"
//	"sigs.k8s.io/controller-runtime/pkg/reconcile"
//	"sigs.k8s.io/controller-runtime/pkg/source"
)


/* Return true if str is an element of array of string */
func stringInArray(stringArray []string, str string) bool {
    for _, stringInArray := range stringArray {
       if stringInArray == str {
            return true
        }
    }
    return false
}

func namespaceNameHash(namespace string, name string ) string {
    return namespace + "/" + name
}

func mediatorHash(mediator *eventsv1alpha1.EventMediator) string {
    return namespaceNameHash(mediator.Namespace, mediator.Name)
}

func mediationsHash(mediations *eventsv1alpha1.EventMediations) string {
    return namespaceNameHash(mediations.Namespace, mediations.Name)
}

/* EventMediationImplManager is responsible for running one instance of mediation.  Each instance is scoped within
  a Mediator.
*/
type EventMediationImplManager struct {
    manager *EventManager // top level manager
    mediator *eventsv1alpha1.EventMediator // the mediator that imports or contains this mediation impl
    mediations *eventsv1alpha1.EventMediations // The meditionas resource that contains this impl. May be null
    mediationImpl  *eventsv1alpha1.EventMediationImpl // the mediation impl to be run
}

func (executor *EventMediationImplManager ) Start() {
}

func (executor *EventMediationImplManager) Stop () {
}


/* Responsible for managing the life cycle of all mediations contained within an Mediations resource*/
type MediationsManager struct {
   manager *EventManager    // top level manager
   mediator *eventsv1alpha1.EventMediator // the mediator that imports the mediations
   mediations *eventsv1alpha1.EventMediations // The mediations resource being imported.
   implManagers map[string]*EventMediationImplManager  // manager for each MediationIMpl
}

func (mediationsManager * MediationsManager) initialize() {
    mediationImpls := mediationsManager.mediations.Spec.Mediations
    for _, oneMediationImpl := range mediationImpls {
        if oneMediationImpl.Mediation != nil {
             hash :=  oneMediationImpl.Mediation.Name
             mediationImplMgr := &EventMediationImplManager {
                                  manager: mediationsManager.manager,
                                  mediator: mediationsManager.mediator,
                                  mediations: mediationsManager.mediations,
                                  mediationImpl:  oneMediationImpl.Mediation,
                            }
             mediationsManager.implManagers[hash] =  mediationImplMgr
             mediationImplMgr.Start()
        }
    }
}


/* Mnages the mediations for one Mediator */
type MediatorManager struct {
    manager *EventManager // top level manager
    mediator *eventsv1alpha1.EventMediator // the mediator whose mediations we are managing
    importMediations map[string]*MediationsManager // imported mediations
    containedEventMediationImplMgr map[string]*EventMediationImplManager // mediations contained within
}

/* Add a new mediations */
func (mediatorMgr *MediatorManager) addMediations(mediations *eventsv1alpha1.EventMediations) {
    mediator := mediatorMgr.mediator
    if mediator.Namespace != mediations.Namespace {
        // ignore if not in the same namespace
        return
    }
    if mediator.Spec.ImportMediations != nil {
        if stringInArray(*mediator.Spec.ImportMediations, mediations.Name) {
            /* mediations now available to use */
            mediationsMgr := &MediationsManager {
                manager: mediatorMgr.manager,
                mediator: mediatorMgr.mediator,
                mediations : mediations,
           }
           hash := mediationsHash(mediations)
           mediatorMgr.importMediations[hash] = mediationsMgr
           mediationsMgr.initialize()
        }
    }
}

/* Add a new mediator */
func (mediatorMgr *MediatorManager) initialize() {

    /* initialize imported mediationss */
    if mediatorMgr.mediator.Spec.ImportMediations != nil {
        for _, importName := range *mediatorMgr.mediator.Spec.ImportMediations {
           hash := namespaceNameHash(mediatorMgr.mediator.Namespace, importName)
           mediations := mediatorMgr.manager.mediations[hash]
           if mediations != nil {
               /* found */
               mediationsMgr := &MediationsManager {
                    manager: mediatorMgr.manager,
                    mediator: mediatorMgr.mediator,
                    mediations : mediations,
               }
               mediatorMgr.importMediations[hash] = mediationsMgr
               mediationsMgr.initialize()
           }
        }
    }

    /* initialize contained mediations */
    if mediatorMgr.mediator.Spec.Mediations != nil {
        mediator := mediatorMgr.mediator
        for _, containedMediationsImpl := range *mediator.Spec.Mediations {
            if containedMediationsImpl.Mediation != nil {
                mediationImplMgr := &EventMediationImplManager {
                                  manager: mediatorMgr.manager,
                                  mediator: mediatorMgr.mediator,
                                  mediations: nil,
                                  mediationImpl:  containedMediationsImpl.Mediation,
                            }
                 mediatorMgr.containedEventMediationImplMgr[containedMediationsImpl.Mediation.Name] = mediationImplMgr
                 mediationImplMgr.Start()
            }
        }
    }
}


type EventManager struct {
    mediatorMgrs map[string] *MediatorManager
    mediations map[string]*eventsv1alpha1.EventMediations // cache of EventMediations objects
    functions map[string]*eventsv1alpha1.EventFunctionImpl
/*
    MediationExecutors *MediationExecutors // mediation executors
    FunctionLibrary *FunctionLibrary // library of functions
*/
    mutex  sync.Mutex // mutex
}


func (mgr *EventManager) GetFunction(name string) (*eventsv1alpha1.EventFunctionImpl, bool) {
    mgr.mutex.Lock()
    defer mgr.mutex.Unlock()
    obj, ok := mgr.functions[name]
    return obj, ok
}

func (mgr *EventManager) addEventMediator(mediator *eventsv1alpha1.EventMediator) {
    mgr.mutex.Lock()
    defer mgr.mutex.Unlock()

    /* add the functions */
    if mediator.Spec.Mediations != nil {
        for _, mediatorImpl := range *mediator.Spec.Mediations {
            if mediatorImpl.Function != nil {
                mgr.functions[mediatorImpl.Function.Name] =  mediatorImpl.Function
            }
        }
    }

    hash := mediatorHash(mediator)
    mediatorMgr := mgr.mediatorMgrs[hash]
    if mediatorMgr != nil {
        /* TODO: Stop all executors for this mediator */
    }
    /* Add new entry */
    mediatorMgr = &MediatorManager {
        manager: mgr,
        mediator: mediator,
    }
    hash = mediatorHash(mediator)
    mgr.mediatorMgrs[hash] = mediatorMgr
    mediatorMgr.initialize()
}

func (mgr *EventManager) addEventMediations(mediations *eventsv1alpha1.EventMediations) {
    mgr.mutex.Lock()
    defer mgr.mutex.Unlock()

    /* TODO: check for updates */
    hash := mediationsHash(mediations)
    mgr.mediations[hash] = mediations

   /* add functions */
    for _, mediationImpl := range mediations.Spec.Mediations  {
        if mediationImpl.Function != nil {
             mgr.functions[mediationImpl.Function.Name] = mediationImpl.Function
        }
    }

    /* refresh all existing mediators */
    for _, mediatorMgr := range mgr.mediatorMgrs {
        mediatorMgr.addMediations(mediations)
    }
}

func (mgr *EventManager) print () {
    mgr.mutex.Lock()
    defer mgr.mutex.Unlock()
}
