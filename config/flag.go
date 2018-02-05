//create: 2018/01/02 18:52:14 change: 2018/01/22 13:58:14 lijiaocn@foxmail.com
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/golang/glog"
	"os"
)

const (
	//auth methods
	AUTH_INCLUSTER  = "incluster"
	AUTH_KUBECONFIG = "kubeconfig"
	AUTH_TOKEN      = "token"
)

func init() {
	flag.BoolVar(&cmdline.Help, "help", false, "show usage")
	flag.StringVar(&cmdline.Auth, "auth", "", "auth method: "+
		AUTH_INCLUSTER+","+AUTH_KUBECONFIG+","+AUTH_TOKEN)
	flag.StringVar(&cmdline.KubeConfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&cmdline.Host, "host", "", "kubernetes api host")
	flag.StringVar(&cmdline.Token, "token", "", "user's bearer token")
	flag.BoolVar(&cmdline.SkipTLS, "skiptls", true, "don't verify TLS certificate")
	flag.StringVar(&cmdline.Source, "source", "", "source type: k8s")
	flag.StringVar(&cmdline.Filter, "filter", "", "filter type: k8s_tenx")
	flag.StringVar(&cmdline.FilterConfig, "fconfig", "", "filter config")
	flag.StringVar(&cmdline.Executor, "executor", "", "executor type: haproxy")
	flag.StringVar(&cmdline.Namespace, "namespace", "", "namespace")
	flag.StringVar(&cmdline.DefaultIP, "defaultip", "0.0.0.0", "defaultip")
	flag.StringVar(&cmdline.Template, "template", "", "tempalte file")
	flag.StringVar(&cmdline.Result, "result", "./result.conf", "output result file")
	flag.StringVar(&cmdline.Script, "script", "", "callback script")
}

type CmdLine struct {
	Help         bool
	Auth         string
	KubeConfig   string
	Host         string
	Token        string
	SkipTLS      bool
	Source       string
	Filter       string
	FilterConfig string
	Executor     string
	Namespace    string

	DefaultIP string
	Template  string
	Result    string
	Script    string
}

var cmdline CmdLine

func checkAuth(cmdline *CmdLine) error {
	switch cmdline.Auth {
	case "":
		return errors.New("auth method is not set by -auth")
	case AUTH_INCLUSTER:

	case AUTH_KUBECONFIG:
		if cmdline.KubeConfig == "" {
			return errors.New("must specify the kubeconfig file by -kubeconfig")
		}
	case AUTH_TOKEN:
		if cmdline.Host == "" {
			return errors.New("must specify the host by -host")
		}
	default:
		return errors.New("unknown auth method: " + cmdline.Auth)
	}
	return nil
}

func checkSource(cmdline *CmdLine) error {
	switch cmdline.Source {
	case "":
		return errors.New("must specify one source:  k8s")
	case "k8s":
	default:
		return errors.New("unknown source: " + cmdline.Source)
	}
	return nil
}

func checkFilter(cmdline *CmdLine) error {
	switch cmdline.Filter {
	case "":
		return errors.New("must specify one filter:  k8s_tenx")
	case "k8s_tenx":
	default:
		return errors.New("unknown filter: " + cmdline.Filter)
	}
	return nil
}

func checkExecutor(cmdline *CmdLine) error {
	switch cmdline.Executor {
	case "":
		return errors.New("must specify one executor:  haproxy")
	case "haproxy":
	default:
		return errors.New("unknown executor: " + cmdline.Executor)
	}
	return nil
}

func checkTemplate(cmdline *CmdLine) error {
	if cmdline.Template == "" {
		return errors.New("must specify one template")
	}
	if _, err := os.Stat(cmdline.Template); os.IsNotExist(err) {
		return errors.New("template file not found: " + cmdline.Template)
	}
	return nil
}

func checkScript(cmdline *CmdLine) error {
	if cmdline.Script == "" {
		return errors.New("must specify one script")
	}
	if _, err := os.Stat(cmdline.Script); os.IsNotExist(err) {
		return errors.New("script file not found: " + cmdline.Script)
	}
	return nil
}

func ValidCheck() error {

	if b, err := json.Marshal(cmdline); err != nil {
		glog.Exitf("marshal cmdline fail: %s\n", err.Error())
	} else {
		glog.Infof("cmdline is: %s\n", string(b))
	}

	if err := checkAuth(&cmdline); err != nil {
		return err
	}

	if err := checkSource(&cmdline); err != nil {
		return err
	}

	if err := checkFilter(&cmdline); err != nil {
		return err
	}

	if err := checkExecutor(&cmdline); err != nil {
		return err
	}

	if err := checkTemplate(&cmdline); err != nil {
		return err
	}

	if err := checkScript(&cmdline); err != nil {
		return err
	}

	return nil
}

func GetCmdLine() *CmdLine {
	return &cmdline
}
