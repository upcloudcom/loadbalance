package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/upcloudcom/loadbalance/client"
	"github.com/upcloudcom/loadbalance/config"
	"github.com/upcloudcom/loadbalance/executor"
	"github.com/upcloudcom/loadbalance/executor/haproxy"
	"github.com/upcloudcom/loadbalance/filter/k8s_tenx"
	"github.com/upcloudcom/loadbalance/source"
	"github.com/upcloudcom/loadbalance/source/k8s"
	"k8s.io/api/core/v1"
	"os"
	os_exec "os/exec"
	"os/signal"
	"syscall"
	"time"
)

func displayListeners(src source.Source) {
	glog.Infof("num name  bindip bindport encryption certfile l4proto l7proto")
	for i, l := range src.List() {
		glog.Infof("%d %s %s %s %t %s %s %s",
			i, l.Name, l.BindIP, l.BindPort, l.Encryption, l.CertFile, l.L4Proto, l.L7Proto)
	}
}

func notify(src source.Source, exec executor.Executor, script string, notifyStopChan chan struct{}) {
	glog.Info("start notify")
	update := src.Notify()
	for {
		select {
		case <-update:
			if l := len(update); l >= 0 {
				glog.Infof("receive %d times update", l+1)
				for i := 0; i < l; i++ {
					<-update
				}
			}

			glog.Info("prepare to generate a new config file")
			if err := exec.Convert2Config(src.List()); err != nil {
				glog.Errorf("convert to config fail: %s", err.Error())
			} else {
				glog.Info("convert to config success")
				/* if run in pod, Output() will block
				if out, err := os_exec.Command(script).Output(); err != nil {
					glog.Errorf("callback script \"%s\" fail: %s", script, err.Error())
				} else {
					glog.Infof("callback script \"%s\" output : %s", script, out)
				}
				*/
				if err := os_exec.Command(script).Run(); err != nil {
					glog.Errorf("callback script \"%s\" fail: %s", script, err.Error())
				} else {
					glog.Infof("callback script \"%s\" success", script)
				}
			}
			time.Sleep(0 * time.Second)
		case <-notifyStopChan:
			return
		}
	}
}

func main() {
	flag.Parse()
	cmdline := config.GetCmdLine()
	if cmdline.Help {
		flag.Usage()
		return
	}
	if err := config.ValidCheck(); err != nil {
		glog.Exitln(err.Error())
	}

	//init client
	if err := client.InitClientSet(cmdline); err != nil {
		glog.Exitln(err.Error())
	}
	if err := client.InitRESTClient(cmdline, v1.SchemeGroupVersion); err != nil {
		glog.Exitln(err.Error())
	}

	srcStopChan := make(chan struct{})
	notifyStopChan := make(chan struct{})
	set := client.GetClientSet()
	rest := client.GetRESTClient(v1.SchemeGroupVersion)
	var err error
	var src source.Source
	var exec executor.Executor

	//Set Source
	switch cmdline.Source {
	case "k8s":
		src, err = k8s.NewKubeSource(set, rest)
		if err != nil {
			glog.Exitln("kube source start fail: ", err.Error())
		}
		if err := src.AddWatch("watchService", k8s.SERVICES, cmdline.Namespace, nil); err != nil {
			glog.Exitln("kube source add watch fail: ", err.Error())
		}
		if err := src.AddWatch("watchPod", k8s.PODS, cmdline.Namespace, nil); err != nil {
			glog.Exitln("kube source add watch fail: ", err.Error())
		}
	default:
		glog.Exitln("unknown source: ", cmdline.Source)
	}

	//Set Filter
	switch cmdline.Filter {
	case "k8s_tenx":
		if err := src.AddFilter(k8s.ADD, k8s_tenx.DEBUG_ADD, k8s_tenx.AddCheck, k8s_tenx.Add); err != nil {
			glog.Exitln("kube source add ADD filter fail: ", err.Error())
		}
		if err := src.AddFilter(k8s.DELETE, k8s_tenx.DEBUG_DELETE, k8s_tenx.DeleteCheck, k8s_tenx.Delete); err != nil {
			glog.Exitln("kube source add DELETE filter fail: ", err.Error())
		}
		if err := src.AddFilter(k8s.UPDATE, k8s_tenx.DEBUG_UPDATE, k8s_tenx.UpdateCheck, k8s_tenx.Update); err != nil {
			glog.Exitln("kube source add UPDATE filter fail: ", err.Error())
		}
		if cmdline.FilterConfig != "" {
			go k8s_tenx.WatchConfig(cmdline.FilterConfig, 1*time.Second)
		}
	default:
		glog.Exitln("unknown source: ", cmdline.Source)
	}

	//Set Executor
	switch cmdline.Executor {
	case "haproxy":
		exec = haproxy.NewHaproxy(cmdline.DefaultIP, cmdline.Template, cmdline.Result)
	default:
	}

	go notify(src, exec, cmdline.Script, notifyStopChan)
	go src.Run(srcStopChan)

	//Set signal
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGUSR1)
	for {
		select {
		case s := <-signalChan:
			switch s {
			case syscall.SIGQUIT:
				fallthrough
			case syscall.SIGKILL:
				fallthrough
			case syscall.SIGTERM:
				var stop struct{}
				srcStopChan <- stop
				notifyStopChan <- stop
			case syscall.SIGUSR1:
				displayListeners(src)
			default:
				continue
			}
		}
	}
}
