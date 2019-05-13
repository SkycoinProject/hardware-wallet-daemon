#!/usr/bin/env bash

if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then

echo "start to build daemon..."
pushd "build"
./build-daemon-release.sh  'darwin-10.10/amd64,windows/386,windows/amd64,linux/arm,linux/amd64' ;
ls release/
popd

fi

