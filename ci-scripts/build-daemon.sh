#!/usr/bin/env bash
set -e -o pipefail

echo "start to build daemon..."
pushd "build"

if [[ "$OSTYPE" == "linux"* ]]; then ./build-daemon-release.sh  'windows/386,windows/amd64,linux/386,linux/amd64,linux/arm-7' ;fi

if [[ "$OSTYPE" == "darwin"* ]]; then
	echo "load osx config"
	. build-conf.sh "darwin-10.10/amd64"
	OSX64="${DMN_OUTPUT_DIR}/${OSX64_DMN}"

	echo "make output directories"
	rm -rf "$OSX64"
	mkdir -p "$OSX64"
	rm -rf release
	mkdir -p release

	echo "build daemon ${OSX64}/${BIN_NAME}"
	go build -o "${OSX64}/${BIN_NAME}" ../cmd/daemon/daemon.go
	rm -rf osx/build

	echo "create osx/build/ directory"
	mkdir -p osx/build

	echo "copy binary file to osx/build/"
	cp "${OSX64}/${BIN_NAME}" osx/build/

	echo "set version: ${APP_VERSION}"
	echo "${APP_VERSION}" > osx/build/VERSION
	./osx/release.sh
	# cp osx/build/*.pkg release/
	PKG=$(ls osx/build/*.pkg)
	INSTALLER=$(basename "${PKG}")

	echo "signing binary"
	echo "${INSTALLER}"
	productsign --sign "Developer ID Installer: yunfei mao" "osx/build/${INSTALLER}" "release/${INSTALLER}"

	echo "clear temporary builds"
	rm -rf osx/build
fi

ls release/
popd >/dev/null
