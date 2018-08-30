//create: 2018/01/12 13:01:48 change: 2018/01/24 16:40:52 upcloudcom@foxmail.com
package source

import (
	"github.com/upcloudcom/loadbalance/filter"
	L "github.com/upcloudcom/loadbalance/listener"
)

type Source interface {
	List() []L.Listener
	AddWatch(name, resource, namespace string, selector map[string]string) error
	AddFilter(action string, funcs ...filter.FilterFunc) error
	Run(stopChan chan struct{})
	Notify() chan struct{}
}
