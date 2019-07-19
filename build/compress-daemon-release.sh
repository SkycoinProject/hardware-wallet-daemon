#!/usr/bin/env bash
set -e -o pipefail

# Compresses packaged daemon release after
# ./package-daemon-release.sh is done

XGO_TARGETS="$@"

. build-conf.sh "$XGO_TARGETS"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$DMN_OUTPUT_DIR" >/dev/null

FINALS=()

# OS X
if [ -e "$OSX64_DMN" ]; then
    if [ -e "$OSX64_DMN_ZIP" ]; then
        echo "Removing old $OSX64_DMN_ZIP"
        rm "$OSX64_DMN_ZIP"
    fi
    echo "Zipping $OSX64_DMN_ZIP"
    # -y preserves symlinks,
    # so that the massive .framework library isn't duplicated
    zip -r -y --quiet "$OSX64_DMN_ZIP" "$OSX64_DMN"
    FINALS+=("$OSX64_DMN_ZIP")
fi

# Windows 64bit
if [ -e "$WIN64_DMN" ]; then
    if [ -e "$WIN64_DMN_ZIP" ]; then
        echo "Removing old $WIN64_DMN_ZIP"
        rm "$WIN64_DMN_ZIP"
    fi
    echo "Zipping $WIN64_DMN_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        zip -r --quiet -X "$WIN64_DMN_ZIP"  "$WIN64_DMN"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        zip -r --quiet "$WIN64_DMN_ZIP" "$WIN64_DMN"
    elif [[ "$OSTYPE" == "msys"* ]]; then
        7z a "$WIN64_DMN_ZIP" "$WIN64_DMN"
    fi
    FINALS+=("$WIN64_DMN_ZIP")
fi

# Windows 32bit
if [ -e "$WIN32_DMN" ]; then
    if [ -e "$WIN32_DMN_ZIP" ]; then
        echo "Removing old $WIN32_DMN_ZIP"
        rm "$WIN32_DMN_ZIP"
    fi
    echo "Zipping $WIN32_DMN_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        zip -r --quiet -X "$WIN32_DMN_ZIP"  "$WIN32_DMN"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        zip -r --quiet "$WIN32_DMN_ZIP" "$WIN32_DMN"
    elif [[ "$OSTYPE" == "msys"* ]]; then
        7z a "$WIN32_DMN_ZIP" "$WIN32_DMN"
    fi
    FINALS+=("$WIN32_DMN_ZIP")
fi

# Linux i386
if [ -e "$LNX32_DMN" ]; then
    if [ -e "$LNX32_DMN_ZIP" ]; then
        echo "Removing old $LNX32_DMN_ZIP"
        rm "$LNX32_DMN_ZIP"
    fi
    echo "Zipping $LNX32_DMN_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar cjf "$LNX32_DMN_ZIP" --owner=0 --group=0 "$LNX32_DMN"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar cjf "$LNX32_DMN_ZIP" -C ${LNX32_DMN} .
    fi
    FINALS+=("$LNX32_DMN_ZIP")
fi

# Linux AMD64
if [ -e "$LNX64_DMN" ]; then
    if [ -e "$LNX64_DMN_ZIP" ]; then
        echo "Removing old $LNX64_DMN_ZIP"
        rm "$LNX64_DMN_ZIP"
    fi
    echo "Zipping $LNX64_DMN_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar cjf "$LNX64_DMN_ZIP" --owner=0 --group=0 "$LNX64_DMN"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar cjf "$LNX64_DMN_ZIP" -C ${LNX64_DMN} .
    fi
    FINALS+=("$LNX64_DMN_ZIP")
fi

# Linux arm-7
if [ -e "$LNX_ARM_DMN" ]; then
    if [ -e "$LNX_ARM_DMN_ZIP" ]; then
        echo "Removing old $LNX_ARM_DMN_ZIP"
        rm "$LNX_ARM_DMN_ZIP"
    fi
    echo "Zipping $LNX_ARM_DMN_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar cjf "$LNX_ARM_DMN_ZIP" --owner=0 --group=0 "$LNX_ARM_DMN"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      tar cjf "$LNX_ARM_DMN_ZIP" -C ${LNX_ARM_DMN} .
    fi
    FINALS+=("$LNX_ARM_DMN_ZIP")
fi

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT_DIR"
for var in "${FINALS[@]}"; do

    mv "${DMN_OUTPUT_DIR}/${var}" "$FINAL_OUTPUT_DIR"
done

if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
# Create linux packages
echo "Create linux packages"
if [ -e "$DMN_OUTPUT_DIR/$LNX64_DMN" ]; then
  echo "Create linux amd64 packagess"
  ./linux/fpm-package.sh linux-amd64 deb
  ./linux/fpm-package.sh linux-amd64 rpm
fi

if [ -e "$DMN_OUTPUT_DIR/$LNX32_DMN" ]; then
  echo "Creating linux 386 packages"
  ./linux/fpm-package.sh linux-386 deb
  ./linux/fpm-package.sh linux-386 rpm
fi

if [ -e "$DMN_OUTPUT_DIR/$LNX_ARM_DMN" ]; then
  echo "Create linux arm-7 packages"
  ./linux/fpm-package.sh linux-arm-7 deb
  ./linux/fpm-package.sh linux-arm-7 rpm
fi

fi

popd >/dev/null
