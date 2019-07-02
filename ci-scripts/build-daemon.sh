#!/usr/bin/env bash
set -e -o pipefail

echo "start to build daemon..."
pushd "build"


if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then ./build-daemon-release.sh  'windows/386,windows/amd64,linux/arm,linux/amd64' ;fi
if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    echo "load osx config"
    . build-conf.sh "darwin-10.10/amd64"
    OSX64="${DMN_OUTPUT_DIR}/${OSX64_DMN}"
    echo "make output directories"
    mkdir -p "$OSX64"
    mkdir -p release
    echo "build daemon ${OSX64}/${BIN_NAME}"
    go build -o "${OSX64}/${BIN_NAME}" ../cmd/daemon/daemon.go
    echo "signing binary"
    codesign --force --sign "Developer ID Application: yunfei mao" "${OSX64}/${BIN_NAME}"
    cp ../VERSION "$OSX64"
    echo "------------------------------"
    echo "Compressing daemon release"
    ./compress-daemon-release.sh "darwin-10.10/amd64"
fi

ls release/
popd >/dev/null
