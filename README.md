![daemon logo](https://user-images.githubusercontent.com/8619106/55698557-0bfc2400-59e4-11e9-972d-e33b640ff9a6.png)
 
# Hardware Wallet Daemon
[![Build Status](https://travis-ci.com/skycoin/hardware-wallet-daemon.svg)](https://travis-ci.com/skycoin/hardware-wallet-daemon)
[![GoDoc](https://godoc.org/github.com/skycoin/hardware-wallet-daemon?status.svg)](https://godoc.org/github.com/skycoin/hardware-walletd-daemon)
[![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/hardware-wallet-daemon)](https://goreportcard.com/report/github.com/skycoin/hardware-wallet-daemon)

The hardware walllet daemon provides an HTTP API to interface with the wallets supported by skycoin.
It uses the go bindings provided by the hardware wallet go [library](https://github.com/skycoin/hardware-wallet-go).

## Table of contents

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->
- [Installation](#installation)
	- [Go 1.10+ Installation and Setup](#go-110-installation-and-setup)
	- [Run Daemon from the command line](#run-daemon-from-the-command-line)
		- [Modes](#modes)
	- [Show Daemon options](#show-daemon-options)
- [API Documentation](#api-documentation)
	- [REST API](#rest-api)
- [Development guidelines](#development-guidelines)
	- [Client libraries](#client-libraries)
	- [Running tests](#running-tests)
	- [Running Integration Tests](#running-integration-tests)
		- [Emulator Integration Tests](#emulator-integration-tests)
		- [Wallet Integration Tests](#wallet-integration-tests)
		- [Debugging Integration Tests](#debugging-integration-tests)
		- [Update golden files in integration testdata](#update-golden-files-in-integration-testdata)
	- [Test coverage](#test-coverage)
	- [Formatting](#formatting)
	- [Code Linting](#code-linting)
	- [Profiling](#profiling)
	- [Dependency Management](#dependency-management)
	- [Releases](#releases)
		- [Update the version](#update-the-version)
		- [Pre-release testing](#pre-release-testing)
		- [Creating release builds](#creating-release-builds)
- [Responsible Disclosure](#responsible-disclosure)    
<!-- /MarkdownTOC -->

## Installation

Hardware Daemon supports go1.10+.

### Go 1.10+ Installation and Setup

[Golang 1.10+ Installation/Setup](https://github.com/skycoin/skycoin/blob/develop/INSTALLATION.md)

## Run Daemon from the command line
### Modes
The API has two modes:
1. **USB**: Communicate with the hardware wallet.
2. **EMULATOR**: Communicate with the emulator.

You can use the `-daemon-mode` flag to enable the required mode or use the `make` commands.

Example(USB Mode):
```sh
$ cd $GOPATH/src/github.com/skycoin/hardware-wallet-daemon
$ make run-usb
```

Example(Emulator Mode):
```sh
$ cd $GOPATH/src/github.com/skycoin/hardware-wallet-daemon
$ make run-emulator
```

### Show Daemon options

```sh
$ cd $GOPATH/src/github.com/skycoin/hardware-wallet-daemon
$ make run-help
```

## API Documentation


### REST API

[REST API](src/api/README.md).

# Development guidelines

Code added in this repository should comply to development guidelines documented in [Skycoin wiki](https://github.com/skycoin/skycoin/wiki).

The project has two branches: `master` and `develop`.

- `develop` is the default branch and will always have the latest code.
- `master` will always be equal to the current stable release on the website, and should correspond with the latest release tag.

### Client libraries

Hardware wallet daemon uses swagger for the API specification. It follows OpenAPI 2.0.

The go client has been automatically generated using [go-swagger](https://github.com/go-swagger/go-swagger) and the swagger specification.
>Note: The go client uses a slightly modified response [template](https://github.com/skycoin/hardware-wallet-daemon/blog/master/templates/client).

The swagger specification can be used to generate more libraries in the required language.

# Running tests

Library test suite can be run by running the following command

```sh
$ make test
```

### Running Integration Tests

There are integration tests for the HTTP API interfaces. They have two
run modes, "wallet" and "emulator".

The `emulator` integration tests are run only when an emulator is running.

The `wallet` integration tests are run only when a physical skycoin hardware wallet is connected.

#### Emulator Integration Tests

```sh
$ make test-integration-emulator
```

or

```sh
$ ./ci-scripts/integration-test.sh -m emulator -v
```

The `-m emulator` option, runs emulator integration tests.

The `-v` option, shows verbose logs.

#### Wallet Integration Tests

```sh
$ make test-integration-wallet
```

or

```sh
$ ./ci-scripts/integration-test.sh -m wallet -v
```

The `-m wallet` option, runs wallet integration tests.

The `-v` option, shows verbose logs.

#### Debugging Integration Tests

Run specific test case:

It's annoying and a waste of time to run all tests to see if the test we real care
is working correctly. There's an option: `-r`, which can be used to run specific test case.
For example: if we only want to test `TestEmulatorFeatures` and see the result, we can run:

```sh
$ ./ci-scripts/integration-test.sh -m emulator -v -r TestEmulatorFeatures
```


#### Update golden files in integration testdata

Golden files are expected data responses from the HTTP API saved to disk.
When some of the tests are run, their output is compared to the golden files.

To update golden files, use the provided `make` command:

**Wallet**
```bash
$ make DEVICE_TYPE=USB update-golden-files
```

**EMULATOR**
```bash
$ make DEVICE_TYPE=USB update-golden-files
```

We can also update a specific test case's golden file with the `-r` option.
For example:
```bash
$ ./ci-scripts/integration-test.sh -m emulator -v -u -r TestEmulatorFeatures
```

### Test coverage

Coverage is automatically generated for `make test` and integration tests run against a stable node.
This includes integration test coverage. The coverage output files are placed in `coverage/`.

To merge coverage from all tests into a single HTML file for viewing:

```bash
$ make check
$ make merge-coverage
```

Then open `coverage/all-coverage.html` in the browser.

### Formatting

All `.go` source files should be formatted `goimports`.  You can do this with:

```bash
$ make format
```

### Code Linting

Install prerequisites:

```bash
$ make install-linters
```

Run linters:

```bash
$ make lint
```

### Profiling

A full CPU profile of the program from start to finish can be obtained by running the node with the `-profile-cpu` flag.
Once the node terminates, a profile file is written to `-profile-cpu-file` (defaults to `cpu.prof`).
This profile can be analyzed with

```bash
go tool pprof cpu.prof
```

The HTTP interface for obtaining more profiling data or obtaining data while running can be enabled with `-http-prof`.
The HTTP profiling interface can be controlled with `-http-prof-host` and listens on `localhost:6060` by default.

See https://golang.org/pkg/net/http/pprof/ for guidance on using the HTTP profiler.

Some useful examples include:

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10
go tool pprof http://localhost:6060/debug/pprof/heap
```

A web page interface is provided by http/pprof at http://localhost:6060/debug/pprof/.

### Dependency Management

Dependencies are managed with [dep](https://github.com/golang/dep).

To [install `dep` for development](https://github.com/golang/dep/blob/master/docs/installation.md#development):

```bash
go get -u github.com/golang/dep/cmd/dep
```

`dep` vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `dep ensure`.

Use `dep help` for instructions on vendoring a specific version of a dependency, or updating them.

When updating or initializing, `dep` will find the latest version of a dependency that will compile.

Examples:

Initialize all dependencies:

```bash
dep init
```

Update all dependencies:

```bash
dep ensure -update -v
```

Add a single dependency (latest version):

```bash
dep ensure github.com/foo/bar
```

Add a single dependency (more specific version), or downgrade an existing dependency:

```bash
dep ensure github.com/foo/bar@tag
```

### Releases

#### Update the version

0. If the `master` branch has commits that are not in `develop` (e.g. due to a hotfix applied to `master`), merge `master` into `develop` (and fix any build or test failures)
0. Switch to a new release branch named `release-X.Y.Z` for preparing the release.
0. Update `CHANGELOG.md`: move the "unreleased" changes to the version and add the date.
0. Follow the steps in [pre-release testing](#pre-release-testing)
0. Make a PR merging the release branch into `master`
0. Review the PR and merge it
0. Tag the `master` branch with the version number. Version tags start with `v`, e.g. `v0.20.0`. Sign the tag. If you have your GPG key in github, creating a release on the Github website will automatically tag the release. It can be tagged from the command line with `git tag -as v0.20.0 $COMMIT_ID`, but Github will not recognize it as a "release".
0. Release builds are created and uploaded by travis. To do it manually, checkout the master branch and follow the [create release builds instructions](#creating-release-builds).
0. Checkout `develop` branch and bump `VERSION` to next [`dev` version number](https://www.python.org/dev/peps/pep-0440/#developmental-releases).

#### Pre-release testing

Pre-release testing procedure requires [skycoin-cli](https://github.com/skycoin/skycoin/tree/develop/cmd/cli). Please [install it](https://github.com/skycoin/skycoin/blob/develop/cmd/cli/README.md#install) if not available in your system. Some operations in the process require [running a Skycoin node](https://github.com/skycoin/skycoin/tree/master/INTEGRATION.md#running-the-skycoin-node). Also clone [Skywallet firmware repository](https://github.com/skycoin/hardware-wallet/) in advance.

The instructions that follow are meant to be followed for Skywallet devices flashed without memory protection. If your device memory is protected then some values might be different e.g. `firmwareFeatures`.

During the process beware of the fact that running an Skycoin node in the background can block the Skywallet from running.

Some values need to be known during the process. They are represented by the following variables:

- `WALLET1`, `WALLET2`, ... names of wallets created by `skycoin_cli`
- `ADDRESS1`, `ADDRESS2`, ... Skycoin addresses
- `TXN1_RAW`, `TXN2_RAW`, ... transactions data encoded in hex dump format
- `TXN1_JSON`, `TXN2_JSON`, ... transactions data encoded in JSON format, if numeric index value matches the one of another variable with `RAW` prefix then both refer to the same transaction
- `TXN1_ID`, `TXN2_ID`, ... Hash ID of transactions after broadcasting to the P2P network
- `AMOUNT` represents an arbitrary number of coins
- `ID1`, `ID2`, `ID3`, ... unique ID values , usually strings identifying hardware or software artifacts
- `AUTHTOKEN` is a CSRF token for the Skycoin REST API

Perform these actions before releasing.

**Note** : In all cases `skycoin-cli` would be equivalent to `go run cmd/cli/cli.go` if current working directory set to `$GOPATH/src/github.com/skycoin/skycoin`.

##### Run project test suite

Perform these actions before releasing:

- Open a terminal window and run Skywallet emulator. Wait for emulator UI to display.
- Connect SkyWallet device
- From a separate terminal window run the test suite as follows, and make sure that all the integration tests are passing.
```sh
make lint
make check
```

##### Run transaction tests

- Create new wallet e.g. with `skycoin-cli` (or reuse existing wallet for testing purposes)
```sh
skycoin-cli walletCreate -f $WALLET1.wlt -l $WALLET1
```
- From command output take note of the seed `SEED1` and address `ADDRESS1`
- List wallet addresses and confirm that `ADDRESS1` is the only value in the list.
```sh
skycoin-cli listAddresses $WALLET1.wlt
```
- Transfer funds to `ADDRESS1` in new wallet in two transactions
- Check balance
```sh
skycoin-cli addressBalance $ADDRESS1
```
- List unspent outputs for `WALLET1` and check in response that `head_outputs` includes two outputs with `address` set to `ADDRESS1`
```
skycoin-cli walletOutputs $WALLET1.wlt
```
- [Get device features](src/api/README.md#get-features) and check that:
  * `vendor` is set to `Skycoin Foundation`
  * `deviceId` is a string of 24 chars, which we'll refer to as `ID1`
  * write down the value of `bootloaderHash` i.e. `ID2`
  * `model` is set to `'1'`
  * `fwMajor` is set to expected firmware major version number
  * `fwMinor` is set to expected firmware minor version number
  * `fwPatch` is set to expected firmware patch version number
  * `firmwareFeatures` is set to `0`
- Ensure device is seedless by [wiping it](src/api/README.md#wipe). Check that device ends up in home screen with `NEEDS SEED!` message at the top.
- [Recover seed](src/api/README.md#recover-old-wallet) `SEED1` in Skywallet device (`dry-run=false`).
- [Get device features](src/api/README.md#get-features) and check that:
  * `vendor` is set to `Skycoin Foundation`
  * `deviceId` is set to `ID1`
  * `pinProtection` is seto to `false`
  * `passphraseProtection` is set to `false`
  * `label` is set to `ID1`
  * `initialized` is set to `true`
  * `bootloaderHash` is set to `ID2`
  * `passphraseCached` is set to `false`
  * `needsBackup` is set to `false`
  * `model` is set to `'1'`
  * `fwMajor` is set to expected firmware major version number
  * `fwMinor` is set to expected firmware minor version number
  * `fwPatch` is set to expected firmware patch version number
  * `firmwareFeatures` is set to `0`
- [Set device label](src/api/README.md#apply-settings) to a new value , say `ID3`. Specify `usePassphrase=false`.
- [Get device features](src/api/README.md#get-features) and check that:
  * `label` is set to `ID3`
  * all other values did not change with respect to previous step, especially `deviceId`
  * as a result device label is displayed on top of Skycoin logo in device home screen
- Ensure you know at least two addresses for test wallet, if not, [generate some](src/api/README.md#generate-address). Choose the second and third in order, hereinafter referred to as `ADDRESS2`, `ADDRESS3` respectively
- Check that address sequence generated by SkyWallet matches the values generated by `skycoin-cli`
```sh
skycoin-cli walletAddAddresses -f $WALLET1.wlt -n 5
```
- Check once again with desktop wallet
- Create new transaction from `ADDRESS1` to `ADDRESS2` in test wallet (say `TXN_RAW1`) for an spendable amount higher than individual output's coins
```sh
export TXN1_RAW="$(skycoin-cli createRawTransaction -a $ADDRESS1 -f $WALLET1.wlt $ADDRESS2 $AMOUNT)"
echo $TXN1_RAW
```
- Display transaction details and confirm that it contains at least two inputs
```sh
export TXN1_JSON=$(skycoin-cli decodeRawTransaction $TXN1_RAW)
echo $TXN1_JSON
```
- [Sign transaction](src/api/README.md#transaction-sign) with Skywallet by putting together a message using values resulting from previous step as follows.
  * Set message `nbIn` to the length of transaction `inputs` array
  * Set message `nbOut` to the length of transaction `outputs` array
  * For each hash in transaction `inputs` array there should be an item in messsage `inputs` array with `hashIn` field set to the very same hash and `index` set to `0`.
  * For each source item in transaction `outputs` array there should be an item in messsage `outputs` array with fields set as follows:
    - `address` : source item's `dst`
    - `coin` : source item's `coins`
    - `hour` : source item's `hours`
    - `address_index` : set to `0` if source item `address` equals `ADDRESS1` or to `1` otherwise
- Check that `signatures` array returned by hardware wallet includes entries for each and every transaction input
- [Check signatures](src/api/README.md#check-message-signature) were signed by corresponding addresses
- Create transaction `TXN2_JSON` by replacing `TXN1_JSON` signatures with the array returned by SkyWallet
- Use `TXN2_JSON` to obtain encoded transaction `TXN2_RAW`
```sh
export $TXN2_RAW=$( echo "$TXN2_JSON" | skycoin-cli encodeJsonTransaction - | grep '"rawtx"' | cut -d '"' -f4)
echo $TXN2_RAW
```
- Broadcast transaction. Refer to its id as `TXN2_ID`
```sh
export TXN2_ID=$(skycoin-cli broadcastTransaction $TXN2_RAW)
```
- After a a reasonable time check that balance changed.
```sh
skycoin-cli walletBalance $WALLET1.wlt
```
- Create a second wallet i.e. `WALLET2`
```sh
skycoin-cli walletCreate -f $WALLET2.wlt -l $WALLET2
```
- From command output take note of the seed `SEED2` and address `ADDRESS4`
- List `WALLET2` addresses and confirm that `ADDRESS4` is the only value in the list.
```sh
skycoin-cli lisAddresses $WALLET2.wlt
```
- Transfer funds to `WALLET2` (if not already done) and check balance
```sh
skycoin-cli addressBalance $ADDRESS4
```
- Request CSRF token (i.e. `AUTHTOKEN`) using Skycoin REST API.
```sh
curl http://127.0.0.1:6420/api/v1/csrf
```
- Use Skycoin REST API to create one transaction grabbing all funds from `ADDRESS1` (i.e. first address in `WALLET1` previously recovered in Skywallet device), `ADDRESS2` (i.e. second address in `WALLET1` previously recovered in Skywallet device) and `ADDRESS4` (i.e. first address in `WALLET2`) so as to transfer to the third address of `WALLET1` (i.e. `ADDRESS3`). Change should be sent back to `ADDRESS1`. In server response transaction JSON object (a.k.a `TNX3_JSON`) would be the object at `data.transaction` JSON path. If Skycoin node was started with default parameters this can be achieved as follows:
```sh
curl -X POST http://127.0.0.1:6420/api/v2/transaction -H 'content-type: application/json' -H "X-CSRF-Token: $AUTHTOKEN" -d "{
    \"hours_selection\": {
        \"type\": \"auto\",
        \"mode\": \"share\",
        \"share_factor\": \"0.5\"
    },
    \"addresses\": [\"$ADDRESS1\", \"$ADDRESS2\", \"$ADDRESS4\"],
    \"change_address\": \"$ADDRESS1\",
    \"to\": [{
        \"address\": \"$ADDRESS3\",
        \"coins\": \"$AMOUNT\"
    }],
    \"ignore_unconfirmed\": false
}"
```
- [Sign transaction](src/api/README.md#transaction-sign) represented by `TXN3_JSON` for inputs owned by Skywallet (i.e. `WALLET1`)
  * Set message `nbIn` to the 
  * Set message `nbOut` to the length of transaction `outputs` array
  * For each object in transaction `inputs` array there should be an item in messsage `inputs` array with:
    - `hashIn` field set to the value bound to object's `uxid` key
    - address index as follows
      * not set if input `address` is `ADDRESS4`
      * `0` if input `address` is `ADDRESS1`
      * `1` if input `address` is `ADDRESS2`
  * For each source item in transaction `outputs` array there should be an item in messsage `outputs` array with fields set as follows:
    - `address` : source item's `address`
    - `coin` : source item's `coins * 1000000`
    - `hour` : source item's `hours`
    - `address_index` set to 2 (since destination address is `ADDRESS3`)
- Check that signatures array includes one entry for every input except the one associated to `ADDRESS4`, which should be an empty string
- [Recover seed](src/api/README.md#recover-old-wallet) `SEED2` in Skywallet device (`dry-run=false`).
- [Sign transaction](src/api/README.md#transaction-sign) represented by `TXN3_JSON` for inputs owned by Skywallet (i.e. `WALLET2`)
  * Set message `nbIn` to the length of transaction `inputs` array
  * Set message `nbOut` to the length of transaction `outputs` array
  * For each hash in transaction `inputs` array there should be an item in messsage `inputs` array with:
    - `hashIn` field set to the value bound to object's `uxid` key
    - address index as follows
      * not set if input `address` is `ADDRESS1`
      * not set if input `address` is `ADDRESS2`
      * `0` if input `address` is `ADDRESS4`
  * For each source item in transaction `outputs` array there should be an item in messsage `outputs` array with fields set as follows:
    - `address` : source item's `dst`
    - `coin` : source item's `coins * 1000000`
    - `hour` : source item's `hours`
    - `address_index` set to 2 (since destination address is `ADDRESS3`)
- Create a new transaction JSON object (a.k.a `TXN4_JSON`) from `TXN3_JSON` and the previous signatures like this
  * `type` same as in `TXN3_JSON`
  * `inner_hash` should be an empty string
  * `sigs` returned by SkyWallet in same order as corresponding input
  * `inputs` is an array of strings. For each item in `TXN3_JSON` `inputs` include the value of its `uxid` field in `TXN4_JSON` `inputs` array. Respect original order.
  * `outputs` is an array of objects constructed out of `TXN3_JSON` `outputs` items, in te same order, as follows
    - `dst` : source item's `address`
    - `coins` : source item's `coins * 1000000`
    - `hours` : source item's `hours` as integer
- Use `TXN4_JSON` to obtain encoded transaction `TXN4_RAW`
```sh
export $TXN4_RAW=$( echo "$TXN4_JSON" | skycoin-cli encodeJsonTransaction - | grep '"rawtx"' | cut -d '"' -f4)
echo $TXN4_RAW
```
- Broadcast transaction. Refer to its id as `TXN4_ID`
```sh
export TXN4_ID=$(skycoin-cli broadcastTransaction $TXN4_RAW)
```
- After a a reasonable time check that wallets balance changed.
```sh
skycoin-cli walletBalance $WALLET1.wlt
skycoin-cli walletBalance $WALLET2.wlt
```

#### Creating release builds

The following instruction creates a full release:

```bash
make release
```
Daemon version will be take from `VERSION`.

## Responsible Disclosure

Security flaws in the source of any software provided by skycoin or infrastructure can be sent to security@skycoin.net.
Bounties are available for accepted critical bug reports.

PGP Key for signing:

```
-----BEGIN PGP PUBLIC KEY BLOCK-----

mDMEWaj46RYJKwYBBAHaRw8BAQdApB44Kgde4Kiax3M9Ta+QbzKQQPoUHYP51fhN
1XTSbRi0I0daLUMgU0tZQ09JTiA8dG9rZW5AcHJvdG9ubWFpbC5jb20+iJYEExYK
AD4CGwMFCwkIBwIGFQgJCgsCBBYCAwECHgECF4AWIQQQpyK3by/+e9I4AiJYAWMb
0nx4dAUCWq/TNwUJCmzbzgAKCRBYAWMb0nx4dKzqAP4tKJIk1vV2bO60nYdEuFB8
FAgb5ITlkj9PyoXcunETVAEAhigo4miyE/nmE9JT3Q/ZAB40YXS6w3hWSl3YOF1P
VQq4OARZqPjpEgorBgEEAZdVAQUBAQdAa8NkEMxo0dr2x9PlNjTZ6/gGwhaf5OEG
t2sLnPtYxlcDAQgHiH4EGBYKACYCGwwWIQQQpyK3by/+e9I4AiJYAWMb0nx4dAUC
Wq/TTQUJCmzb5AAKCRBYAWMb0nx4dFPAAQD7otGsKbV70UopH+Xdq0CDTzWRbaGw
FAoZLIZRcFv8zwD/Z3i9NjKJ8+LS5oc8rn8yNx8xRS+8iXKQq55bDmz7Igw=
=5fwW
-----END PGP PUBLIC KEY BLOCK-----
```

Key ID: [0x5801631BD27C7874](https://pgp.mit.edu/pks/lookup?search=0x5801631BD27C7874&op=index)

The fingerprint for this key is:

```
pub   ed25519 2017-09-01 [SC] [expires: 2023-03-18]
      10A7 22B7 6F2F FE7B D238  0222 5801 631B D27C 7874
uid                      GZ-C SKYCOIN <token@protonmail.com>
sub   cv25519 2017-09-01 [E] [expires: 2023-03-18]
```

Keybase.io account: https://keybase.io/gzc
