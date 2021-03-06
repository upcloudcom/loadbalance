#!/bin/bash
mkdir /var/lib/haproxy
mkdir /var/run/haproxy

./loadbalance -auth=kubeconfig -host=https://10.39.1.62:6443 -kubeconfig=./kubeconfig-62 -skiptls=true -namespace="" -v=1 -source=k8s -template=./images/tenx/haproxy.tpl -result=/etc/haproxy/haproxy.cfg.new -script=./images/tenx/update.sh -filter=k8s_tenx -fconfig=./filter/k8s_tenx/extention.conf -executor=haproxy  -alsologtostderr

#Usage of ./loadbalance:
#  -alsologtostderr
#    	log to standard error as well as files
#  -auth string
#    	auth method: incluster,kubeconfig,token
#  -defaultip string
#    	defaultip (default "0.0.0.0")
#  -executor string
#    	executor type: haproxy
#  -filter string
#    	filter type: k8s_tenx
#  -help
#    	show usage
#  -host string
#    	kubernetes api host
#  -kubeconfig string
#    	kubeconfig file
#  -log_backtrace_at value
#    	when logging hits line file:N, emit a stack trace
#  -log_dir string
#    	If non-empty, write log files in this directory
#  -logtostderr
#    	log to standard error instead of files
#  -namespace string
#    	namespace
#  -result string
#    	output result file (default "./result.conf")
#  -script string
#    	callback script
#  -skiptls
#    	don't verify TLS certificate (default true)
#  -source string
#    	source type: k8s
#  -stderrthreshold value
#    	logs at or above this threshold go to stderr
#  -template string
#    	tempalte file
#  -token string
#    	user's bearer token
#  -v value
#    	log level for V logs
#  -vmodule value
#    	comma-separated list of pattern=N settings for file-filtered logging
