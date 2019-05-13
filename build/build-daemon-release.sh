#!/usr/bin/env bash
set -e -o pipefail

XGO_TARGETS="$@"

. build-conf.sh "$XGO_TARGETS"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

GOPATH=~/go ./xgo.sh "$XGO_TARGETS" "$XGO_DMN_OUTPUT_DIR"

echo "==========================="
echo "Packaging daemon release"
./package-daemon-release.sh "$XGO_TARGETS"

echo "------------------------------"
echo "Compressing daemon release"
./compress-daemon-release.sh "$XGO_TARGETS"

popd >/dev/null
