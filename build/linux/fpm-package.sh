#!/usr/bin/env bash

set -e

pushd "linux" >/dev/null

TARGET=$1
TYPE=$2
NAME=skyhwd
VERSION=$(cat ../../VERSION)
  case "$TARGET-$TYPE" in
      linux-386-*)
          ARCH=i386
          ARCH_NAME=linux-x86
          ;;
      linux-amd64-deb)
          ARCH=amd64
          ARCH_NAME=linux-x64
          ;;
      linux-amd64-rpm)
          ARCH=x86_64
          ARCH_NAME=linux-x64
          ;;
      linux-arm-7-deb)
          ARCH=armhf
          ARCH_NAME=linux-arm-7
          ;;
      linux-arm-7-rpm)
          ARCH=armv7hl
          ARCH_NAME=linux-arm-7
          ;;
      linux-arm64-*)
          ARCH=arm64
          ARCH_NAME=linux-arm
          ;;
  esac
  fpm \
      -s tar \
      -t $TYPE \
      -a $ARCH \
      -n $NAME \
      -v $VERSION \
      -d systemd \
      -p ../release \
      --vendor "Skycoin Foundation" \
      --description "Communication daemon for Skycoin Hardware Wallet" \
      --maintainer "therealssj <therealssj@sskycoinmail.com>" \
      --url "https://skycoin.net/" \
      --category "Productivity/Security" \
      --before-install fpm.before-install.sh \
      --after-install fpm.after-install.sh \
      --before-remove fpm.before-remove.sh \
      ../release/$NAME-$VERSION-$ARCH_NAME.tar.bz2

popd >/dev/null
