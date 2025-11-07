# Keeper Installer - Deployment Guide

## Current Status

✅ **GitHub Repository:** https://github.com/qnepff/keeper-installer
✅ **GitHub Actions:** Building now (Linux, Windows, macOS)
✅ **Website Updated:** /tmp/quickneasy.com/index.html (ready to upload)
✅ **Upload Script:** upload-to-cdn.sh (ready to use)

## Step 1: Wait for GitHub Actions Build

**Check build status:**
https://github.com/qnepff/keeper-installer/actions

Build takes ~5-10 minutes. You'll get 3 artifacts:
- keeper-installer-linux-amd64
- keeper-installer-windows-amd64.exe
- keeper-installer-darwin-amd64

## Step 2: Download Build Artifacts

1. Go to the completed workflow run
2. Scroll to "Artifacts" section
3. Download each artifact:
   - Click "keeper-installer-linux-amd64" → saves zip
   - Click "keeper-installer-windows-amd64.exe" → saves zip
   - Click "keeper-installer-darwin-amd64" → saves zip
4. Extract to ~/Downloads/

```bash
cd ~/Downloads
unzip keeper-installer-linux-amd64.zip
unzip keeper-installer-windows-amd64.exe.zip
unzip keeper-installer-darwin-amd64.zip
```

## Step 3: Upload to BunnyCDN

### Option A: Using upload-to-cdn.sh script

```bash
cd /home/paulf/qne-repos/keeper-installer

# Set API key
export BUNNY_API_KEY=your-bunny-api-key

# Run upload script
./upload-to-cdn.sh
```

### Option B: Manual upload with bunnycdn-cli

```bash
export BUNNY_API_KEY=your-bunny-api-key

cd /home/paulf/cascade-chameleon/bunnycdnutils

./bunnycdn-cli upload -zone qne-installers-cdn \
  -files ~/Downloads/keeper-installer-linux-amd64

./bunnycdn-cli upload -zone qne-installers-cdn \
  -files ~/Downloads/keeper-installer-windows-amd64.exe

./bunnycdn-cli upload -zone qne-installers-cdn \
  -files ~/Downloads/keeper-installer-darwin-amd64
```

## Step 4: Upload Updated Website

```bash
export BUNNY_API_KEY=your-bunny-api-key

cd /home/paulf/cascade-chameleon/bunnycdnutils

./bunnycdn-cli upload -zone quickneasy-info \
  -files /tmp/quickneasy.com/index.html
```

## Step 5: Test Downloads

Visit https://quickneasy-info-cdn.b-cdn.net and test:

1. **Linux:** Button should download keeper-installer-linux-amd64
2. **Windows:** Visit from Windows, should get .exe
3. **macOS:** Visit from Mac, should get darwin version

## Step 6: DNS Switch (After Nov 9 Transfer)

Once quickneasy.info transfers to Namecheap:

```bash
export NAMECHEAP_API_USER=qnecommunity
export NAMECHEAP_API_KEY=884ba3b18adb4e69a0ee5dd5945653f4

cd /home/paulf/cascade-chameleon/bunnycdnutils
./bunnycdn-cli switch-ns -domain quickneasy.info
```

Wait 5-30 minutes for DNS propagation, then:
- Site accessible at https://quickneasy.info ✓

## Verification Checklist

- [ ] GitHub Actions build completed successfully
- [ ] All 3 platform installers downloaded
- [ ] Installers uploaded to BunnyCDN
- [ ] CDN URLs accessible:
  - [ ] https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-linux-amd64
  - [ ] https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-windows-amd64.exe
  - [ ] https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-darwin-amd64
- [ ] Updated website uploaded to BunnyCDN
- [ ] Website accessible at https://quickneasy-info-cdn.b-cdn.net
- [ ] Download button works and detects OS correctly
- [ ] Test Linux installer end-to-end (download → install → run)

## Current CDN URLs

**Installers:**
- Linux: `https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-linux-amd64`
- Windows: `https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-windows-amd64.exe`
- macOS: `https://qne-installers-cdn-cdn.b-cdn.net/keeper-installer-darwin-amd64`

**Keeper AppImage** (what installer downloads):
- `https://qne-installers-cdn.b-cdn.net/keeper-1.0.0.AppImage`
- ⚠️ Note: This needs to be built from qne-browser → keeper rebrand (separate task)

**Website:**
- CDN: `https://quickneasy-info-cdn.b-cdn.net`
- Domain: `https://quickneasy.info` (after DNS)

## Next Major Task: Keeper Application

After installer is fully deployed, the qne-browser Electron app needs rebranding:

1. Copy /home/paulf/cascade-chameleon/qne-browser → /home/paulf/qne-repos/keeper
2. Rebrand all UI text "QNE Browser" → "Keeper"
3. Update package.json, build configs
4. Build new AppImage: keeper-1.0.0.AppImage
5. Upload to CDN
6. Installer will then download the correctly named Keeper

## Support

For issues:
- GitHub: https://github.com/qnepff/keeper-installer/issues
- Progress tracking: /home/paulf/qne-repos/PROGRESS.md
