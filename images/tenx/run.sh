#!/bin/sh

#NAMESPACE=""
#LOGLEVEL=1

./kube-lb -auth=incluster -namespace=$NAMESPACE -v=$LOGLEVEL -source=k8s -template=/etc/default/hafolder/haproxy.tpl -result=/etc/haproxy/haproxy.cfg.new -script=/update.sh -filter=k8s_tenx -fconfig=/etc/tenx/extention.conf -executor=haproxy  -alsologtostderr
