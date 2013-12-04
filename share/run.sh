#!/bin/bash
ulimit -n 655350

export GOMAXPROCS=`cat /proc/cpuinfo | grep -v grep | grep process | wc -l`
#export GODEBUG='gctrace=1,schedtrace=1,scheddetail=1'
export GODEBUG='gctrace=1,scheddetail=1'
#/srv/aupd/bin/aupd 
/srv/aupd/bin/aupd  -log=/srv/aupd/log/stdout.log -pid=/srv/aupd/tmp/aupd.pid


#Usage: /srv/aupd/bin/aupd
#  -h="0.0.0.0": listen host
#  -log="": log file
#  -p=8000: listen port
#  -pid="": pid file
