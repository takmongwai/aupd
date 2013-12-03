#!/bin/bash
ulimit -n 655350

export GOMAXPROCS=`cat /proc/cpuinfo | grep -v grep | grep process | wc -l`
export GODEBUG='gctrace=1'

GODEBUG='gctrace=1,schedtrace=1,scheddetail'  GOMAXPROCS=8 /srv/aupd/bin/aupd 
