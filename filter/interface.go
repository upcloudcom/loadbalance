//create: 2018/01/15 19:06:03 change: 2018/01/24 14:30:15 lijiaocn@foxmail.com
package filter

import (
	"github.com/golang/glog"
	L "github.com/lijiaocn/kube-lb/listener"
)

const (
	//filter func return value
	CONTINUE = iota
	IGNORE   = iota
	UPDATE   = iota
	NOUPDATE = iota
)

type FilterFunc func(action string, ls *L.Listeners, objs ...interface{}) int

func Debug(action string, ls *L.Listeners, objs ...interface{}) int {
	glog.V(2).Infof("debug filter: %s", action)
	return CONTINUE
}

func Iterate_filters(action string, ls *L.Listeners, filters []FilterFunc, objs ...interface{}) int {
	ret := NOUPDATE
	for _, f := range filters {
		switch f(action, ls, objs...) {
		case CONTINUE:
			continue
		case IGNORE:
			return NOUPDATE
		case UPDATE:
			ret = UPDATE
		case NOUPDATE:
			continue
		default:
			continue
		}
	}
	return ret
}
