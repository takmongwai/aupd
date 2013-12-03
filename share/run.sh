#!/bin/bash
ulimit -n 655350

export GOMAXPROCS=`cat /proc/cpuinfo | grep -v grep | grep process | wc -l`
#export GODEBUG='gctrace=1,schedtrace=1,scheddetail=1'
export GODEBUG='gctrace=1,scheddetail=1'
/srv/aupd/bin/aupd 
