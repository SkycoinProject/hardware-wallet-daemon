#!/usr/bin/env bash

getent group skycoin >/dev/null || groupadd -r skycoin
getent group plugdev >/dev/null || groupadd -r plugdev
getent passwd skycoin >/dev/null || useradd -r -g skycoin -d /var -s /bin/false -c "Skycoin Hardware Wallet Daemon" skycoin
usermod -a -G plugdev skycoin
touch /var/log/skyhwd.log
chown skycoin:skycoin /var/log/skyhwd.log
chmod 660 /var/log/skyhwd.log
