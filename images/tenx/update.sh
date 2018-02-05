#!/bin/sh

haproxy -c -f /etc/haproxy/haproxy.cfg.new
ret=$?
if [[ "$ret" != "0" ]];then
	echo "/etc/haproxy/haprox.cfg.new is wrong"
	exit 1
fi
mv /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.bak
cp /etc/haproxy/haproxy.cfg.new  /etc/haproxy/haproxy.cfg

ChildPidFile=/var/run/haproxy.pid

function reloadchild
{
  echo "Reloading"
  EXT_CMD=
  CHPID=
  if [ -f "$ChildPidFile" ]; then
     CHPID=$(cat ${ChildPidFile})
     if [ -n "$CHPID" ] && [ -n $(ps -o pid | grep  "$CHPID") ]; then
     CHPIDS=`ps aux|grep "haproxy -f"|grep -v grep|awk '{print $1}'`
     EXT_CMD="-sf $CHPID $CHPIDS"  
     fi
  fi
  haproxy -f /etc/haproxy/haproxy.cfg -db ${EXT_CMD} &
  CHPID=$!
  echo "$CHPID" >${ChildPidFile}
}

reloadchild
