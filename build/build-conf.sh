#!/usr/bin/env bash
set -e -o pipefail

# daemon version
APP_VERSION=$(cat ./VERSION 2> /dev/null)

# package name
PKG_NAME="daemon"
# binary name
BIN_NAME="skyd"

if [ -n "$1" ]; then
    GOX_OSARCH="$@"
else
    GOX_OSARCH="linux/amd64 linux/arm windows/amd64 windows/386 darwin/amd64"
fi

GOX_OUTPUT_DIR=".gox_output"
GOX_DMN_OUTPUT_DIR="${GOX_OUTPUT_DIR}/daemon"

DMN_OUTPUT_DIR=".daemon_output"

FINAL_OUTPUT_DIR="release"


# Variable suffix guide:
# _DMN -- our name for daemon releases
# _DMN_ZIP -- our compressed name for daemon releases

if [[ $GOX_OSARCH == *"darwin/amd64"* ]]; then
    OSX64_DMN="${PKG_NAME}-${APP_VERSION}-osx-darwin-x64"
    OSX64_DMN_ZIP="${OSX64_DMN}.zip"
    OSX64_OUT="mac_x64"
fi

if [[ $GOX_OSARCH == *"linux/amd64"* ]]; then
    LNX64_DMN="${PKG_NAME}-${APP_VERSION}-linux-x64"
    LNX64_DMN_ZIP="${LNX64_DMN}.tar.gz"
    LNX64_OUT="linux_x64"
fi

if [[ $GOX_OSARCH == *"windows/amd64"* ]]; then
    WIN64_DMN="${PKG_NAME}-${APP_VERSION}-win-x64"
    WIN64_DMN_ZIP="${WIN64_DMN}.zip"
    WIN64_OUT="win_x64"
fi

if [[ $GOX_OSARCH == *"windows/386"* ]]; then
    WIN32_DMN="${PKG_NAME}-${APP_VERSION}-win-x86"
    WIN32_DMN_ZIP="${WIN32_DMN}.zip"
    WIN32_OUT="win_ia32"
fi

if [[ $GOX_OSARCH == *"linux/arm"* ]]; then
    LNX_ARM_DMN="${PKG_NAME}-${APP_VERSION}-linux-arm"
    LNX_ARM_DMN_ZIP="${LNX_ARM_DMN}.tar.gz"
    LNX_ARM_OUT="linux_arm"
fi
