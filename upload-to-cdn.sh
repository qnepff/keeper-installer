#!/bin/bash
# Upload keeper installer binaries to BunnyCDN

set -e

BUNNYCDN_CLI="/home/paulf/cascade-chameleon/bunnycdnutils/bunnycdn-cli"
STORAGE_ZONE="qne-installers-cdn"
INSTALLER_DIR="$HOME/Downloads"  # Where GitHub Actions artifacts are downloaded

echo "🚀 Uploading Keeper Installers to BunnyCDN"
echo "============================================"
echo ""

# Check if bunnycdn-cli exists
if [ ! -f "$BUNNYCDN_CLI" ]; then
    echo "❌ Error: bunnycdn-cli not found at $BUNNYCDN_CLI"
    exit 1
fi

# Check if BUNNY_API_KEY is set
if [ -z "$BUNNY_API_KEY" ]; then
    echo "❌ Error: BUNNY_API_KEY environment variable not set"
    exit 1
fi

echo "📁 Looking for installer files in: $INSTALLER_DIR"
echo ""

# Upload Linux installer
if [ -f "$INSTALLER_DIR/keeper-installer-linux-amd64" ]; then
    echo "📤 Uploading Linux installer..."
    $BUNNYCDN_CLI upload -zone $STORAGE_ZONE -files "$INSTALLER_DIR/keeper-installer-linux-amd64"
    echo "✅ Linux installer uploaded"
else
    echo "⚠️  Linux installer not found"
fi

echo ""

# Upload Windows installer
if [ -f "$INSTALLER_DIR/keeper-installer-windows-amd64.exe" ]; then
    echo "📤 Uploading Windows installer..."
    $BUNNYCDN_CLI upload -zone $STORAGE_ZONE -files "$INSTALLER_DIR/keeper-installer-windows-amd64.exe"
    echo "✅ Windows installer uploaded"
else
    echo "⚠️  Windows installer not found"
fi

echo ""

# Upload macOS installer
if [ -f "$INSTALLER_DIR/keeper-installer-darwin-amd64" ]; then
    echo "📤 Uploading macOS installer..."
    $BUNNYCDN_CLI upload -zone $STORAGE_ZONE -files "$INSTALLER_DIR/keeper-installer-darwin-amd64"
    echo "✅ macOS installer uploaded"
else
    echo "⚠️  macOS installer not found"
fi

echo ""
echo "============================================"
echo "✅ Upload complete!"
echo ""
echo "CDN URLs:"
echo "  Linux:   https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-linux-amd64"
echo "  Windows: https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-windows-amd64.exe"
echo "  macOS:   https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-darwin-amd64"
echo ""
echo "Next: Update quickneasy.info website with new URLs"
