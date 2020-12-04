#!/usr/bin/env bash
set -Eeuo pipefail

MASTER_PORT=${MASTER_PORT:-22}

SCANNER_ID=$(cat /data/scannerid)

until ssh -N -T -i /data/ssh/key -o UserKnownHostsFile=/data/ssh/known_hosts -p $MASTER_PORT -R /sockets/$SCANNER_ID.sock:/data/ospd.sock gvm@$MASTER_ADDRESS && echo "Connected to GVM."; do
	echo "Connection disrupted, retrying in 10 seconds..." >> /usr/local/var/log/gvm/ssh-connection.log
	sleep 10
done
