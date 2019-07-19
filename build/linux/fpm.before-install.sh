#!/usr/bin/env bash

getent group skyhwd >/dev/null || groupadd -r skyhwd
getent group plugdev >/dev/null || groupadd -r plugdev
getent passwd skyhwd >/dev/null || useradd -r -g skyhwd -d /var -s /bin/false -c "Skycoin Hardware Wallet Daemon" skyhwd
usermod -a -G plugdev skyhwd
touch /var/log/skyhwd.log
chown skyhwd:skyhwd /var/log/skyhwd.log
chmod 660 /var/log/skyhwd.log
