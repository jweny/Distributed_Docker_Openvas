#!/bin/bash

apt-get update

{ cat <<EOF
bison
build-essential
ca-certificates
cmake
curl
gcc
gcc-mingw-w64
geoip-database
gnutls-bin
heimdal-dev
ike-scan
libgcrypt20-dev
libglib2.0-dev
libgnutls28-dev
libgpgme11-dev
libgpgme-dev
libhiredis-dev
libical-dev
libksba-dev
libldap2-dev
libmicrohttpd-dev
libnet-snmp-perl
libpcap-dev
libpopt-dev
libsnmp-dev
libssh-gcrypt-dev
libxml2-dev
net-tools
nmap
openssh-client
perl-base
python3-bcrypt
python3-cffi
python3-cryptography
python3-defusedxml
python3-lxml
python3-packaging
python3-paramiko
python3-pip
python3-psutil
python3-pycparser
python3-pyparsing
python3-redis
python3-setuptools
python3-six
redis-server
redis-tools
rsync
smbclient
uuid-dev
wapiti
wget
EOF
} | xargs apt-get install -yq --no-install-recommends


rm -rf /var/lib/apt/lists/*
