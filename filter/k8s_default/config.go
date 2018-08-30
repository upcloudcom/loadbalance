//create: 2018/01/30 19:19:21 change: 2018/01/30 19:38:19 upcloudcom@foxmail.com
package k8s_default

type LocalConfig struct {
	hostname string
	group    string
}

var localConfig LocalConfig

func init() {
	localConfig.hostname = os.HostName()
}
