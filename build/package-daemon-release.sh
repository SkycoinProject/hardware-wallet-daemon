#!/usr/bin/env bash
set -e -o pipefail

# Builds the daemon release

XGO_TARGETS="$@"

echo "In package daemon release: $XGO_TARGETS"

. build-conf.sh "$XGO_TARGETS"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

function copy_if_exists {
    if [ -z "$1" -o -z "$2" ]; then
        echo "copy_if_exists requires 2 args"
        exit 1
    fi

    TARGET="$1"
    DESTDIR="$2"

    if [ -f "$TARGET" ]; then
        if [ -e "$DESTDIR" ]; then
            rm -r "$DESTDIR"
        fi
        mkdir -p "$DESTDIR"

        # Copy target to dest
        echo "Copying $TARGET to $DESTDIR"
        cp "$TARGET" "$DESTDIR"
    else
        echo "$TARGET does not exsit"
    fi
}

echo "Copying ${PKG_NAME} binaries"

# Linux i386
if [ ! -z "$LNX64_DMN" ]; then
    LNX32_BIN_DIR="${DMN_OUTPUT_DIR}/${LNX32_DMN_BIN_DIR}"
    LNX32_RULES_DIR="${DMN_OUTPUT_DIR}/${LNX32_DMN_RULES_DIR}"
    LNX32_SERVICE_DIR="${DMN_OUTPUT_DIR}/${LNX32_DMN_SERVICE_DIR}"
    copy_if_exists "${XGO_DMN_OUTPUT_DIR}/${LNX32_OUT}/${BIN_NAME}" "$LNX32_BIN_DIR"
    copy_if_exists "./linux/skyhwd.rules" "$LNX32_RULES_DIR"
    copy_if_exists "./linux/skyhwd.service" "$LNX32_SERVICE_DIR"
fi

# Linux amd64
if [ ! -z "$LNX64_DMN" ]; then
    LNX64_BIN_DIR="${DMN_OUTPUT_DIR}/${LNX64_DMN_BIN_DIR}"
    LNX64_RULES_DIR="${DMN_OUTPUT_DIR}/${LNX64_DMN_RULES_DIR}"
    LNX64_SERVICE_DIR="${DMN_OUTPUT_DIR}/${LNX64_DMN_SERVICE_DIR}"
    copy_if_exists "${XGO_DMN_OUTPUT_DIR}/${LNX64_OUT}/${BIN_NAME}" "$LNX64_BIN_DIR"
    copy_if_exists "./linux/skyhwd.rules" "$LNX64_RULES_DIR"
    copy_if_exists "./linux/skyhwd.service" "$LNX64_SERVICE_DIR"
fi

# Linux arm
if [ ! -z "$LNX_ARM_DMN" ]; then
    LNX_ARM_BIN_DIR="${DMN_OUTPUT_DIR}/${LNX_ARM_DMN_BIN_DIR}"
    LNX_ARM_RULES_DIR="${DMN_OUTPUT_DIR}/${LNX_ARM_DMN_RULES_DIR}"
    LNX_ARM_SERVICE_DIR="${DMN_OUTPUT_DIR}/${LNX_ARM_DMN_SERVICE_DIR}"
    copy_if_exists "${XGO_DMN_OUTPUT_DIR}/${LNX_ARM_OUT}/${BIN_NAME}" "$LNX_ARM_BIN_DIR"
    copy_if_exists "./linux/skyhwd.rules" "$LNX_ARM_RULES_DIR"
    copy_if_exists "./linux/skyhwd.service" "$LNX_ARM_SERVICE_DIR"
fi

# Windows amd64
if [ ! -z "$WIN64_DMN" ]; then
    WIN64="${DMN_OUTPUT_DIR}/${WIN64_DMN}"
    copy_if_exists "${XGO_DMN_OUTPUT_DIR}/${WIN64_OUT}/${BIN_NAME}.exe" "$WIN64"
fi

# Windows 386
if [ ! -z "$WIN32_DMN" ]; then
    WIN32="${DMN_OUTPUT_DIR}/${WIN32_DMN}"
    copy_if_exists "${XGO_DMN_OUTPUT_DIR}/${WIN32_OUT}/${BIN_NAME}.exe" "$WIN32"
fi
