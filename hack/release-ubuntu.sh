#!/bin/bash
# release-ubuntu.sh produces a cat-o-licious debian package containing the
# game binary and assets.
#
# Tested on Ubuntu 18.04 LTS, ubuntu/bionic64 on vagrant.

set -ex

VERSION=1.0-1
RELEASE=release/cat-o-licious_${VERSION}

rm -rf pkgs $RELEASE || true
trap "rm -rf pkgs $RELEASE" EXIT

PKGS=" \
    golang \
    libasound2-dev \
    libx11-dev \
    libxcursor-dev \
    libxext-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxss-dev \
    libxxf86vm-dev \
    zlib1g-dev \
"

sudo apt-get install -y $PKGS

DEPS=""
for PKG in $PKGS
do
    [ "$(echo $PKG | grep -e '-dev')" == "" ] && continue

    IFS=',' read -r -a array <<< $(dpkg -s $PKG | grep Depends: | sed 's/Depends: //g')
    for element in "${array[@]}"
    do
        [ "$(echo $element | grep -e '\(-dev\|-doc\)')" != "" ] && continue
        [ "$DEPS" != "" ] && DEPS="${DEPS},"
        DEPS="$DEPS $element"
    done
done

mkdir -p $RELEASE/{DEBIAN,usr/bin,usr/games/cat-o-licious}
cp -r assets $RELEASE/usr/games/cat-o-licious
go build -tags static -o $RELEASE/usr/games/cat-o-licious/cat-o-licious

cat << EOF > $RELEASE/DEBIAN/control
Package: cat-o-licious
Version: $VERSION
Section: base
Priority: optional
Architecture: amd64
Depends: $DEPS
Maintainer: Alexandre Fiori <fiorix@gmail.com>
Description: Simple cat game written in Go and SDL.
 See https://github.com/fiorix/cat-o-licious for details.
EOF

cat << EOF > $RELEASE/usr/bin/cat-o-licious
#!/bin/bash
cd /usr/games/cat-o-licious
exec ./cat-o-licious
EOF

chmod +x $RELEASE/usr/bin/cat-o-licious

dpkg-deb --build $RELEASE