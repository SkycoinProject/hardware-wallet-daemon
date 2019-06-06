![hwgo-01](https://user-images.githubusercontent.com/8619106/57969063-e6255000-798f-11e9-973e-9fdf7a0bd5fc.png)

# Go bindings and CLI tool for the Skycoin hardware wallet

[![Build Status](https://travis-ci.com/skycoin/hardware-wallet-go.svg?branch=master)](https://travis-ci.com/skycoin/hardware-wallet-go)

## Table of contents

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->
- [Installation](#installation)
- [Usage](#usage)
  - [Download source code](#download-source-code)
  - [Dependancies management](#dependancies-management)
  - [Run](#run)
- [Development guidelines](#development-guidelines)
  - [Versioning policies](#versioning-policies)
  - [Running tests](#running-tests)
  - [Releases](#releases)
    - [Update the version](#update-the-version)
    - [Pre-release testing](#pre-release-testing)
    - [Creating release builds](#creating-release-builds)
- [Wiki](#wiki)
<!-- /MarkdownTOC -->

## Installation

### Install golang

https://github.com/golang/go/wiki/Ubuntu

## Usage

### Download source code

```bash
$ go get github.com/skycoin/hardware-wallet-go
```

### Dependancies management

This project uses dep [dependancy manager](https://github.com/golang/dep).

Don't modify anything under vendor/ directory without using [dep commands](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md).

Download dependencies using command:

```bash
$ make dep
```

### Run

```bash
$ go run cmd/cli/cli.go
```

See also [CLI README](https://github.com/skycoin/hardware-wallet-go/blob/master/cmd/cli/README.md) for information about the Command Line Interface.

# Development guidelines

Code added in this repository should comply to development guidelines documented in [Skycoin wiki](https://github.com/skycoin/skycoin/wiki).

The project has two branches: `master` and `develop`.

- `develop` is the default branch and will always have the latest code.
- `master` will always be equal to the current stable release on the website, and should correspond with the latest release tag.

# Versioning policies

This go client library follows the [versioning rules for SkyWallet client libraries](https://github.com/skycoin/hardware-wallet/tree/develop/README.md#versioning-libraries) .

# Running tests

Library test suite can be run by running the following command

```
make test
```

If a physical SkyWallet device is connected then the test suite will exchange test messages with it. Likewise, the emulator will reply to test messages if executed in advance as follows.

```
git clone https://github.com/skycoin/hardware-wallet /some/path/to/hardware-wallet
make clean -C /some/path/to/hardware-wallet && make -C /some/path/to/hardware-wallet run-emulator
```

If neither the emulator nor a physical device are connected then tests will be skipped silently.

# Releases

# Update the version

0. If the `master` branch has commits that are not in `develop` (e.g. due to a hotfix applied to `master`), merge `master` into `develop` (and fix any build or test failures)
0. Switch to a new release branch named `release-X.Y.Z` for preparing the release.
0. Run `make build` to make sure that the code base is up to date
0. Update `CHANGELOG.md`: move the "unreleased" changes to the version and add the date.
0. Follow the steps in [pre-release testing](#pre-release-testing)
0. Make a PR merging the release branch into `master`
0. Review the PR and merge it
0. Tag the `master` branch with the version number. Version tags start with `v`, e.g. `v0.20.0`. Sign the tag. If you have your GPG key in github, creating a release on the Github website will automatically tag the release. It can be tagged from the command line with `git tag -as v0.20.0 $COMMIT_ID`, but Github will not recognize it as a "release".
0. Tag the changeset of the `protob` submodule checkout with the same version number as above.
0. Release builds are created and uploaded by travis. To do it manually, checkout the master branch and follow the [create release builds instructions](#creating-release-builds).
0. Checkout `develop` branch and bump `tiny-firmware/VERSION` and `tiny-firmware/bootloader/VERSION` to next [`dev` version number](https://www.python.org/dev/peps/pep-0440/#developmental-releases).


# Pre-release testing

Follow [hardware wallet firmware pre-release instructions](https://github.com/skycoin/hardware-wallet/tree/develop/README.md#pre-release-testing) by using library commands.

## Wiki

More information in [the wiki](https://github.com/skycoin/hardware-wallet-go/wiki)
