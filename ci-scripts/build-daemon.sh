#!/usr/bin/env bash

echo "start to build daemon..."
pushd "build"
if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then ./build-daemon-release.sh  'linux/amd64 linux/arm' ;fi
if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then ./build-daemon-release.sh 'darwin/amd64' ;fi
ls release/
popd
