#!/bin/bash
#kill -9 `ps aux | grep -v grep | grep "/srv/aupd/bin/aupd" | awk '{print $2}'`
[[ -s "/srv/aupd/tmp/aupd.pid" ]] &&  kill -9 `cat /srv/aupd/tmp/aupd.pid`
#nohup /srv/aupd/share/run.sh 2>/dev/null 1>/dev/null &
nohup /srv/aupd/share/run.sh 2>&1 >> /srv/aupd/log/stdout.log &
