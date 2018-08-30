package client

import (
	"github.com/golang/glog"
	"github.com/upcloudcom/loadbalance/config"
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
