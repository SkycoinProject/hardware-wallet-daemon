#!/usr/bin/env bash
set -e -o pipefail

# daemon version
APP_VERSION=$(cat ../VERSION 2> /dev/null)

# package name
PKG_NAME="daemon"
# binary name
BIN_NAME="skyhwd"

if [ -n "$1" ]; then
    XGO_TARGETS="$@"
else
    XGO_TARGETS="linux/386,linux/amd64,linux/arm-7,windows/amd64,windows/386"
fi

XGO_OUTPUT_DIR=".xgo_output"
XGO_DMN_OUTPUT_DIR="${XGO_OUTPUT_DIR}/daemon"

DMN_OUTPUT_DIR=".daemon_output"

FINAL_OUTPUT_DIR="release"


# Variable suffix guide:
# _DMN -- our name for daemon releases
# _DMN_ZIP -- our compressed name for daemon releases

if echo "$XGO_TARGETS" | grep -Eq 'linux(-[0-9]{1,2}\.[0-9]{1,2})?\/386'; then
    LNX32_DMN="${BIN_NAME}-${APP_VERSION}-linux-x86"
    LNX32_DMN_BIN_DIR="${LNX32_DMN}/usr/bin/"
    LNX32_DMN_RULES_DIR="${LNX32_DMN}/lib/udev/rules.d/"
    LNX32_DMN_SERVICE_DIR="${LNX32_DMN}/usr/lib/systemd/system/"
    LNX32_DMN_ZIP="${LNX32_DMN}.tar.bz2"
    LNX32_OUT="linux_i386"
fi

if echo "$XGO_TARGETS" | grep -Eq 'linux(-[0-9]{1,2}\.[0-9]{1,2})?\/amd64'; then
    LNX64_DMN="${BIN_NAME}-${APP_VERSION}-linux-x64"
    LNX64_DMN_BIN_DIR="${LNX64_DMN}/usr/bin/"
    LNX64_DMN_RULES_DIR="${LNX64_DMN}/lib/udev/rules.d/"
    LNX64_DMN_SERVICE_DIR="${LNX64_DMN}/usr/lib/systemd/system/"
    LNX64_DMN_ZIP="${LNX64_DMN}.tar.bz2"
    LNX64_OUT="linux_x64"
fi

if echo "$XGO_TARGETS" | grep -Eq 'linux(-[0-9]{1,2}\.[0-9]{1,2})?\/arm-7'; then
    LNX_ARM_DMN="${BIN_NAME}-${APP_VERSION}-linux-arm-7"
    LNX_ARM_DMN_BIN_DIR="${LNX_ARM_DMN}/usr/bin/"
    LNX_ARM_DMN_RULES_DIR="${LNX_ARM_DMN}/lib/udev/rules.d/"
    LNX_ARM_DMN_SERVICE_DIR="${LNX_ARM_DMN}/usr/lib/systemd/system/"
    LNX_ARM_DMN_ZIP="${LNX_ARM_DMN}.tar.bz2"
    LNX_ARM_OUT="linux_arm-7"
fi

if echo "$XGO_TARGETS" | grep -Eq 'windows(-[0-9]{1,2}\.[0-9]{1,2})?\/amd64'; then
    WIN64_DMN="${BIN_NAME}-${APP_VERSION}-win-x64"
    WIN64_DMN_ZIP="${WIN64_DMN}.zip"
    WIN64_OUT="win_x64"
fi

if echo "$XGO_TARGETS" | grep -Eq 'windows(-[0-9]{1,2}\.[0-9]{1,2})?\/386'; then
    WIN32_DMN="${BIN_NAME}-${APP_VERSION}-win-x86"
    WIN32_DMN_ZIP="${WIN32_DMN}.zip"
    WIN32_OUT="win_ia32"
fi

