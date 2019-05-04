#!/usr/bin/env bash

echo "start to build daemon..."
pushd "build"
if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then ./build-daemon-release.sh  'darwin/amd64,windows/386,windows/amd64,linux/386,linux/amd64' ;fi
ls release/
popd
