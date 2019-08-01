#!/bin/bash

set -e -o pipefail

cd $(dirname $0)

VERSION=$(cat build/VERSION)

INSTALLER=skyhwd-$VERSION-osx-darwin-x64.pkg

# make uninstaller
rm -rf build/uninstall
cp -r uninstall build
pushd "build" >/dev/null
pkgbuild --identifier net.skycoin.skyhwd --version=$VERSION --root uninstall/ROOT --scripts uninstall/scripts uninstall.pkg
rm -rf uninstall

# make installer
rm -rf install
cp -r ../install .
cp skyhwd install/ROOT/Applications/Utilities/SKYHWD/
mv uninstall.pkg install/ROOT/Applications/Utilities/SKYHWD/
pkgbuild --identifier net.skycoin.skyhwd --version=$VERSION --root install/ROOT --scripts install/scripts $INSTALLER
rm -rf install
popd >/dev/null
