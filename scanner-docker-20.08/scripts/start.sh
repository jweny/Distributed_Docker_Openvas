#!/usr/bin/env bash
set -Eeuo pipefail

MASTER_PORT=${MASTER_PORT:-22}

if [ -z $MASTER_ADDRESS ]; then
	echo "ERROR: The environment variable \"MASTER_ADDRESS\" is not set"
	exit 1
fi

if  [ ! -d /data ]; then
	echo "Creating Data folder..."
        mkdir /data
fi

if [ ! -f "/firstrun" ]; then
	echo "Running first start configuration..."

	echo "Creating Openvas NVT sync user..."
	useradd --home-dir /usr/local/share/openvas openvas-sync
	chown openvas-sync:openvas-sync -R /usr/local/share/openvas
	chown openvas-sync:openvas-sync -R /usr/local/var/lib/openvas

	touch /firstrun
fi

if [ ! -f "/data/scannerid" ]; then
	echo "Generating scanner id..."
	
	echo $(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 10 | head -n 1) > /data/scannerid
fi

if  [ ! -d /data/ssh ]; then
	mkdir /data/ssh
fi

if  [ ! -f /data/ssh/known_hosts ]; then
	echo "Getting Master SSH key..."
	ssh-keyscan -t ed25519 -p $MASTER_PORT $MASTER_ADDRESS > /data/ssh/known_hosts.temp
	mv /data/ssh/known_hosts.temp /data/ssh/known_hosts
fi

if  [ ! -f /data/ssh/key ]; then
	echo "Setup SSH key..."
	ssh-keygen -t ed25519 -f /data/ssh/key -N "" -C "$(cat /data/scannerid)"
fi

if [ ! -d "/run/redis" ]; then
	mkdir /run/redis
fi
if  [ -S /run/redis/redis.sock ]; then
        rm /run/redis/redis.sock
fi
redis-server --unixsocket /run/redis/redis.sock --unixsocketperm 700 --timeout 0 --databases 128 --maxclients 4096 --daemonize yes --port 6379 --bind 0.0.0.0

echo "Wait for redis socket to be created..."
while  [ ! -S /run/redis/redis.sock ]; do
        sleep 1
done

echo "Testing redis status..."
X="$(redis-cli -s /run/redis/redis.sock ping)"
while  [ "${X}" != "PONG" ]; do
        echo "Redis not yet ready..."
        sleep 1
        X="$(redis-cli -s /run/redis/redis.sock ping)"
done
echo "Redis ready."


if [ ! -h /usr/local/var/lib/openvas/plugins ]; then
	echo "Fixing NVT Plugins folder..."
	rm -rf /usr/local/var/lib/openvas/plugins
	ln -s /data/plugins /usr/local/var/lib/openvas/plugins
	chown openvas-sync:openvas-sync -R /data/plugins
	chown openvas-sync:openvas-sync -R /usr/local/var/lib/openvas/plugins
fi

sleep 5

if [ -f /var/run/ospd.pid ]; then
  rm /var/run/ospd.pid
fi

if [ -S /data/ospd.sock ]; then
  rm /data/ospd.sock
fi

if [ ! -d /var/run/ospd ]; then
  mkdir /var/run/ospd
fi

echo "Starting Open Scanner Protocol daemon for OpenVAS..."
ospd-openvas --log-file /usr/local/var/log/gvm/ospd-openvas.log --unix-socket /data/ospd.sock --log-level INFO

while  [ ! -S /data/ospd.sock ]; do
	sleep 1
done

chmod 666 /data/ospd.sock

touch /usr/local/var/log/gvm/ssh-connection.log
/connect.sh &

echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++"
echo "+ Your Scanner is now ready to use! +"
echo "+ UPDATE BY jweny(https://github.com/jweny)+"
echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++"
echo ""
echo "-------------------------------------------------------"
echo "Scanner id: $(cat /data/scannerid)"
echo "Public key: $(cat /data/ssh/key.pub)"
echo "Master host key (Check that it matches the public key from the master):"
cat /data/ssh/known_hosts
echo "-------------------------------------------------------"
tail -F /usr/local/var/log/gvm/*
