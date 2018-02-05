//create: 2018/01/04 16:46:04 change: 2018/01/25 14:00:24 lijiaocn@foxmail.com
package haproxy

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	L "github.com/lijiaocn/kube-lb/listener"
)

func checkProto(l *L.Listener) error {
	switch l.L4Proto {
	case "":
		return errors.New("l4 proto is not set")
	case L.PROTO_TCP:
	case L.PROTO_UDP:
		return errors.New("udp is not supported")
	default:
		return errors.New("unknown l4 proto " + l.L4Proto)
	}

	switch l.L7Proto {
	case "":
	case L.PROTO_SSL:
	case L.PROTO_HTTP:
	case L.PROTO_HTTPS:
	default:
		return errors.New("unknown l7 proto " + l.L7Proto)
	}
	return nil
}

func checkCert(l *L.Listener) error {
	if l.Encryption && l.CertFile == "" {
		return errors.New("use secure mode but not provide pem file")
	}
	return nil
}

func checkSticky(s *L.ServerGroup) error {
	switch s.Sticky {
	case "":
	case L.STI_COOKIE_POST:
	default:
		return errors.New("unknown sticky mode")
	}
	return nil
}

func checkAlgorithm(s *L.ServerGroup) error {
	switch s.Algorithm {
	case "":
		return errors.New("algorithm is not set")
	case L.ALG_RR:
	case L.ALG_LEAST:
	case L.ALG_FIRST:
	case L.ALG_SRC:
	case L.ALG_URI:
	default:
		return errors.New("unknown algorithm")
	}
	return nil
}

func checkServerGroup(s *L.ServerGroup, l *L.Listener) error {
	if err := checkSticky(s); err != nil {
		return err
	}

	if err := checkAlgorithm(s); err != nil {
		return err
	}
	return nil
}

func checkListener(l *L.Listener) error {

	if err := checkProto(l); err != nil {
		return err
	}

	if err := checkCert(l); err != nil {
		return err
	}

	names := make(map[string]int, len(l.ServerGroups))
	for i, s := range l.ServerGroups {
		if err := checkServerGroup(s, l); err != nil {
			return errors.New(fmt.Sprintf("serverGroup %d %s", i, err.Error()))
		}
		if pre, ok := names[s.Name]; !ok {
			names[s.Name] = i
		} else {
			return errors.New(fmt.Sprintf("serverGroup %d and serverGroup %d has a same name: %s", i, pre, s.Name))
		}
	}

	return nil
}

func valid(ls []L.Listener) error {
	glog.V(2).Infof("valid check for listeners")
	//	names := make(map[string]int, len(ls))

	for i, l := range ls {
		if err := checkListener(&l); err != nil {
			return errors.New(fmt.Sprintf("listener %d %s", i, err.Error()))
		}
		/*
			if pre, ok := names[l.Name]; !ok {
				names[l.Name] = i
			} else {
				return errors.New(fmt.Sprintf("listener %d and listener %d has a same name: %s", i, pre, l.Name))
			}
		*/

		switch l.L4Proto {
		case "":
			return errors.New("l4 proto is not set")
		case L.PROTO_TCP:
		case L.PROTO_UDP:
			return errors.New("udp is not supported")
		default:
			return errors.New("unknown l4 proto " + l.L4Proto)
		}
	}
	return nil
}
