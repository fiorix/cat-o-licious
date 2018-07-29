#!/bin/bash
# release-windows.sh produces the cat-o-licious zip containing an executable
# file and all required assets to for users to unpack and run the game.
#
# Requires a working Go compiler, mingw64, git, and zip utilities.
#
# Tested on Windows 10.

set -ex

VERSION=1.0
RELEASE=release/cat-o-licious-${VERSION}

rm -rf $RELEASE || true
trap "rm -rf $RELEASE" EXIT

mkdir -p $RELEASE

cp -r assets $RELEASE
go build -tags static -o $RELEASE/cat-o-licious.exe

APP=$(basename $RELEASE)
ZIP=$APP-windows.zip
(cd release && zip -r $ZIP $APP)
