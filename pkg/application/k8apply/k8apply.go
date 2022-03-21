package k8apply

import (
	"flag"
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
	"path/filepath"
	"strings"
	"time"
)

var restConfig *rest.Config
var clientset *kubernetes.Clientset
var decoder runtime.Decoder

func InitKubeClient() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var err error
	restConfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}

	clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	decoder = scheme.Codecs.UniversalDeserializer()
	return nil
}

func decodeObj(data []byte) (runtime.Object, error) {
	obj, _, err := decoder.Decode(data, nil, nil)
	return obj, err
}

var applyChan chan []byte
var lastApply = time.Now()

func A(data []byte) {
	for d := range applyChan {
		Apply(d)
	}
}

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

	//_, err = restHelper.DeleteWithOptions(namespace, name, &metav1.DeleteOptions{})

	_, err = restHelper.Create(namespace, true, obj)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return restHelper.Replace(namespace, name, true, obj)
	}
	return obj, err
}

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
