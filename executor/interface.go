//create: 2018/01/04 11:04:08 change: 2018/01/24 16:42:20 lijiaocn@foxmail.com
package executor

import (
	L "github.com/lijiaocn/kube-lb/listener"
)

type Executor interface {
	Convert2Config(listeners []L.Listener) error
}
