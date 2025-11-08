package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	browserURL     = "https://qne-installers-cdn-cdn.b-cdn.net/keeper-1.0.0.AppImage"
	browserSHA256  = "" // TODO: Add checksum
	appName        = "keeper"
	appVersion     = "1.0.0"
)

func main() {
	// Check if installer itself needs to be installed
	if !isInstallerInstalled() {
		// First run - perform self-installation
		if err := performSelfInstall(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to install: %v\n", err)
			os.Exit(1)
		}
	}

	myApp := app.NewWithID("com.qnecommunity.browser.installer")
	myWindow := myApp.NewWindow("Keeper Installer")
	myWindow.Resize(fyne.NewSize(600, 500))
	myWindow.CenterOnScreen()

	// Create main UI
	ui := newInstallerUI(myApp, myWindow)
	myWindow.SetContent(ui.makeWelcomeScreen())
	
	myWindow.ShowAndRun()
}

// getInstallDir returns the installation directory for the installer itself
func getInstallDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "QNE", "local")
}

// getInstallerPath returns the full path where installer should be installed
func getInstallerPath() string {
	installerName := "keeper-installer"
	if runtime.GOOS == "windows" {
		installerName += ".exe"
	}
	return filepath.Join(getInstallDir(), installerName)
}

// isInstallerInstalled checks if the installer is already installed
func isInstallerInstalled() bool {
	installedPath := getInstallerPath()
	currentExe, err := os.Executable()
	if err != nil {
		return false
	}
	
	// If we're running from the install location, we're already installed
	currentExe, _ = filepath.EvalSymlinks(currentExe)
	installedPath, _ = filepath.EvalSymlinks(installedPath)
	
	return currentExe == installedPath
}

// copyFileSimple copies a file from src to dst
func copyFileSimple(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// performSelfInstall installs the installer to the user's system
func performSelfInstall() error {
	installDir := getInstallDir()
	installerPath := getInstallerPath()
	
	fmt.Println("Installing Keeper Installer...")
	fmt.Printf("Install location: %s\n", installDir)
	
	// Create install directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %v", err)
	}
	
	// Copy executable to install location
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %v", err)
	}
	
	if err := copyFileSimple(currentExe, installerPath); err != nil {
		return fmt.Errorf("failed to copy installer: %v", err)
	}
	
	// Make it executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(installerPath, 0755); err != nil {
			return fmt.Errorf("failed to set executable permission: %v", err)
		}
	}
	
	// Create desktop entry (Linux only for now)
	if runtime.GOOS == "linux" {
		if err := createInstallerDesktopEntry(installerPath); err != nil {
			fmt.Printf("Warning: Could not create desktop entry: %v\n", err)
		}
	}
	
	fmt.Println("✓ Installation complete!")
	fmt.Println("You can now launch 'Keeper Installer' from your application menu")
	fmt.Println()
	
	return nil
}

// createInstallerDesktopEntry creates a desktop entry for the installer itself
func createInstallerDesktopEntry(installerPath string) error {
	home, _ := os.UserHomeDir()
	desktopDir := filepath.Join(home, ".local", "share", "applications")
	
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		return err
	}
	
	desktopFile := filepath.Join(desktopDir, "keeper-installer.desktop")
	content := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=Keeper Installer
Comment=Install Chameleon Keeper
Exec=%s
Icon=system-software-install
Terminal=false
Categories=Network;WebBrowser;Utility;
Keywords=browser;web;internet;installer;qne;
StartupNotify=true
`, installerPath)
	
	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		return err
	}
	
	// Update desktop database
	exec.Command("update-desktop-database", desktopDir).Run()
	
	return nil
}

type installerUI struct {
	app    fyne.App
	window fyne.Window
}

func newInstallerUI(app fyne.App, window fyne.Window) *installerUI {
	return &installerUI{
		app:    app,
		window: window,
	}
}

func (ui *installerUI) makeWelcomeScreen() fyne.CanvasObject {
	// Header with logo/title
	title := canvas.NewText("QuickNEasy Browser", theme.ForegroundColor())
	title.TextSize = 32
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("Secure P2P Web Browser")
	subtitle.Alignment = fyne.TextAlignCenter

	version := widget.NewLabel("Version " + appVersion)
	version.Alignment = fyne.TextAlignCenter
	version.TextStyle = fyne.TextStyle{Italic: true}

	// Welcome message
	welcome := widget.NewLabel(
		"Welcome to the Keeper installer!\n\n" +
			"This will install the Chameleon Keeper to your system.\n" +
			"The Keeper manages your chameleons and provides skills like browsing, email, and messaging.",
	)
	welcome.Wrapping = fyne.TextWrapWord

	// Features list
	features := widget.NewCard("Features", "", widget.NewLabel(
		"✓ Peer-to-peer architecture\n"+
			"✓ End-to-end encryption\n"+
			"✓ No central servers\n"+
			"✓ Privacy-focused\n"+
			"✓ Open source",
	))

	// Install info
	installPath := filepath.Join(os.Getenv("HOME"), "QNE", "local", appName)
	installInfo := widget.NewLabel(fmt.Sprintf("Install location: %s", installPath))
	installInfo.TextStyle = fyne.TextStyle{Italic: true}

	// Buttons
	installBtn := widget.NewButton("Install Keeper", func() {
		ui.window.SetContent(ui.makeInstallScreen())
	})
	installBtn.Importance = widget.HighImportance

	uninstallBtn := widget.NewButton("Uninstall", func() {
		ui.confirmUninstall()
	})
	
	quitBtn := widget.NewButton("Cancel", func() {
		ui.app.Quit()
	})

	buttonBox := container.NewHBox(
		layout.NewSpacer(),
		quitBtn,
		uninstallBtn,
		installBtn,
	)

	// Layout
	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(title),
		container.NewCenter(subtitle),
		container.NewCenter(version),
		layout.NewSpacer(),
		container.NewPadded(welcome),
		container.NewPadded(features),
		layout.NewSpacer(),
		container.NewCenter(installInfo),
		layout.NewSpacer(),
		buttonBox,
		layout.NewSpacer(),
	)

	return container.NewPadded(content)
}

func (ui *installerUI) makeInstallScreen() fyne.CanvasObject {
	title := canvas.NewText("Installing Keeper", theme.ForegroundColor())
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	statusLabel := widget.NewLabel("Preparing installation...")
	statusLabel.Alignment = fyne.TextAlignCenter

	progressBar := widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100

	cancelBtn := widget.NewButton("Cancel", func() {
		// TODO: Implement cancellation
		ui.app.Quit()
	})
	cancelBtn.Disable() // Disable during download

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(title),
		layout.NewSpacer(),
		container.NewCenter(statusLabel),
		container.NewPadded(progressBar),
		layout.NewSpacer(),
		container.NewCenter(cancelBtn),
		layout.NewSpacer(),
	)

	// Start installation in background
	go ui.performInstallation(statusLabel, progressBar, cancelBtn)

	return container.NewPadded(content)
}

func (ui *installerUI) performInstallation(statusLabel *widget.Label, progressBar *widget.ProgressBar, cancelBtn *widget.Button) {
	// Determine install directory - use ~/QNE/local
	installDir := filepath.Join(os.Getenv("HOME"), "QNE", "local")
	installPath := filepath.Join(installDir, appName)

	// Create install directory
	statusLabel.SetText("Creating installation directory...")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		ui.showError("Failed to create installation directory", err)
		return
	}
	progressBar.SetValue(10)

	// Download browser
	statusLabel.SetText("Downloading Keeper (126 MB)...")
	tmpFile, err := os.CreateTemp("", "keeper-*.AppImage")
	if err != nil {
		ui.showError("Failed to create temporary file", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if err := ui.downloadWithProgress(browserURL, tmpFile, progressBar, 10, 80); err != nil {
		ui.showError("Download failed", err)
		return
	}
	tmpFile.Close()

	// Verify checksum if provided
	if browserSHA256 != "" {
		statusLabel.SetText("Verifying download...")
		if err := ui.verifyChecksum(tmpFile.Name(), browserSHA256); err != nil {
			ui.showError("Checksum verification failed", err)
			return
		}
	}
	progressBar.SetValue(85)

	// Copy to install location
	statusLabel.SetText("Installing browser...")
	if err := ui.copyFile(tmpFile.Name(), installPath, 0755); err != nil {
		ui.showError("Failed to install browser", err)
		return
	}
	progressBar.SetValue(90)

	// Create desktop entry
	statusLabel.SetText("Creating shortcuts...")
	if err := ui.createDesktopEntry(); err != nil {
		// Non-fatal, just log
		fmt.Printf("Warning: Failed to create desktop entry: %v\n", err)
	}
	progressBar.SetValue(100)

	// Show success screen
	ui.window.SetContent(ui.makeSuccessScreen(installPath))
}

func (ui *installerUI) downloadWithProgress(url string, dest *os.File, progressBar *widget.ProgressBar, startPct, endPct float64) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	totalSize := resp.ContentLength
	downloaded := int64(0)
	buffer := make([]byte, 32*1024)

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := dest.Write(buffer[:n]); writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)
			
			// Update progress
			if totalSize > 0 {
				pct := float64(downloaded) / float64(totalSize)
				actualPct := startPct + (pct * (endPct - startPct))
				progressBar.SetValue(actualPct)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (ui *installerUI) verifyChecksum(filePath, expectedHash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	actualHash := hex.EncodeToString(h.Sum(nil))
	if actualHash != expectedHash {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

func (ui *installerUI) copyFile(src, dst string, perm os.FileMode) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return os.Chmod(dst, perm)
}

func (ui *installerUI) createDesktopEntry() error {
	installPath := filepath.Join(os.Getenv("HOME"), "QNE", "local", appName)

	// Simple desktop entry format - direct execution works best
	content := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=Keeper
Comment=Chameleon Keeper - P2P Platform
Exec=%s --no-sandbox
Icon=web-browser
Terminal=false
Categories=Network;WebBrowser;
Keywords=qne;browser;p2p;
StartupNotify=true
`, installPath)

	// Create in applications directory (for app menu)
	desktopDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		return err
	}
	
	desktopFile := filepath.Join(desktopDir, "keeper.desktop")
	if err := os.WriteFile(desktopFile, []byte(content), 0644); err != nil {
		return err
	}

	// Also create on Desktop (for desktop icon)
	desktopPath := filepath.Join(os.Getenv("HOME"), "Desktop", "keeper.desktop")
	if err := os.WriteFile(desktopPath, []byte(content), 0755); err != nil {
		// Non-fatal if Desktop doesn't exist or isn't writable
		fmt.Printf("Note: Could not create desktop shortcut: %v\n", err)
	} else {
		// Mark desktop file as trusted so it's executable
		exec.Command("gio", "set", desktopPath, "metadata::trusted", "true").Run()
	}

	return nil
}

func (ui *installerUI) makeSuccessScreen(installPath string) fyne.CanvasObject {
	title := canvas.NewText("Installation Complete!", theme.SuccessColor())

successMsg := widget.NewLabel(
"Keeper has been successfully installed!\n\n" +
"You can now launch it from your application menu\n" +
"or run it from the terminal.",
)
successMsg.Alignment = fyne.TextAlignCenter
successMsg.Wrapping = fyne.TextWrapWord

pathInfo := widget.NewLabel(fmt.Sprintf("Installed to: %s", installPath))
pathInfo.Alignment = fyne.TextAlignCenter
pathInfo.TextStyle = fyne.TextStyle{Italic: true}

// Buttons
launchBtn := widget.NewButton("Launch Browser", func() {
ui.launchBrowser(installPath)
})
launchBtn.Importance = widget.HighImportance

closeBtn := widget.NewButton("Close", func() {
ui.app.Quit()
})

buttonBox := container.NewHBox(
layout.NewSpacer(),
closeBtn,
launchBtn,
)

content := container.NewVBox(
layout.NewSpacer(),
container.NewCenter(title),
layout.NewSpacer(),
container.NewPadded(successMsg),
container.NewCenter(pathInfo),
layout.NewSpacer(),
buttonBox,
layout.NewSpacer(),
)

return container.NewPadded(content)
}

func (ui *installerUI) launchBrowser(browserPath string) {
// Launch browser in background with --no-sandbox flag (required for AppImage)
cmd := exec.Command(browserPath, "--no-sandbox")
if err := cmd.Start(); err != nil {
ui.showError("Failed to launch browser", err)
return
}
	
// Close installer after a brief delay
go func() {
time.Sleep(500 * time.Millisecond)
ui.app.Quit()
}()
}

func (ui *installerUI) showError(message string, err error) {
dialog.ShowError(fmt.Errorf("%s: %v", message, err), ui.window)
}

func (ui *installerUI) confirmUninstall() {
installPath := filepath.Join(os.Getenv("HOME"), "QNE", "local", appName)
	
// Check if browser is installed
if _, err := os.Stat(installPath); os.IsNotExist(err) {
dialog.ShowInformation("Not Installed", 
"Keeper is not currently installed.", ui.window)
return
}
	
dialog.ShowConfirm("Uninstall Keeper",
"Are you sure you want to uninstall Keeper?\n\n"+
"This will remove:\n"+
"• Browser binary\n"+
"• Desktop shortcut\n"+
"• Application menu entry",
func(confirmed bool) {
if confirmed {
ui.performUninstall()
}
}, ui.window)
}

func (ui *installerUI) performUninstall() {
installPath := filepath.Join(os.Getenv("HOME"), "QNE", "local", appName)
installDir := filepath.Join(os.Getenv("HOME"), "QNE", "local")
qneDir := filepath.Join(os.Getenv("HOME"), "QNE")
desktopFile := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications", "keeper.desktop")
desktopIcon := filepath.Join(os.Getenv("HOME"), "Desktop", "keeper.desktop")
	
var errors []string
	
// Remove browser binary
if err := os.Remove(installPath); err != nil && !os.IsNotExist(err) {
errors = append(errors, fmt.Sprintf("Failed to remove browser: %v", err))
}
	
// Remove desktop file
if err := os.Remove(desktopFile); err != nil && !os.IsNotExist(err) {
errors = append(errors, fmt.Sprintf("Failed to remove app menu entry: %v", err))
}
	
// Remove desktop icon
if err := os.Remove(desktopIcon); err != nil && !os.IsNotExist(err) {
errors = append(errors, fmt.Sprintf("Failed to remove desktop icon: %v", err))
}
	
// Try to remove directories if empty
os.Remove(installDir)  // Will only succeed if empty
os.Remove(qneDir)      // Will only succeed if empty
	
// Update desktop database
exec.Command("update-desktop-database", 
filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")).Run()
	
if len(errors) > 0 {
dialog.ShowError(fmt.Errorf("Uninstall completed with errors:\n%s", 
strings.Join(errors, "\n")), ui.window)
} else {
ui.window.SetContent(ui.makeUninstallSuccessScreen())
}
}

func (ui *installerUI) makeUninstallSuccessScreen() fyne.CanvasObject {
title := canvas.NewText("Uninstall Complete!", theme.SuccessColor())
title.TextSize = 28
title.TextStyle = fyne.TextStyle{Bold: true}
title.Alignment = fyne.TextAlignCenter

successMsg := widget.NewLabel(
"Keeper has been successfully uninstalled.\n\n" +
"All files and shortcuts have been removed.",
)
successMsg.Alignment = fyne.TextAlignCenter
successMsg.Wrapping = fyne.TextWrapWord

closeBtn := widget.NewButton("Close", func() {
		ui.app.Quit()
	})
	closeBtn.Importance = widget.HighImportance

	buttonBox := container.NewHBox(
		layout.NewSpacer(),
		closeBtn,
		layout.NewSpacer(),
	)

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(title),
		layout.NewSpacer(),
		container.NewPadded(successMsg),
		layout.NewSpacer(),
		buttonBox,
		layout.NewSpacer(),
	)

	return container.NewPadded(content)
}
