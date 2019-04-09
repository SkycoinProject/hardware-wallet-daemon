#!/bin/bash

#Set Script Name variabl
SCRIPT=`basename ${BASH_SOURCE[0]}`

PORT=9510
HOST="http://127.0.0.1:$PORT"
MODE="emulator"
UPDATE=""
TIMEOUT="60m"
# run go test with -v flag
VERBOSE=""
# run go test with -run flag
RUN_TESTS=""
FAILFAST=""
NAME=""
USE_CSRF=""
ENABLE_CSRF=""

usage () {
  echo "Usage: $SCRIPT"
  echo "Optional command line arguments"
  echo "-m <string>  -- Testmode to run, emulator or wallet;"
  echo "-r <string>  -- Run test with -run flag"
  echo "-n <string>  -- Specific name for this test, affects coverage output files"
  echo "-u <boolean> -- Update testdata"
  echo "-v <boolean> -- Run test with -v flag"
  echo "-f <boolean> -- Run test with -failfast flag"
  echo "-c <boolean> -- Pass this argument if the node has CSRF enabled"
  exit 1
}

while getopts "h?m:r:n:uvfc" args; do
case $args in
    h|\?)
        usage;
        exit;;
    m ) MODE=${OPTARG};;
    r ) RUN_TESTS="-run ${OPTARG}";;
    n ) NAME="-${OPTARG}";;
    u ) UPDATE="--update";;
    v ) VERBOSE="-v";;
    f ) FAILFAST="-failfast";;
    c ) ENABLE_CSRF="-enable-csrf"; USE_CSRF="1";;
  esac
done

BINARY="daemon-integration${NAME}.test"

COVERAGEFILE="coverage/${BINARY}.coverage.out"
if [ -f "${COVERAGEFILE}" ]; then
    rm "${COVERAGEFILE}"
fi

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
CMDPKG=$(go list ./cmd/daemon)
COVERPKG=$(dirname $(dirname ${CMDPKG}))
GOLDFLAGS="-X ${CMDPKG}.Commit=${COMMIT} -X ${CMDPKG}.Branch=${BRANCH}"

set -euxo pipefail

DATA_DIR=$(mktemp -d -t daemon-data-dir.XXXXXX)

# Compile the daemon
# We can't use "go run" because that creates two processes which doesn't allow us to kill it at the end
echo "compiling daemon with coverage"
go test -c -ldflags "${GOLDFLAGS}" -tags testrunmain -o "$BINARY" -coverpkg="${COVERPKG}/..." ./cmd/daemon/

mkdir -p coverage/

# Run daemon
echo "starting daemon node in background with http listener on $HOST"

./"$BINARY" -web-interface-port=$PORT \
            -data-dir="$DATA_DIR" \
            -test.run "^TestRunMain$" \
            -test.coverprofile="${COVERAGEFILE}" \
            $ENABLE_CSRF \
            &

DAEMON_PID=$!

echo "daemon pid=$DAEMON_PID"

set +e

HW_DAEMON_INTEGRATION_TESTS=1 HW_DAEMON_INTEGRATION_TEST_MODE=$MODE USE_CSRF=$USE_CSRF \
    go test ./src/api/integration/... $FAILFAST $UPDATE -timeout=$TIMEOUT $VERBOSE $RUN_TESTS

TEST_FAIL=$?

echo "shutting down daemon"

# Shutdown daemon
kill -s SIGINT $DAEMON_PID
wait $DAEMON_PID

rm "$BINARY"

exit $TEST_FAIL
