# Keeper Installer

Installer for **Chameleon Keeper** - Personal P2P Platform & Skills Manager

**Platform Support**: Linux (Windows/macOS coming soon - requires native Keeper builds)

## About Chameleon Keeper

Chameleon Keeper is your platform for managing chameleons and accessing the decentralized QNE Network. It provides:

- **P2P Web Browser** - Browse peer content using custom P2P technology
- **Skills Platform** - Email, messaging, file sharing, and more
- **Personal Agent** - Protects your privacy and manages your identity
- **QNE Network** - Closed, members-only network with no passwords, spam, or ads

## Features

- ✅ **Self-installing** - Copies itself to `~/QNE/local` on first run
- ✅ **No sudo required** - User-level installation
- ✅ **GUI installer** - Built with Fyne for a native experience
- ✅ **Linux support** - Works on any modern Linux distribution
- ✅ **Auto-downloads** - Fetches Keeper AppImage from CDN
- ✅ **Creates shortcuts** - Desktop and app menu entries
- ✅ **Uninstaller** - Clean removal of all files

## Installation

### Linux

1. Download the installer:
   ```bash
   wget https://qne-installers-cdn.b-cdn.net/keeper-installer-linux-amd64
   ```

2. Make it executable:
   ```bash
   chmod +x keeper-installer-linux-amd64
   ```

3. Run it:
   ```bash
   ./keeper-installer-linux-amd64
   ```

4. The installer will:
   - Copy itself to `~/QNE/local/`
   - Create an app menu entry
   - Open the GUI installer window
   - Allow you to install the Keeper

### Windows & macOS

**Status**: Not yet available

The installer framework supports cross-platform builds, but the **Keeper application itself** currently exists only as a Linux AppImage. To support Windows and macOS, we need to:

1. Build Keeper as a Windows `.exe` (using electron-builder)
2. Build Keeper as a macOS `.app` (using electron-builder)
3. Upload platform-specific builds to CDN
4. Update installer to download the appropriate version

**Interested in helping?** See the Contributing section below!

## Building from Source

### Prerequisites

- Go 1.24 or later
- Fyne dependencies:
  - **Linux**: `sudo apt install libgl1-mesa-dev xorg-dev`
  - **Windows**: No additional dependencies
  - **macOS**: No additional dependencies

### Build

```bash
git clone https://github.com/quickneasy/keeper-installer.git
cd keeper-installer
go build
```

## Development

### Project Structure

```
keeper-installer/
├── main.go           # Installer application code
├── go.mod            # Go module definition
├── assets/           # Icons and resources
├── .github/
│   └── workflows/    # CI/CD automation
└── README.md         # This file
```

### Key Features

- **Self-installation logic** - Detects if running from install location
- **Fyne GUI** - Native cross-platform UI
- **Progress tracking** - Visual feedback during download
- **Desktop integration** - Creates `.desktop` files on Linux
- **Uninstall support** - Clean removal with confirmation

## Architecture

### Install Process

1. User downloads installer binary
2. On first run: Installer copies itself to `~/QNE/local/`
3. Installer creates app menu shortcut
4. GUI opens with Install/Uninstall options
5. User clicks "Install Keeper"
6. Downloads Keeper AppImage from CDN (125 MB)
7. Installs to `~/QNE/local/keeper`
8. Creates desktop and app menu shortcuts
9. User can launch Keeper

### File Locations

- **Installer**: `~/QNE/local/keeper-installer`
- **Keeper**: `~/QNE/local/keeper`
- **App Menu**: `~/.local/share/applications/`
- **Desktop**: `~/Desktop/`

## Contributing

Contributions welcome! This is an open-source project under the QuickNEasy organization.

### Areas for Contribution

- Windows installer improvements
- macOS installer support
- Localization/translations
- UI enhancements
- Bug fixes

## License

Copyright © 2025 QNE Community

## Links

- **Website**: https://quickneasy.info
- **GitHub Organization**: https://github.com/quickneasy
- **Keeper Main Repo**: https://github.com/quickneasy/keeper (private)

## Support

For issues or questions:
- Open an issue on GitHub
- Visit quickneasy.info for community support
