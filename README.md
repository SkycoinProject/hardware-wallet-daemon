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

Performs these actions before releasing:

* `make lint`
* `make check` (make sure that all the integration tests are passing)

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
