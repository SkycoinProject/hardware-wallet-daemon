#!/usr/bin/env bash

getent group skyhwd >/dev/null || groupadd -r skyhwd
getent group plugdev >/dev/null || groupadd -r plugdev
getent passwd skyhwd >/dev/null || useradd -r -g skyhwd -d /var -s /bin/false -c "Skycoin Hardware Wallet Daemon" skyhwd
usermod -a -G plugdev skyhwd

# set home directory
mkdir /home/skyhwd
chown skyhwd:skyhwd /home/skyhwd
usermod -d /home/skyhwd skyhwd

# set log file
touch /var/log/skyhwd.log
chown skyhwd:skyhwd /var/log/skyhwd.log
chmod 660 /var/log/skyhwd.log
