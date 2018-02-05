//create: 2018/01/12 13:01:48 change: 2018/01/24 16:40:52 lijiaocn@foxmail.com
package source

import (
	"github.com/lijiaocn/kube-lb/filter"
	L "github.com/lijiaocn/kube-lb/listener"
)

type Source interface {
	List() []L.Listener
	AddWatch(name, resource, namespace string, selector map[string]string) error
	AddFilter(action string, funcs ...filter.FilterFunc) error
	Run(stopChan chan struct{})
	Notify() chan struct{}
}
