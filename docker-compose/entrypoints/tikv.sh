#!/bin/sh
set -e

exec ./tikv-server \
    --addr 0.0.0.0:20160 \
    --status-addr 0.0.0.0:20180 \
    --advertise-addr 192.168.65.3:20160 \
    --pd-endpoints pd.tikv:2379 \
    --log-file /tmp/logs/kv.log
