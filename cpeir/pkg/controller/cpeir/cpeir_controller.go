package cpeir

import (
	"context"
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/resource"
	"gopkg.in/yaml.v2"
	//"strconv"
	cloudv1alpha1 "github.ibm.com/CASE/cpeir/pkg/apis/cloud/v1alpha1"
	//appsv1 "k8s.io/api/apps/v1"
	//corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

type capacity struct {
	TotCpu resource.Quantity `json:totcpu`
	TotMem resource.Quantity `json:totmem`
	MaxCpu resource.Quantity `json:maxcpu`
	MaxMem resource.Quantity `json:maxmem`
	Arch   string `json:arch`
	Kubelet string `json:kubever`
	NumNode int   `json:numnode`
}

type version struct {
	OCPVersion string `json:version`
	UpgradeChannel string `json:channel`
}

type registry struct {
	Configured bool `json:configured`
	External bool `json:external`
	Capacity resource.Quantity `json:capacity`
}

type installed struct {
	Name string `json:name,omitempty`
	Installed bool `json:installed`
	Version string `json:version,omitempty`
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

func connected() (ok bool) {
    _, err := http.Get("https://cp.icr.io")
    if err != nil {
        return false
    }
    return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
	statusMessage := ""
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

	// Read configuration file - calculate installation requirements
	configType := instance.Spec.CPSizeType;
	if configType == "" {
		configType = "default"
		instance.Spec.CPSizeType = "default"
	}

	if instance.Status.CPStatus == "" {
		instance.Status.CPStatus = "Initial"
	}

	if instance.Spec.Action == "" {
		instance.Spec.Action = "Check"
	}

	// Collect cluster information

	var capjson capacity
	var regjson registry
	var verjson version

  rcapacity, err := r.getRest("capacity","", "")
	json.Unmarshal(rcapacity, &capjson)
	reqLogger.Info(string(rcapacity),"json",capjson)

	rregistry, err := r.getRest("registry","", "")
	json.Unmarshal(rregistry,&regjson)
	reqLogger.Info(string(rregistry),"json",regjson)

	rversion, err := r.getRest("version", "", "")
	json.Unmarshal(rversion,&verjson)
	reqLogger.Info(string(rversion),"json",verjson)

	instance.Status.ClusterStatus.ClusterCPU = capjson.TotCpu
	instance.Status.ClusterStatus.ClusterMemory = capjson.TotMem
	instance.Status.ClusterStatus.ClusterArch = capjson.Arch
	instance.Status.ClusterStatus.ClusterWorkerNum = capjson.NumNode
	instance.Status.ClusterStatus.ClusterKubelet = capjson.Kubelet
	instance.Status.ClusterStatus.OCPVersion = verjson.OCPVersion
	instance.Status.OnlineInstall = connected()
	// Set to false for now - compilation problem
	instance.Status.OfflineInstall = regjson.External

	// Perform check for status < Installed

  numfeat := 0
  installedfeat := 0
	var cpureq resource.Quantity;
	var memreq resource.Quantity;
	var pvreq  resource.Quantity;

	if ((instance.Status.CPStatus == "Initial") ||
	    (instance.Status.CPStatus == "ReadyToInstall") ||
			(instance.Status.CPStatus == "NotInstallable" )) {
		// check the status and all installed
		rcheck, err := r.getRest("check",instance.Spec.CPType + "-" + instance.Spec.CPVersion,instance.Name)
		if err != nil {
			reqLogger.Error(err, "Check result error "+instance.Spec.CPType + "-" + instance.Spec.CPVersion)
		}
		var instjson installed
		json.Unmarshal(rcheck, &instjson)
		reqLogger.Info(string(rcheck)," Prod check json",instjson)
		numfeat++

		if !(instjson.Installed) {
			configFile := "/cfgdata/" + instance.Spec.CPType + "-" + instance.Spec.CPVersion +".yaml"
			yamlFile, err := ioutil.ReadFile(configFile)
			if err != nil {
				reqLogger.Error(err, "yamlFile read error "+instance.Spec.CPType + "-" + instance.Spec.CPVersion +".yaml")
				instance.Status.CPStatus = "ValidationFailed"
				statusMessage = statusMessage + "cloudpak unsupported yet "+instance.Spec.CPType + "-" + instance.Spec.CPVersion + "\n"
			}
			var c CPrequirements
			err = yaml.Unmarshal(yamlFile, &c)
			if err != nil {
				reqLogger.Error(err,"Cannot un marshall file")
			}
			cpureq, err = resource.ParseQuantity(c.Requirements[configType].Cpu)
			memreq, err = resource.ParseQuantity(c.Requirements[configType].Memory)
			pvreq,  err = resource.ParseQuantity(c.Requirements[configType].PV)
		} else {
			installedfeat++
			cpureq, err = resource.ParseQuantity("0")
			memreq, err = resource.ParseQuantity("0")
			pvreq,  err = resource.ParseQuantity("0")
		}

		/* Adding requriement calculated from sub-features */
		if len(instance.Spec.CPFeatures) > 0 {
			for _, feature := range instance.Spec.CPFeatures {
				reqLogger.Info("Processing feature", "feature", feature)
				rcheckfeat, err := r.getRest("check",instance.Spec.CPType + "-" + instance.Spec.CPVersion + "-" + feature.Name, instance.Name)
				if err != nil {
					reqLogger.Error(err, "Check result error "+instance.Spec.CPType + "-" + instance.Spec.CPVersion+"-"+feature.Name,instance.Name)
				}
				var instjsonfeat installed
				json.Unmarshal(rcheckfeat, &instjsonfeat)
				reqLogger.Info(string(rcheckfeat),"Feature check json",instjsonfeat)
				numfeat++
				if (instjsonfeat.Installed) {
					installedfeat++
					if (!contains(instance.Status.InstalledFeatures,feature.Name)) {
						instance.Status.InstalledFeatures = append(instance.Status.InstalledFeatures, feature.Name)
					}
				} else {
					yamlFeatureFile, err := ioutil.ReadFile("/cfgdata/" + instance.Spec.CPType + "-" + instance.Spec.CPVersion + "-" + feature.Name +".yaml")
					if err == nil {
						var cf CPrequirements
						err = yaml.Unmarshal(yamlFeatureFile, &cf)
						if err != nil {
							reqLogger.Error(err,"Cannot un marshall file " + instance.Spec.CPType + "-" + instance.Spec.CPVersion + "-" + feature.Name)
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
					}	else {
						reqLogger.Error(err, "feature yamlFile read error "+instance.Spec.CPType + "-" + instance.Spec.CPVersion + "-" + feature.Name + ".yaml")
						instance.Status.CPStatus = "ValidationFailed"
						statusMessage = statusMessage + "cloudpak feature unsupported yet "+instance.Spec.CPType + "-" + instance.Spec.CPVersion + " - " + feature.Name + "\n"
					}
				}
			}
		}
	}
	// Store the total requirements
	instance.Status.CPRequirement.CPReqCPU = cpureq
	instance.Status.CPRequirement.CPReqMemory = memreq
	instance.Status.CPRequirement.CPReqStorage = pvreq

	if (numfeat == installedfeat) {
		instance.Status.CPStatus = "Installed"
	} else if ((capjson.TotCpu.Cmp(cpureq)<0) || (capjson.TotMem.Cmp(memreq)<0)) {
		instance.Status.CPStatus = "NotInstallable"
	} else {
		instance.Status.CPStatus = "ReadyToInstall"
	}

	instance.Status.StatusMessages = statusMessage //Allocatable worker nodes capacity is CPU="+strconv.FormatInt(totCpu.MilliValue(),10)+"m and memory="+strconv.FormatInt(totMemory.Value(),10)+"\n"+"Requirement is CPU="+strconv.FormatInt(cpureq.MilliValue(),10)+"m and memory="+strconv.FormatInt(memreq.Value(),10)

	if (instance.Spec.Action == "Install") {

		if (instance.Status.CPStatus == "ReadyToInstall") {
		// perform installation
			rinstall, err := r.getRest("install",instance.Spec.CPType + "-" + instance.Spec.CPVersion,instance.Name)
			if err != nil {
				reqLogger.Error(err, "Install result error "+instance.Spec.CPType + "-" + instance.Spec.CPVersion)
			}
			var instjson installed
			json.Unmarshal(rinstall, &instjson)
			reqLogger.Info(string(rinstall),"json",instjson)
		}


		reqLogger.Info("Checking features for install", "feature", instance.Spec.CPFeatures)

		if (len(instance.Spec.CPFeatures) > 0) {
			/* installing sub-features */
				for _, feature := range instance.Spec.CPFeatures {

					reqLogger.Info("Installing feature", "feature", feature.Name)
					if (!contains(instance.Status.InstalledFeatures,feature.Name)) {
						reqLogger.Info("Running install for feature", "feature", feature.Name)
						rinstfeat, err := r.getRest("install",instance.Spec.CPType + "-" + instance.Spec.CPVersion + "-" + feature.Name,instance.Name)
						if err != nil {
							reqLogger.Error(err, "Check result error "+instance.Spec.CPType + "-" + instance.Spec.CPVersion+"-"+feature.Name)
						}
						var instjsonfeat installed
						json.Unmarshal(rinstfeat, &instjsonfeat)
						reqLogger.Info(string(rinstfeat),"json",instjsonfeat)
						if !instjsonfeat.Installed {
							break
						}
					}
				}
		}
	}

	reqLogger.Info("Instance info","Instance",instance)
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Status update failed")
		return reconcile.Result{}, err
	}

	reqLogger.Info("Finished reconcile")
	// Recheck status every minutes - may need to revisit
	return reconcile.Result{RequeueAfter: time.Second*300}, nil
}

func (r *ReconcileCPeir) getRest(restOper, restArg, objName string) (resp []byte, err error) {
		reqLogger := log.WithValues("RestOper", restOper, "RestArg", restArg, "object", objName)

    reqLogger.Info("Starting the application...")
		restUrl := ""
		if (objName == "") {
			restUrl = strings.Join([]string{"http://127.0.0.1:8080/",restOper,"/",restArg},"")
		} else {
			restUrl = strings.Join([]string{"http://127.0.0.1:8080/",restOper,"/",objName,"/",restArg},"")
		}
    response, err := http.Get(restUrl)
    if err != nil {
        reqLogger.Error(err,"Rest call failed")
				return nil, err
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        reqLogger.Info(string(data))
				return data, nil
    }
}
