#!/bin/bash

set -e

APP_DIR="dist/macos/Homepage.app"
DMG_PATH="dist/macos/Homepage.dmg"
DMG_DIR="dist/macos/dmg"

CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
HELPERS_DIR="$CONTENTS_DIR/Helpers"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

rm -rf "$APP_DIR" "$DMG_PATH" "$DMG_DIR"
mkdir -p "$MACOS_DIR" "$HELPERS_DIR" "$RESOURCES_DIR"

GOOS=darwin go build -o "$HELPERS_DIR/homepage-server" .
cp AppIcon.png "$RESOURCES_DIR/AppIcon.png"

cat > "$MACOS_DIR/HomepageLauncher" <<'EOF'
#!/bin/bash

set -e

SERVER_PATH="\${0%/*}/../Helpers/homepage-server"
PLIST_PATH="\$HOME/Library/LaunchAgents/com.nielsprins.homepage.agent.plist"

mkdir -p "\${PLIST_PATH%/*}"

cat > "\$PLIST_PATH" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.nielsprins.homepage.agent</string>
  <key>ProgramArguments</key>
  <array>
    <string>\$SERVER_PATH</string>
  </array>
  <key>KeepAlive</key>
  <true/>
  <key>RunAtLoad</key>
  <true/>
</dict>
</plist>
PLIST

launchctl bootout "gui/\$UID" "\$PLIST_PATH" >/dev/null 2>&1 || true
launchctl bootstrap "gui/\$UID" "\$PLIST_PATH"
EOF

cat > "$CONTENTS_DIR/Info.plist" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleExecutable</key>
  <string>HomepageLauncher</string>
  <key>CFBundleIconFile</key>
  <string>AppIcon.png</string>
  <key>CFBundleIdentifier</key>
  <string>com.nielsprins.homepage.launcher</string>
  <key>CFBundleName</key>
  <string>Homepage</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleShortVersionString</key>
  <string>1.0</string>
  <key>CFBundleVersion</key>
  <string>1</string>
  <key>LSMinimumSystemVersion</key>
  <string>13.0</string>
  <key>LSUIElement</key>
  <true/>
</dict>
</plist>
EOF

chmod +x "$MACOS_DIR/HomepageLauncher" "$HELPERS_DIR/homepage-server"

mkdir -p "$DMG_DIR"
cp -R "$APP_DIR" "$DMG_DIR/Homepage.app"
ln -s /Applications "$DMG_DIR/Applications"

hdiutil create \
  -volname "Homepage" \
  -srcfolder "$DMG_DIR" \
  -ov \
  -format UDZO \
  "$DMG_PATH" >/dev/null

printf "Built:\n  %s\n  %s\n" "$APP_DIR" "$DMG_PATH"
