package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// ImportExportCLI Version. Value is set during the build process.
var importExportCLIVersion string
var importExportCLIPOS string
var importExportCLIPArch string

//const resourceStore = "https://cdn-updates.private.wso2.com/importExportCLI-build-resources/"
const resourceStore = "https://github.com/wso2/product-apim-tooling/" // TODO:: Add a working resourceStore

var windowsData = map[string]string{
	"installer.wxs": `<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
<?if $(var.Arch) = 386 ?>
  <?define ProdId = {E2BD598E-154C-4202-9DD5-6F040877F627} ?>
  <?define UpgradeCode = {02DBD2E9-334B-4955-AAB2-50BD78BABC9D} ?>
  <?define SysFolder=SystemFolder ?>
  <?define PlatformArch=x86 ?>
  <?define ProgramFilesDir=ProgramFilesFolder ?>
<?else?>
  <?define ProdId = {813f8d1e-fd51-4d79-bc8d-a216c68aef54} ?>
  <?define UpgradeCode = {76a0492c-ada4-4b2c-be21-55ce8da4bdcc} ?>
  <?define SysFolder=System64Folder ?>
  <?define PlatformArch=x64 ?>
  <?define ProgramFilesDir=ProgramFiles64Folder ?>
<?endif?>
<Product
    Id="*"
    Name="ImportExportCLI $(var.)"
    Language="1033"
    Version="$(var.WixImportExportCLIVersion)"
    Manufacturer="https://wso2.com"
    UpgradeCode="$(var.UpgradeCode)" >
<Package
    Id='*'
    Keywords='Installer'
    Description="The ImportExportCLI Installer"
    Comments="ImportExportCLI is a CLI tool which enables importing and exporting APIs between WSO2 API Manager
	environments"
    InstallerVersion="300"
    Compressed="yes"
    InstallScope="perMachine"
    Languages="1033"
    Platform="$(var.PlatformArch)" />
<Property Id="ARPCOMMENTS" Value="ImportExportCLI is a CLI tool which enables importing and exporting APIs between WSO2
	API Manager environments" />
<Property Id="ARPCONTACT" Value="https://wso2.com/contact" />
<Property Id="ARPHELPLINK" Value="https://wso2.com/update" />
<Property Id="ARPREADME" Value="https://wso2.com/update" />
<Property Id="ARPURLINFOABOUT" Value="https://wso2.com/update" />
<Property Id="LicenseAccepted">1</Property>
<Media Id='1' Cabinet="importExportCLI.cab" EmbedCab="yes" CompressionLevel="high" />
<Condition Message="Windows XP or greater required."> VersionNT >= 500</Condition>
<MajorUpgrade AllowDowngrades="yes" />
<SetDirectory Id="INSTALLDIRROOT" Value="[$(var.ProgramFilesDir)]"/>
<CustomAction
    Id="SetApplicationRootDirectory"
    Property="ARPINSTALLLOCATION"
    Value="[INSTALLDIR]" />
<!-- Define the directory structure and environment variables -->
<Directory Id="TARGETDIR" Name="SourceDir">
  <Directory Id="INSTALLDIRROOT">
    <Directory Id="INSTALLDIR" Name="ImportExportCLI"/>
  </Directory>
  <Directory Id="ProgramMenuFolder">
    <Directory Id="ImportExportCLIProgramShortcutsDir" Name="ImportExportCLI"/>
  </Directory>
  <Directory Id="EnvironmentEntries">
    <Directory Id="ImportExportCLIEnvironmentEntries" Name="ImportExportCLI"/>
  </Directory>
</Directory>
<!-- Programs Menu Shortcuts -->
<DirectoryRef Id="ImportExportCLIProgramShortcutsDir">
  <Component Id="Component_ImportExportCLIProgramShortCuts" Guid="{87b51cae-d159-46a0-9a8b-cf206e37c35d}">
    <Shortcut
        Id="UninstallShortcut"
        Name="Uninstall ImportExportCLI"
        Description="Uninstalls ImportExportCLI"
        Target="[$(var.SysFolder)]msiexec.exe"
        Arguments="/x [ProductCode]" />
    <RemoveFolder
        Id="ImportExportCLIProgramShortcutsDir"
        On="uninstall" />
    <RegistryValue
        Root="HKCU"
        Key="Software\ImportExportCLI"
        Name="ShortCuts"
        Type="integer"
        Value="1"
        KeyPath="yes" />
  </Component>
</DirectoryRef>
<!-- Registry & Environment Settings -->
<DirectoryRef Id="ImportExportCLIEnvironmentEntries">
  <Component Id="Component_ImportExportCLIEnvironment" Guid="{b23b29a2-4b2d-4f91-82c1-0e98aa4544e5}">
    <RegistryKey
        Root="HKCU"
        Key="Software\ImportExportCLI"
        Action="create" >
            <RegistryValue
                Name="installed"
                Type="integer"
                Value="1"
                KeyPath="yes" />
            <RegistryValue
                Name="installLocation"
                Type="string"
                Value="[INSTALLDIR]" />
    </RegistryKey>
    <Environment
        Id="ImportExportCLIPathEntry"
        Action="set"
        Part="last"
        Name="PATH"
        Permanent="no"
        System="yes"
        Value="[INSTALLDIR]bin" />
    <Environment
        Id="ImportExportCLIRoot"
        Action="set"
        Part="all"
        Name="ImportExportCLIROOT"
        Permanent="no"
        System="yes"
        Value="[INSTALLDIR]" />
    <RemoveFolder
        Id="ImportExportCLIEnvironmentEntries"
        On="uninstall" />
  </Component>
</DirectoryRef>
<!-- Install the files -->
<Feature
    Id="ImportExportCLITools"
    Title="ImportExportCLI"
    Level="1">
      <ComponentRef Id="Component_ImportExportCLIEnvironment" />
      <ComponentGroupRef Id="AppFiles" />
      <ComponentRef Id="Component_ImportExportCLIProgramShortCuts" />
</Feature>
<!-- Update the environment -->
<InstallExecuteSequence>
    <Custom Action="SetApplicationRootDirectory" Before="InstallFinalize" />
</InstallExecuteSequence>
<!-- Include the user interface -->
<WixVariable Id="WixUILicenseRtf" Value="LICENSE.rtf" />
<WixVariable Id="WixUIBannerBmp" Value="Banner.jpg" />
<WixVariable Id="WixUIDialogBmp" Value="Dialog.jpg" />
<Property Id="WIXUI_INSTALLDIR" Value="INSTALLDIR" />
<UIRef Id="WixUI_InstallDir" />
</Product>
</Wix>
`,
	"LICENSE.rtf": resourceStore + "windows/LICENSE.rtf",
	"Banner.jpg":  resourceStore + "windows/importExportCLI-banner.jpg",
	"Dialog.jpg":  resourceStore + "windows/importExportCLI-dialog.jpg",
}

var darwinData = map[string]string{

	"scripts/postinstall": `#!/bin/bash
ImportExportCLIROOT=/usr/local/importExportCLI
echo "Fixing permissions"
cd $ImportExportCLIROOT
find . -exec chmod ugo+r \{\} \;
find bin -exec chmod ugo+rx \{\} \;
find . -type d -exec chmod ugo+rx \{\} \;
chmod o-w .
`,

	"scripts/preinstall": `#!/bin/bash
ImportExportCLIROOT=/usr/local/importExportCLI
echo "Removing previous installation"
if [ -d $ImportExportCLIROOT ]; then
	rm -r $ImportExportCLIROOT
fi
`,

	"Distribution": `<?xml version="1.0" encoding="utf-8" standalone="no"?>
<installer-script minSpecVersion="1.000000">
    <title>ImportExportCLI</title>
    <background mime-type="image/png" file="dialog.png"/>
    <license file="LICENSE"/>
    <welcome file="WELCOME" />
    <options customize="never" allow-external-scripts="no"/>
    <domains enable_localSystem="true" />
    <installation-check script="installCheck();"/>
    <script>
function installCheck() {
    if(!(system.compareVersions(system.version.ProductVersion, '10.6.0') >= 0)) {
        my.result.title = 'Unable to install';
        my.result.message = 'ImportExportCLI requires Mac OS X 10.6 or later.';
        my.result.type = 'Fatal';
        return false;
    }
    if(system.files.fileExistsAtPath('/usr/local/importExportCLI/bin/importExportCLI')) {
	    my.result.title = 'Previous Installation Detected';
	    my.result.message = 'A previous installation of ImportExportCLI exists at /usr/local/importExportCLI. This installer will remove the previous installation prior to installing. Please back up any data before proceeding.';
	    my.result.type = 'Warning';
	    return false;
	}
    return true;
}
    </script>
    <choices-outline>
        <line choice="com.wso2.importExportCLI.choice"/>
    </choices-outline>
    <choice id="com.wso2.importExportCLI.choice" title="ImportExportCLI">
        <pkg-ref id="com.wso2.importExportCLI.pkg"/>
    </choice>
    <pkg-ref id="com.wso2.importExportCLI.pkg" auth="Root">com.wso2.importExportCLI.pkg</pkg-ref>
</installer-script>
`,
	"Resources/dialog.png": resourceStore + "darwin/importExportCLI-dialog.png",
	"Resources/LICENSE":    resourceStore + "darwin/LICENSE",
	"Resources/WELCOME":    resourceStore + "darwin/WELCOME",
}

var versionRe = regexp.MustCompile(`^importExportCLI(\d+(\.\d+)*)`)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	switch runtime.GOOS {
	case "windows":
		err = windowsMSI()
	case "darwin":
		err = darwinPKG()
	}
	if err != nil {
		log.Fatal(err)
	}

}

func windowsMSI() error {
	cwd, version, err := environmentInfo()
	if err != nil {
		return err
	}

	//Install Wix tools.
	wix := filepath.Join(cwd, "wix")
	defer os.RemoveAll(wix)
	if err := installWix(wix); err != nil {
		return err
	}

	// Write out windows data that is used by the packaging process.
	win := filepath.Join(cwd, "windows")

	defer os.RemoveAll(win)
	if err := writeDataFiles(windowsData, win); err != nil {
		return err
	}

	// Gather files.
	importExportCLIDir := filepath.Join(cwd, "importExportCLI")
	appfiles := filepath.Join(win, "AppFiles.wxs")
	if err := runDir(win, filepath.Join(wix, "heat"),
		"dir", importExportCLIDir,
		"-nologo",
		"-gg", "-g1", "-srd", "-sfrag",
		"-cg", "AppFiles",
		"-template", "fragment",
		"-dr", "INSTALLDIR",
		"-var", "var.SourceDir",
		"-out", appfiles,
	); err != nil {
		return err
	}

	// Build package.
	if err := runDir(win, filepath.Join(wix, "candle"),
		"-nologo",
		"-dImportExportCLIVersion="+version,
		"-dWixImportExportCLIVersion="+wixVersion(version),
		"-dArch="+runtime.GOARCH,
		"-dSourceDir="+importExportCLIDir,
		filepath.Join(win, "installer.wxs"),
		appfiles,
	); err != nil {
		return err
	}

	msi := filepath.Join(cwd, "msi")
	if err := os.Mkdir(msi, 0755); err != nil {
		return err
	}
	return runDir(win, filepath.Join(wix, "light"),
		"-nologo",
		"-dcl:high",
		"-ext", "WixUIExtension",
		"-ext", "WixUtilExtension",
		"AppFiles.wixobj",
		"installer.wixobj",
		"-o", filepath.Join(msi, "importExportCLI-"+version+"-"+importExportCLIPOS+"-"+importExportCLIPArch+".msi"),
	)
}

func environmentInfo() (cwd, version string, err error) {
	cwd, err = os.Getwd()
	if err != nil {
		return
	}
	version = strings.TrimSpace(importExportCLIVersion)
	return
}

func installWix(path string) error {
	// Fetch wix binary zip file.
	body, err := httpGet(resourceStore + "windows/wix310-binaries.zip")
	if err != nil {
		return err
	}

	// Unzip to path.
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	if len(zr.File) <= 0 {
		fmt.Println("No zip")
	}
	for _, f := range zr.File {
		name := filepath.FromSlash(f.Name)
		err := os.MkdirAll(filepath.Join(path, filepath.Dir(name)), 0755)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(path, name), b, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func httpGet(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.New(r.Status)
	}
	return body, nil
}

func writeDataFiles(data map[string]string, base string) error {
	for name, body := range data {
		dst := filepath.Join(base, name)
		err := os.MkdirAll(filepath.Dir(dst), 0755)
		if err != nil {
			return err
		}
		b := []byte(body)
		if strings.HasPrefix(body, resourceStore) {
			b, err = httpGet(body)
			if err != nil {
				return err
			}
		}
		// (We really mean 0755 on the next line; some of these files
		// are executable, and there's no harm in making them all so.)
		if err := ioutil.WriteFile(dst, b, 0755); err != nil {
			return err
		}
	}
	return nil
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func runDir(dir, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func wixVersion(v string) string {
	m := versionRe.FindStringSubmatch(v)
	if m == nil {
		return "0.0.0"
	}
	return m[1]
}

func darwinPKG() error {
	cwd, version, err := environmentInfo()
	if err != nil {
		return err
	}

	// Write out darwin data that is used by the packaging process.
	defer os.RemoveAll("darwin")
	if err := writeDataFiles(darwinData, "darwin"); err != nil {
		return err
	}

	// Create a work directory and place inside the files as they should
	// be on the destination file system.
	work := filepath.Join(cwd, "darwinpkg")
	if err := os.MkdirAll(work, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(work)

	// Write out /etc/paths.d/importExportCLI.
	const pathsBody = "/usr/local/importExportCLI/bin"
	pathsDir := filepath.Join(work, "etc/paths.d")
	pathsFile := filepath.Join(pathsDir, "importExportCLI")
	if err := os.MkdirAll(pathsDir, 0755); err != nil {
		return err
	}
	if err = ioutil.WriteFile(pathsFile, []byte(pathsBody), 0644); err != nil {
		return err
	}

	// Copy ImportExportCLI installation to /usr/local/importExportCLI.
	importExportCLIDir := filepath.Join(work, "usr/local/importExportCLI")
	if err := os.MkdirAll(importExportCLIDir, 0755); err != nil {
		return err
	}
	if err := cpDir(importExportCLIDir, "importExportCLI"); err != nil {
		return err
	}

	// Build the package file.
	dest := "package"
	if err := os.Mkdir(dest, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(dest)

	if err := run("pkgbuild",
		"--identifier", "com.wso2.importExportCLI",
		"--version", version,
		"--scripts", "darwin/scripts",
		"--root", work,
		filepath.Join(dest, "com.wso2.importExportCLI.pkg"),
	); err != nil {
		return err
	}

	const pkg = "pkg" // known to cmd/release
	if err := os.Mkdir(pkg, 0755); err != nil {
		return err
	}

	return run("productbuild",
		"--distribution", "darwin/Distribution",
		"--resources", "darwin/Resources",
		"--package-path", dest,
		filepath.Join(cwd, pkg, "importExportCLI-"+version+"-"+importExportCLIPOS+"-"+importExportCLIPArch+".pkg"), // file name irrelevant
	)
}

func cpDir(dst, src string) error {
	walk := func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, srcPath[len(src):])
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}
		return cp(dstPath, srcPath)
	}
	return filepath.Walk(src, walk)
}

func cp(dst, src string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	fi, err := sf.Stat()
	if err != nil {
		return err
	}
	tmpDst := dst + ".tmp"
	df, err := os.Create(tmpDst)
	if err != nil {
		return err
	}
	defer df.Close()

	if runtime.GOOS != "windows" {
		if err := df.Chmod(fi.Mode()); err != nil {
			return err
		}
	}
	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}
	if err := df.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpDst, dst); err != nil {
		return err
	}
	// Ensure the destination has the same mtime as the source.
	return os.Chtimes(dst, fi.ModTime(), fi.ModTime())
}
