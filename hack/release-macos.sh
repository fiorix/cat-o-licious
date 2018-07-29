#!/bin/bash
# release-macos.sh produces the cat-o-licious zip containing a mac bundle
# with all required assets to for users to unpack and run the game.
#
# Requires a working Go compiler.
#
# Tested on macOS High Sierra version 10.13.6.

set -ex

VERSION=1.0
RELEASE=release/cat-o-licious.app
ICONSET=release/cat-o-licious.iconset
ICON=assets/img/player_frame_3.png

rm -rf $ICONSET $RELEASE || true
trap "rm -rf $ICONSET $RELEASE" EXIT

mkdir -p $ICONSET $RELEASE/Contents/{MacOS,Resources}

for SIZE in 16 32 64 128 256 512 1024
do
    sips -z $SIZE $SIZE $ICON \
        --out $ICONSET/icon_${SIZE}x${SIZE}.png

    [ $SIZE -eq 1024 ] && continue

    SIZE2X=$((SIZE+SIZE))

    sips -z $SIZE2X $SIZE2X $ICON --out $ICONSET/icon_${SIZE}x${SIZE}@2x.png
done

iconutil --convert icns --output $RELEASE/Contents/Resources/icon.icns $ICONSET
cp -r assets $RELEASE/Contents/Resources
go build -tags static -o $RELEASE/Contents/Resources/cat-o-licious

cat << EOF > $RELEASE/Contents/Info.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>cat-o-licious</string>
    <key>CFBundleIconFile</key>
    <string>icon</string>
    <key>CFBundleIdentifier</key>
    <string>com.github.fiorix.cat-o-licious</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSUIElement</key>
    <false/>
</dict>
</plist>
EOF

cat << EOF > $RELEASE/Contents/MacOS/cat-o-licious
#!/bin/bash
cd \$(dirname \$0)/../Resources
exec ./cat-o-licious
EOF

chmod +x $RELEASE/Contents/MacOS/cat-o-licious

APP=$(basename $RELEASE)
ZIP=$(basename $RELEASE | sed 's/\.app//g')-${VERSION}-macos.zip
(cd release && zip -r $ZIP $APP)
