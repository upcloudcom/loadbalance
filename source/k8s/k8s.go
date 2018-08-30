//create: 2018/01/12 13:17:01 change: 2018/01/24 16:47:51 upcloudcom@foxmail.com
package k8s

import (
	"errors"
	"github.com/golang/glog"
	"github.com/upcloudcom/loadbalance/filter"
	L "github.com/upcloudcom/loadbalance/listener"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const (
	//auth methods
	AUTH_INCLUSTER  = "incluster"
	AUTH_KUBECONFIG = "kubeconfig"
	AUTH_TOKEN      = "token"

	ADD    = "add"
	UPDATE = "update"
	DELETE = "delete"

	PODS       = "pods"
	SERVICES   = "services"
	NAMESPACES = "namespaces"
)

type Config struct {
	Auth       string //auth methods
	KubeConfig string //kubeconfig file
	Host       string
	Token      string
	SkipTLS    bool
	Namespace  string
}

type Watch struct {
	store      cache.Store
	controller cache.Controller
	stopChan   chan struct{}
}

type KubeSource struct {
	client        *kubernetes.Clientset
	rest          *rest.RESTClient
	addFilters    []filter.FilterFunc
	updateFilters []filter.FilterFunc
	deleteFilters []filter.FilterFunc
	watchs        map[string][]*Watch
	listeners     L.Listeners
	notify        chan struct{}
}

func (self *KubeSource) Notify() chan struct{} {
	return self.notify
}

func (self *KubeSource) AddFilter(action string, funcs ...filter.FilterFunc) error {
	switch action {
	case ADD:
		self.addFilters = append(self.addFilters, funcs...)
	case UPDATE:
		self.updateFilters = append(self.updateFilters, funcs...)
	case DELETE:
		self.deleteFilters = append(self.deleteFilters, funcs...)
	default:
		return errors.New("unknown action: " + action)
	}
	return nil
}

func (self *KubeSource) AddWatch(name, resource, namespace string, selector map[string]string) error {

	lw := cache.NewListWatchFromClient(self.rest, resource, namespace, fields.SelectorFromSet(selector))
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if filter.Iterate_filters(ADD, &self.listeners, self.addFilters, obj) == filter.UPDATE {
				glog.V(2).Infof("send update")
				var s struct{}
				self.notify <- s
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if filter.Iterate_filters(UPDATE, &self.listeners, self.updateFilters, oldObj, newObj) == filter.UPDATE {
				glog.V(2).Infof("send update")
				var s struct{}
				self.notify <- s
			}
		},
		DeleteFunc: func(obj interface{}) {
			if filter.Iterate_filters(DELETE, &self.listeners, self.deleteFilters, obj) == filter.UPDATE {
				glog.V(2).Infof("send update")
				var s struct{}
				self.notify <- s
			}
		},
	}

	var store cache.Store
	var controller cache.Controller

	switch resource {
	case PODS:
		store, controller = cache.NewInformer(lw, &v1.Pod{}, 0, handler)
	case SERVICES:
		store, controller = cache.NewInformer(lw, &v1.Service{}, 0, handler)
	case NAMESPACES:
		store, controller = cache.NewInformer(lw, &v1.Namespace{}, 0, handler)
	default:
		return errors.New("unsupported resource type: " + resource)
	}

	watch := &Watch{
		store:      store,
		controller: controller,
		stopChan:   make(chan struct{}),
	}

	if _, ok := self.watchs[name]; ok {
		self.watchs[name] = append(self.watchs[name], watch)
	} else {
		var slice []*Watch
		self.watchs[name] = append(slice, watch)
	}
	return nil
}

func (self *KubeSource) Run(stopChan chan struct{}) {
	for name, slice := range self.watchs {
		glog.Infof("start watch: %s", name)
		for _, watch := range slice {
			go watch.controller.Run(watch.stopChan)
		}
	}
	for {
		select {
		case c := <-stopChan:
			for _, slice := range self.watchs {
				for _, watch := range slice {
					watch.stopChan <- c
				}
			}
			return
		}
	}
	return
}

func NewKubeSource(client *kubernetes.Clientset, rest *rest.RESTClient) (*KubeSource, error) {
	return &KubeSource{
		client:    client,
		rest:      rest,
		watchs:    make(map[string][]*Watch),
		notify:    make(chan struct{}, 100),
		listeners: L.NewListeners(),
	}, nil
}

func (self *KubeSource) List() []L.Listener {
	return self.listeners.List()
}
