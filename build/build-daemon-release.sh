#!/usr/bin/env bash
set -e -o pipefail

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

./gox.sh "$GOX_OSARCH" "$GOX_DMN_OUTPUT_DIR"

echo "==========================="
echo "Packaging daemon release"
./package-daemon-release.sh "$GOX_OSARCH"

echo "------------------------------"
echo "Compressing daemon release"
./compress-daemon-release.sh "$GOX_OSARCH"

popd >/dev/null
