package cpeir

import (
	"context"
	"time"
	"strconv"
	"io/ioutil"
	cloudv1alpha1 "github.ibm.com/CASE/cpeir/pkg/apis/cloud/v1alpha1"
	//corev1 "k8s.io/api/core/v1"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/api/resource"
	//"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type System struct {
	Cpu string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
	PV string `yaml:"pv,omitempty"`
	Disk string `yaml:"disk,omitempty"`
}

type CPrequirements struct {
  Requirements map[string]System
}

var log = logf.Log.WithName("controller_cpeir")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CPeir Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCPeir{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cpeir-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CPeir
	err = c.Watch(&source.Kind{Type: &cloudv1alpha1.CPeir{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CPeir
	//err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
	//	IsController: true,
	//	OwnerType:    &cloudv1alpha1.CPeir{},
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}

// blank assignment to verify that ReconcileCPeir implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCPeir{}

// ReconcileCPeir reconciles a CPeir object
type ReconcileCPeir struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CPeir object and makes changes based on the state read
// and what is in the CPeir.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCPeir) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CPeir")

	// Fetch the CPeir instance
	instance := &cloudv1alpha1.CPeir{}
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

	config, err := rest.InClusterConfig()
	if err != nil {
		reqLogger.Info(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		reqLogger.Info(err.Error())
	}

	// Start of the reconcile loop -
	// Read configuration file
	configType := instance.Spec.CPSizeType;
	if configType == "" {
		configType = "default"
	}
	configFile := "/cfgdata/" + instance.Spec.CPType + "-" + instance.Spec.CPVersion +".yaml"
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		reqLogger.Info("yamlFile.Get err ", "Error", err)
  }
	var c CPrequirements
  err = yaml.Unmarshal(yamlFile, &c)
  if err != nil {
		reqLogger.Info("Cannot un marshall file","Error", err)
  }

  reqLogger.Info("CP Requirement", "CPR", c, "configType", configType)
	cpureq, err := resource.ParseQuantity(c.Requirements[configType].Cpu)
	memreq, err := resource.ParseQuantity(c.Requirements[configType].Memory)
	pvreq,  err := resource.ParseQuantity(c.Requirements[configType].PV)
	reqLogger.Info("CPU", "error", err, "cpu", cpureq, "memory", memreq, "Features", instance.Spec.CPFeatures)

	/* Adding requriement calculated from sub-features */
	if len(instance.Spec.CPFeatures) > 0 {
		for _, feature := range instance.Spec.CPFeatures {
			reqLogger.Info("Processing feature", "feature", feature)
			yamlFeatureFile, err := ioutil.ReadFile("/cfgdata/" + feature + "-" + instance.Spec.CPVersion +".yaml")
			if err == nil {
				var cf CPrequirements
				err = yaml.Unmarshal(yamlFeatureFile, &cf)
				if err != nil {
					reqLogger.Info("Cannot un marshall file","Error", err, "File", feature + "-" + instance.Spec.CPVersion)
				}
				reqLogger.Info("CP Feature Requirement", "CPR", cf)
				fcpureq, err := resource.ParseQuantity(cf.Requirements[configType].Cpu)
				if err == nil {
					cpureq.Add(fcpureq)
				}
				fmemreq, err := resource.ParseQuantity(cf.Requirements[configType].Memory)
				if err == nil {
					memreq.Add(fmemreq)
				}
				fpvreq,  err := resource.ParseQuantity(cf.Requirements[configType].PV)
				if err == nil {
					pvreq.Add(fpvreq)
				}
			}
		}
	}

	// Collect nodes information
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: "node-role.kubernetes.io/worker"})
	if err != nil {
		reqLogger.Info(err.Error())
	}
	reqLogger.Info("There are nodes in the cluster\n", "Node number", len(nodes.Items))
	totCpu := resource.NewQuantity(0,resource.DecimalSI);
	totMemory := resource.NewQuantity(0,resource.BinarySI);
  var acpu *resource.Quantity
	var amem *resource.Quantity

	if len(nodes.Items) > 0 {

		for _, node := range nodes.Items {
			acpu = node.Status.Allocatable.Cpu()
			amem = node.Status.Allocatable.Memory()

			totCpu.Add(*acpu)
			totMemory.Add(*amem)
		}
	}

  reqLogger.Info("total available capacity:", "CPU", totCpu.MilliValue(), "Memory", totMemory.Value())

  if (totCpu.Cmp(cpureq)<0) && (totMemory.Cmp(memreq)<0) {
		instance.Status.ClusterStatus = "NotInstallable"
	} else {
		instance.Status.ClusterStatus = "ReadyToInstall"
	}
	instance.Status.StatusMessages = "Allocatable worker nodes capacity is CPU="+strconv.FormatInt(totCpu.MilliValue(),10)+"m and memory="+strconv.FormatInt(totMemory.Value(),10)+"\n"+"Requirement is CPU="+strconv.FormatInt(cpureq.MilliValue(),10)+"m and memory="+strconv.FormatInt(memreq.Value(),10)
	instance.Status.CPReqCPU = cpureq
	instance.Status.CPReqMemory = memreq
	instance.Status.CPReqStorage = pvreq
	instance.Status.ClusterCPU = *totCpu
	instance.Status.ClusterMemory = *totMemory

	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Status update failed")
		return reconcile.Result{}, err
	}

	reqLogger.Info("status","status",instance.Status.ClusterStatus)
	reqLogger.Info("Finished reconcile")
	// Recheck status every minutes - may need to revisit
	return reconcile.Result{RequeueAfter: time.Second*300}, nil
}
