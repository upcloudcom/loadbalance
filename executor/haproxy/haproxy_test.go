//create: 2018/01/08 14:24:45 change: 2018/01/11 14:28:05 upcloudcom@foxmail.com
package haproxy

import (
	"fmt"
	_ "github.com/golang/glog"
	L "github.com/upcloudcom/loadbalance/listener"
	"testing"
)

var ha *Haproxy
var ls []L.Listener

func TestConvert2Config(t *testing.T) {
	if err := ha.Convert2Config(ls); err != nil {
		t.Error(err)
	}
}

func http_listener(name, ip, port string, hosts ...string) L.Listener {
	l := L.Listener{
		Name:         name,
		BindIP:       ip,
		BindPort:     port,
		Encryption:   false,
		CertFile:     "",
		L4Proto:      L.PROTO_TCP,
		L7Proto:      L.PROTO_HTTP,
		ServerGroups: []L.ServerGroup{},
	}

	sg1 := L.ServerGroup{
		Name:      "sg1",
		L4:        L.L4Condition{},
		L7:        L.L7Condition{},
		Sticky:    L.STI_COOKIE_POST,
		Algorithm: L.ALG_RR,
		Servers:   []L.Server{},
	}
	server1 := L.Server{
		Name:        "server1",
		Addr:        "1.2.3.3",
		Port:        "33",
		MaxConn:     200,
		HealthCheck: true,
	}
	sg1.Servers = append(sg1.Servers, server1)
	sg1.L7.Hosts = append(sg1.L7.Hosts, hosts...)
	l.ServerGroups = append(l.ServerGroups, sg1)
	return l
}

func https_listener(name, ip, port string, hosts ...string) L.Listener {
	l := L.Listener{
		Name:         name,
		BindIP:       ip,
		BindPort:     port,
		Encryption:   true,
		CertFile:     "cert.pem",
		L4Proto:      L.PROTO_TCP,
		L7Proto:      L.PROTO_HTTPS,
		ServerGroups: []L.ServerGroup{},
	}

	sg1 := L.ServerGroup{
		Name:      "sg1",
		L4:        L.L4Condition{},
		L7:        L.L7Condition{},
		Sticky:    L.STI_COOKIE_POST,
		Algorithm: L.ALG_RR,
		Servers:   []L.Server{},
	}
	server1 := L.Server{
		Name:        "server1",
		Addr:        "1.2.3.3",
		Port:        "33",
		MaxConn:     200,
		HealthCheck: true,
	}

	sg1.Servers = append(sg1.Servers, server1)
	sg1.L7.Hosts = append(sg1.L7.Hosts, hosts...)
	l.ServerGroups = append(l.ServerGroups, sg1)
	return l
}

func tcp_listener(name, ip, port string) L.Listener {
	l := L.Listener{
		Name:         name,
		BindIP:       ip,
		BindPort:     port,
		Encryption:   false,
		CertFile:     "",
		L4Proto:      L.PROTO_TCP,
		L7Proto:      "",
		ServerGroups: []L.ServerGroup{},
	}

	sg1 := L.ServerGroup{
		Name:      "sg1",
		L4:        L.L4Condition{},
		L7:        L.L7Condition{},
		Sticky:    L.STI_COOKIE_POST,
		Algorithm: L.ALG_RR,
		Servers:   []L.Server{},
	}
	server1 := L.Server{
		Name:        "server1",
		Addr:        "1.2.3.3",
		Port:        "33",
		MaxConn:     200,
		HealthCheck: true,
	}

	sg1.Servers = append(sg1.Servers, server1)
	l.ServerGroups = append(l.ServerGroups, sg1)
	return l
}

func ssl_listener(name, ip, port string) L.Listener {
	l := L.Listener{
		Name:         name,
		BindIP:       ip,
		BindPort:     port,
		Encryption:   true,
		CertFile:     "cert.pem",
		L4Proto:      L.PROTO_TCP,
		L7Proto:      L.PROTO_SSL,
		ServerGroups: []L.ServerGroup{},
	}

	sg1 := L.ServerGroup{
		Name:      "sg1",
		L4:        L.L4Condition{},
		L7:        L.L7Condition{},
		Sticky:    L.STI_COOKIE_POST,
		Algorithm: L.ALG_RR,
		Servers:   []L.Server{},
	}
	server1 := L.Server{
		Name:        "server1",
		Addr:        "1.2.3.3",
		Port:        "33",
		MaxConn:     200,
		HealthCheck: true,
	}

	sg1.Servers = append(sg1.Servers, server1)
	l.ServerGroups = append(l.ServerGroups, sg1)
	return l
}

func TestMain(m *testing.M) {
	ha = NewHaproxy("0.0.0.0", "./haproxy.tpl", "./haproxy.conf")
	ls = append(ls, http_listener("http1", "1.1.1.1", "1", "http1.com", "http1.1.com"))
	ls = append(ls, http_listener("http2", "1.1.1.1", "1", "http2.com", "http2.1.com"))
	ls = append(ls, https_listener("https1", "1.1.1.1", "2", "http2.com", "http2.1.com"))
	ls = append(ls, https_listener("https2", "1.1.1.1", "2", "http2.com", "http2.1.com"))
	ls = append(ls, tcp_listener("tcp1", "1.1.1.1", "3"))
	ls = append(ls, tcp_listener("tcp2", "1.1.1.1", "3"))
	ls = append(ls, ssl_listener("ssl1", "1.1.1.1", "4"))
	ls = append(ls, ssl_listener("ssl2", "1.1.1.1", "4"))

	ls = append(ls, http_listener("http11", "1.1.1.1", "11", "http1.com", "http1.1.com"))
	ls = append(ls, http_listener("http22", "1.1.1.1", "11", "http2.com", "http2.1.com"))
	ls = append(ls, https_listener("https11", "1.1.1.1", "22", "http2.com", "http2.1.com"))
	ls = append(ls, https_listener("https22", "1.1.1.1", "22", "http2.com", "http2.1.com"))
	ls = append(ls, tcp_listener("tcp11", "1.1.1.1", "33"))
	ls = append(ls, tcp_listener("tcp22", "1.1.1.1", "33"))
	ls = append(ls, ssl_listener("ssl11", "1.1.1.1", "44"))
	ls = append(ls, ssl_listener("ssl22", "1.1.1.1", "44"))

	fmt.Println("ls is ", ls)
	ha.Convert2Config(ls)
	m.Run()
}
