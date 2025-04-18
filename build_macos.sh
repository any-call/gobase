#!/bin/bash

set -e

APP_NAME="$1"
BIN_PATH="$2"
ICON_FILE="$3"

# 检查必填参数
if [ -z "$APP_NAME" ] || [ -z "$BIN_PATH" ]; then
  echo "❗ 错误：参数不足！"
  echo ""
  echo "用法：./build_app.sh <AppName> <BinaryPath> [Icon.icns]"
  echo "例如：./build_app.sh ClearLight bin/clearlight icons/icon.icns"
  echo ""
  echo "参数说明："
  echo "  AppName       应用显示名称（也是 .app 的名字）"
  echo "  BinaryPath    已编译好的可执行文件路径"
  echo "  Icon.icns     可选图标文件（.icns 格式）"
  echo ""
  exit 1
fi

# 检查可执行文件是否存在
if [ ! -f "$BIN_PATH" ]; then
  echo "❗ 错误：可执行文件 '$BIN_PATH' 不存在！"
  exit 1
fi

# 处理变量
BIN_NAME=$(basename "$BIN_PATH")
VERSION="1.0"
IDENTIFIER="com.yourcompany.$(echo "$APP_NAME" | tr '[:upper:]' '[:lower:]')"

echo "📦 正在打包：$APP_NAME.app"
rm -rf "$APP_NAME.app"
mkdir -p "$APP_NAME.app/Contents/MacOS"
mkdir -p "$APP_NAME.app/Contents/Resources"

# 拷贝可执行文件
cp "$BIN_PATH" "$APP_NAME.app/Contents/MacOS/"

# 处理图标
ICON_PLIST_ENTRY=""
if [ -n "$ICON_FILE" ]; then
  if [ ! -f "$ICON_FILE" ]; then
    echo "⚠️  警告：图标文件 '$ICON_FILE' 不存在，跳过设置图标。"
  else
    cp "$ICON_FILE" "$APP_NAME.app/Contents/Resources/"
    ICON_NAME=$(basename "$ICON_FILE")
    ICON_NAME_NO_EXT="${ICON_NAME%.*}"
    ICON_PLIST_ENTRY=$(printf "  <key>CFBundleIconFile</key>\n  <string>${ICON_NAME_NO_EXT}</string>\n")
    echo "🎨 使用图标: $ICON_FILE"
  fi
fi

# 写 Info.plist
cat > "$APP_NAME.app/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
 "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key>
  <string>$APP_NAME</string>
  <key>CFBundleExecutable</key>
  <string>$BIN_NAME</string>
  <key>CFBundleIdentifier</key>
  <string>$IDENTIFIER</string>
  <key>CFBundleVersion</key>
  <string>$VERSION</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>LSUIElement</key>
  <true/>
${ICON_PLIST_ENTRY}
</dict>
</plist>
EOF

echo "✅ 成功生成 $APP_NAME.app！你可以双击它来运行。"