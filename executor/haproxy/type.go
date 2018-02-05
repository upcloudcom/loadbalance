//create: 2018/01/04 16:52:24 change: 2018/01/10 10:57:01 lijiaocn@foxmail.com
package haproxy

type HaproxyConfig struct {
	FrontendHTTP  []*FrontendHTTP
	FrontendHTTPS []*FrontendHTTPS
	FrontendTCP   []*FrontendTCP
	FrontendSSL   []*FrontendSSL
	Backend       []*Backend
}

type FrontendHTTP struct {
	Name     string
	BindIP   string
	BindPort string
	Backend  []*Backend
}

type FrontendHTTPS struct {
	Name      string
	BindIP    string
	BindPort  string
	CertFiles []string
	Backend   []*Backend
}

type FrontendTCP struct {
	Name     string
	BindIP   string
	BindPort string
	Backend  []*Backend
}

type FrontendSSL struct {
	Name      string
	BindIP    string
	BindPort  string
	CertFiles []string
	Backend   []*Backend
}

type Backend struct {
	Name      string
	ACL       string
	Hosts     []string
	Algorithm string
	Servers   []*Server
}

type Server struct {
	Name   string
	Addr   string
	Port   string
	Option string
}
