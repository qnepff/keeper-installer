# Assets Directory

## Required Files

### qne-browser.png
Place your QNE Browser icon here (256x256 PNG recommended).

This icon will be installed to: `/usr/share/icons/hicolor/256x256/apps/qne-browser-installer.png`

You can create a simple placeholder with ImageMagick:
```bash
convert -size 256x256 xc:blue -fill white -pointsize 72 -gravity center \
  -annotate +0+0 "QNE" qne-browser.png
```

Or use your actual QNE Browser logo/icon.

### qne-browser-installer.desktop
Desktop entry file for the installer (already created).

This creates the application menu entry on Linux systems.
