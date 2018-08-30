//create: 2018/01/16 17:16:58 change: 2018/01/22 15:28:22 upcloudcom@foxmail.com
package k8s_tenx

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"os"
	"sync"
	"time"
)

type NodeItem struct {
	HostName      string `json:"host,omitempty"`
	ListenAddress string `json:"address,omitempty"`
	Domain        string `json:"domain,omitempty"`
	Group         string `json:"group,omitempty"`
}

type Group struct {
	Address   string `json:"address,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
	Id        string `json:"id,omitempty"`
}

type ExtentionConfig struct {
	Nodes  []NodeItem `json:"nodes,omitempty"`
	Groups []Group    `json:"groups,omitempty"`
}

type LocalConfig struct {
	hostname        string
	node            NodeItem
	group           Group
	extentionConfig ExtentionConfig
	rwlock          sync.RWMutex
}

var localConfig LocalConfig

func isMyGroup(group string) bool {
	localConfig.rwlock.RLock()
	defer localConfig.rwlock.RUnlock()
	return group == localConfig.group.Id
}

func initLocalConfig(cfile string, config *LocalConfig) error {
	file, err := os.Open(cfile)
	defer file.Close()
	if err != nil {
		return err
	}

	config.rwlock.Lock()
	defer config.rwlock.Unlock()

	if err := json.NewDecoder(file).Decode(&config.extentionConfig); err != nil {
		return err
	}

	if config.hostname, err = os.Hostname(); err != nil {
		return err
	}

	var node *NodeItem = nil
	for _, n := range config.extentionConfig.Nodes {
		if n.HostName == config.hostname {
			config.node = n
			node = &config.node
		}
	}
	if node == nil {
		return errors.New("host is not found: " + config.hostname + " not in " + cfile)
	}

	var group *Group = nil
	for _, g := range config.extentionConfig.Groups {
		if config.node.Group == g.Id {
			config.group = g
			group = &g
		}
	}
	if group == nil {
		return errors.New("group is not found: " + config.node.Group + " not in " + cfile)
	}
	return nil
}

func WatchConfig(cfile string, internal time.Duration) {
	var lastchange int64 = 0
	for {
		stat, err := os.Stat(cfile)
		if err != nil {
			glog.Exitln(err.Error())
		}

		if stat.Mode()&os.ModeSymlink != 0 {
			realfile, err := os.Readlink(cfile)
			if err != nil {
				glog.Exitln(err.Error())
			}

			stat, err = os.Stat(realfile)
			if err != nil {
				glog.Exitln(err.Error())
			}
		}

		changetime := stat.ModTime().Unix()

		if changetime != lastchange {
			lastchange = changetime
			if err := initLocalConfig(cfile, &localConfig); err != nil {
				glog.Exitln(err.Error())
			}
			glog.V(1).Infof("refresh filter config: %s", localConfig)
		}

		time.Sleep(internal)
	}
}
