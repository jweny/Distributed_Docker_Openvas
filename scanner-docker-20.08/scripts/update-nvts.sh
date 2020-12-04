#!/usr/bin/env bash
# This script update the NVTs in the background every 12 hours.
set -Eeuo pipefail

while true; do
	echo "Running Automatic NVT update..."
	su -c "rsync --compress-level=9 --links --times --omit-dir-times --recursive --partial --quiet rsync://feed.community.greenbone.net:/nvt-feed /usr/local/var/lib/openvas/plugins" openvas-sync
	sleep 43200
done
