//create: 2018/01/02 19:20:09 change: 2018/01/18 13:52:56 lijiaocn@foxmail.com
package client

import (
	"github.com/golang/glog"
	"github.com/lijiaocn/kube-lb/config"
	"k8s.io/client-go/kubernetes"
)

var (
	clientset *kubernetes.Clientset
)

func InitClientSet(cmd *config.CmdLine) error {
	glog.Infof("init clientset")
	kconfig, err := ConvertToRestConfig(cmd)
	if err != nil {
		return err
	}

	clientset, err = kubernetes.NewForConfig(kconfig)
	if err != nil {
		return err
	}
	return nil
}

func GetClientSet() *kubernetes.Clientset {
	return clientset
}
