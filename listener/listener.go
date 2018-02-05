//create: 2018/01/03 16:24:59 change: 2018/01/26 20:32:25 lijiaocn@foxmail.com
package listener

import (
	"github.com/golang/glog"
	"sync"
)

const (
	PROTO_UDP       = "udp"
	PROTO_TCP       = "tcp"
	PROTO_SSL       = "ssl"
	PROTO_HTTP      = "http"
	PROTO_HTTPS     = "https"
	ALG_RR          = "roundrobin"
	ALG_LEAST       = "leastconn"
	ALG_FIRST       = "first"
	ALG_SRC         = "source"
	ALG_URI         = "uri"
	STI_COOKIE_POST = "postonly"
)

type L7Condition struct {
	Hosts []string
}
type L4Condition struct {
	Undefined string
}

type Server struct {
	Name        string
	Addr        string
	MaxConn     int
	HealthCheck bool
}

type ServerGroup struct {
	Name string //default
	Port string
	L4   L4Condition
	L7   L7Condition

	Sticky    string //STI_COOKIE_POST
	Algorithm string //ALG_RR、ALG_LEAST、ALG_FIRST、ALG_SRC、ALG_URI
	Servers   []Server
}

type Listener struct {
	Name       string //Namespace-ServiceName
	BindIP     string
	BindPort   string
	Encryption bool
	CertFile   string

	L4Proto string //PROTO_TCP、PROTO_UDP
	L7Proto string //PROTO_HTTPS、PROTO_HTTP、

	ServerGroups []*ServerGroup
}

func (l *Listener) Key() string {
	return l.Name + "-" + l.BindIP + "-" + l.BindPort
}

type Listeners struct {
	//key: listenerName(namespace-svcname)-bindip-bindport
	Listeners map[string]Listener
	//key1: listenerName(namespace-svcname)-servergroupName , key2: servername
	//svcname is read from labels:
	//  enncloud.com/statefulSetName
	//  tenxcloud.com/petsetName
	//  tenxcloud.com/appName
	//  tenxcloud.com/svcName
	//  name
	Servers map[string]*map[string]Server
	Rwlock  sync.RWMutex
}

func NewListeners() Listeners {
	ls := Listeners{
		Listeners: make(map[string]Listener, 0),
		Servers:   make(map[string]*map[string]Server, 0),
	}
	return ls
}

func (self *Listeners) List() []Listener {
	self.Rwlock.RLock()
	defer self.Rwlock.RUnlock()

	ls := make([]Listener, 0)
	for _, l := range self.Listeners {
		for _, sg := range l.ServerGroups {
			sg.Servers = make([]Server, 0)
			key := l.Name + "-" + sg.Name
			if smap, ok := self.Servers[key]; !ok {
				continue
			} else {
				for _, v := range *smap {
					glog.V(1).Infof("use server: %s %s", key, v.Name)
					sg.Servers = append(sg.Servers, v)
				}
			}
		}
		ls = append(ls, l)
	}
	return ls
}
