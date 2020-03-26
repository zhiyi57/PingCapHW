#! /bin/sh
set -e

if [ $SLOT = 1 ]; then 
    exec ./pd-server \
        --name $NAME \
        --client-urls http://0.0.0.0:2379 \
        --peer-urls http://0.0.0.0:2380 \
        --advertise-client-urls http://pd:2379 \
        --advertise-peer-urls http://pd:2380 \
	--log-file /tmp/logs/pd0.log
else
    exec ./pd-server \
        --name $NAME \
        --client-urls http://0.0.0.0:2379 \
        --peer-urls http://0.0.0.0:2380 \
        --advertise-client-urls http://pd1:2379 \
        --advertise-peer-urls http://pd1:2380 \
        --join http://pd.tikv:2379 \
	--log-file /tmp/logs/pd1.log
fi
