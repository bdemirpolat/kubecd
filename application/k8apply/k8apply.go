package k8apply

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	restConfig *rest.Config
	clientset  *kubernetes.Clientset
	decoder    runtime.Decoder
)

func init() {
	applyChan = make(chan ApplyWithChan)
	go consumeApplies()
}

// InitKubeClient initializes kubernetes client with CLUSTER_TYPE choose
func InitKubeClient() error {
	var kubeConfig *string
	var err error

	if os.Getenv("KUBECD_CLUSTER_TYPE") == "OUT_OF_CLUSTER" {
		if kubeconfigFromEnv := os.Getenv("KUBECONFIG"); kubeconfigFromEnv != "" {
			kubeConfig = &kubeconfigFromEnv
		} else {
			kubeConfigWithHomeDir := filepath.Join(homedir.HomeDir(), ".kube", "config")
			kubeConfig = &kubeConfigWithHomeDir
		}
		restConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
		if err != nil {
			return err
		}
	} else {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}

	clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	decoder = scheme.Codecs.UniversalDeserializer()
	return nil
}

// decodeObj decodes raw yaml or json files to runtime.Object
func decodeObj(data []byte) (runtime.Object, error) {
	obj, _, err := decoder.Decode(data, nil, nil)
	return obj, err
}

type ApplyWithChan struct {
	Data []byte
	C    chan error
}

// applyChan uses for queue mechanism, all kubernetes manifests must apply with this chan because of kubernetes rate limit
var applyChan chan ApplyWithChan

// lastApply uses for detecting rate limit diff
var lastApply = time.Now()

// AddToApplyQueue adds raw manifest to applyChan
func AddToApplyQueue(data []byte, c chan error) {
	applyChan <- ApplyWithChan{
		Data: data,
		C:    c,
	}
}

// consumeApplies consumes applyChan
func consumeApplies() {
	for d := range applyChan {
		diff := time.Since(lastApply)
		if diff < time.Millisecond*2000 {
			time.Sleep(time.Millisecond * 2000)
		}
		d.C <- Apply(d.Data)
	}
}

// Apply like kubectl apply, tries to create new object, if it exists tries to replace
func Apply(data []byte) error {
	lastApply = time.Now()
	decodedObj, err := decodeObj(data)
	if err != nil {
		return err
	}

	_, err = applyObject(clientset, *restConfig, decodedObj)
	if err != nil {
		return err
	}
	return nil
}

// applyObject applies with runtime.Object, tries to create, if it exists tries to replace
func applyObject(kubeClientset kubernetes.Interface, restConfig rest.Config, obj runtime.Object) (runtime.Object, error) {
	groupResources, err := restmapper.GetAPIGroupResources(kubeClientset.Discovery())
	if err != nil {
		return nil, err
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)

	groupVersionKind := obj.GetObjectKind().GroupVersionKind()
	groupKind := schema.GroupKind{Group: groupVersionKind.Group, Kind: groupVersionKind.Kind}
	mapping, err := rm.RESTMapping(groupKind, groupVersionKind.Version)
	if err != nil {
		return nil, err
	}

	restClient, err := newRestClient(restConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		return nil, err
	}

	name, err := meta.NewAccessor().Name(obj)
	if err != nil {
		return nil, err
	}

	namespace, err := meta.NewAccessor().Namespace(obj)
	if err != nil || namespace == "" {
		namespace = "default"
	}
	restHelper := resource.NewHelper(restClient, mapping)

	// _, err = restHelper.DeleteWithOptions(namespace, name, &metav1.DeleteOptions{})

	_, err = restHelper.Create(namespace, true, obj)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return restHelper.Replace(namespace, name, true, obj)
	}
	return obj, err
}

// newRestClient
func newRestClient(restConfig rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
	restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	restConfig.GroupVersion = &gv
	if len(gv.Group) == 0 {
		restConfig.APIPath = "/api"
	} else {
		restConfig.APIPath = "/apis"
	}

	return rest.RESTClientFor(&restConfig)
}
