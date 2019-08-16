package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"gvm/web"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	GvmVersion = "1.0.0"
)

//callback function for looping over files. If true, breaks the loop.
type callback func(file os.FileInfo, reg *regexp.Regexp, proot string) bool

func main() {
	args := os.Args
	osArch := strings.ToLower(os.Getenv("PROCESSOR_ARCHITECTURE"))
	detail := ""

	if osArch == "x86" {
		osArch = "386"
	}

	if len(args) < 2 {
		help()
		return
	}
	if len(args) > 2 {
		detail = args[2]
	}
	if len(args) > 3 {
		fmt.Println("Too many args: gvm expects 2 maximum.")
	}

	// Run the appropriate method
	switch args[1] {
	case "arch":
		fmt.Println("System Architecture: " + osArch)
	case "install":
		success := install(detail, osArch)
		if success {
			fmt.Println("Successfully installed Go version " + detail + ".")
			fmt.Println("To use this version, run gvm use " + detail + ". This will also set your GOROOT.")
		}
	case "goroot":
		goroot(detail)
	case "list":
		listGos()
	case "ls":
		listGos()
	case "uninstall":
		uninstall(detail)
	case "use":
		useGo(detail)
	case "version":
		fmt.Println(GvmVersion)
	case "v":
		fmt.Println(GvmVersion)
	default:
		help()
	}
}

func install(version string, arch string) bool {
	fmt.Println("")
	if os.Getenv("GOROOT") == "" {
		fmt.Println("No GOROOT set. Set a GOROOT for Go installations with gvm goroot <path>.")
		return false
	}
	if version == "" {
		fmt.Println("Version not specified.")
		return false
	}
	gorootroot := filepath.Clean(os.Getenv("GOROOT") + "\\..")
	return web.Download(version, "windows-"+arch, gorootroot)
}

func goroot(path string) {
	fmt.Println("")
	if path == "" {
		if os.Getenv("GOROOT") == "" {
			fmt.Println("No GOROOT set.")
		} else {
			fmt.Println("GOROOT: ", os.Getenv("GOROOT"))
			fmt.Println("Other Go versions installed at: ", filepath.Clean(os.Getenv("GOROOT")+"\\.."))
		}
		return
	}
	newpath := filepath.FromSlash(path)
	//permanently set env var for user and local machine
	//The path should be the same for all windows OSes.
	machineEnvPath := "SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment"
	userEnvPath := "Environment"
	setEnvVar("GOROOT", newpath, machineEnvPath, true)
	setEnvVar("GOROOT", newpath, userEnvPath, false)
	//Also update path for user and local machine
	updatePathVar("PATH", filepath.FromSlash(os.Getenv("GOROOT")), newpath, machineEnvPath, true)
	updatePathVar("PATH", filepath.FromSlash(os.Getenv("GOROOT")), newpath, userEnvPath, false)
	fmt.Println("Set the GOROOT to " + newpath + ". Also updated PATH.")
	fmt.Println("Note: You'll have to start another prompt to see the changes.")
}

func listGos() {
	if os.Getenv("GOROOT") == "" {
		fmt.Println("No GOROOT set. Set a GOROOT for go installations with gvm goroot <path>.")
		return
	}
	//store all Go versions so we don't list duplicates
	goVers := make(map[string]bool)

	callb := func(f os.FileInfo, validDir *regexp.Regexp, gorootroot string) bool {
		if f.IsDir() && validDir.MatchString(f.Name()) {
			goDir := filepath.Join(gorootroot, f.Name())
			version := getDirVersion(goDir)
			//check if the version already exists (different named dirs with the same go version can exist)
			_, exists := goVers[version]
			if exists {
				return false
			}
			str := ""
			if goDir == os.Getenv("GOROOT") {
				str = str + "  * " + version[2:] + " (Using with GOROOT " + goDir + ")"
			} else {
				str = str + "    " + version[2:]
			}
			goVers[version] = true
			fmt.Println(str)
		}
		return false
	}

	loopFiles(callb)
}

func uninstall(unVer string) {
	if unVer == "" {
		fmt.Println("A version to uninstall must be specified.")
		return
	}

	callb := func(f os.FileInfo, validDir *regexp.Regexp, gorootroot string) bool {
		if f.IsDir() && validDir.MatchString(f.Name()) {
			goDir := filepath.Join(gorootroot, f.Name())
			version := getDirVersion(goDir)
			if version == "go"+unVer {
				os.RemoveAll(goDir)
				fmt.Println("Uninstalled Go version " + version[2:] + ".")
				fmt.Println("Note: If this was your GOROOT, make sure to set a new GOROOT with gvm goroot <path>")
				return true
			}
		}
		return false
	}

	found := loopFiles(callb)
	if !found {
		fmt.Println("Couldn't uninstall Go version " + unVer + ". Check Go versions with gvm list.")
	}
}

func useGo(newVer string) {
	if os.Getenv("GOROOT") == "" {
		fmt.Println("No GOROOT set. Set a GOROOT for go installations with gvm goroot <path>.")
		return
	}
	if newVer == "" {
		fmt.Println("A new version must be specified.")
		return
	}
	callb := func(f os.FileInfo, validDir *regexp.Regexp, gorootroot string) bool {
		if f.IsDir() && validDir.MatchString(f.Name()) {
			goDir := filepath.Join(gorootroot, f.Name())
			version := getDirVersion(goDir)
			if version == "go"+newVer {
				//permanently set env var for user and local machine
				//The path should be the same for all windows OSes.
				machineEnvPath := "SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment"
				userEnvPath := "Environment"
				setEnvVar("GOROOT", filepath.FromSlash(goDir), machineEnvPath, true)
				setEnvVar("GOROOT", filepath.FromSlash(goDir), userEnvPath, false)
				//Also update path for user and local machine
				updatePathVar("PATH", filepath.FromSlash(os.Getenv("GOROOT")), goDir, machineEnvPath, true)
				updatePathVar("PATH", filepath.FromSlash(os.Getenv("GOROOT")), goDir, userEnvPath, false)
				fmt.Println("Now using Go version " + version[2:] + ". Set GOROOT to " + goDir + ". Also updated PATH.")
				fmt.Println("Note: You'll have to start another prompt to see the changes.")
				return true
			}
		}
		return false
	}
	found := loopFiles(callb)
	if !found {
		fmt.Println("Couldn't use Go version " + newVer + ". Check Go versions with gvm list.")
	}
}

func loopFiles(fn callback) bool {
	validDir := regexp.MustCompile(`go(\d\.\d\.\d){0,1}`)
	gorootroot := filepath.Clean(os.Getenv("GOROOT") + "\\..")
	files, _ := ioutil.ReadDir(gorootroot)
	fmt.Println("")
	for _, f := range files {
		shouldBreak := fn(f, validDir, gorootroot)
		if shouldBreak {
			return true
		}
	}
	return false
}

func setEnvVar(envVar string, newVal string, envPath string, machine bool) {
	//this sets the environment variable (GOROOT in this case) for either LOCAL_MACHINE or CURRENT_USER.
	//They are set in the registry. both must be set since the GOROOT could be used from either location.
	regplace := registry.CURRENT_USER
	if machine {
		regplace = registry.LOCAL_MACHINE
	}
	key, err := registry.OpenKey(regplace, envPath, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer key.Close()

	err = key.SetStringValue(envVar, newVal)
	if err != nil {
		fmt.Println("error", err)
	}
}

func updatePathVar(envVar string, oldVal string, newVal string, envPath string, machine bool) {
	//this sets the environment variable for either LOCAL_MACHINE or CURRENT_USER.
	//They are set in the registry. both must be set since the GOROOT could be used from either location.
	regplace := registry.CURRENT_USER
	if machine {
		regplace = registry.LOCAL_MACHINE
	}
	key, err := registry.OpenKey(regplace, envPath, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	defer key.Close()

	val, _, kerr := key.GetStringValue(envVar)
	if kerr != nil {
		fmt.Println("error", err)
		return
	}
	pvars := strings.Split(val, ";")
	for i, pvar := range pvars {
		if pvar == newVal+"\\bin" {
			//the requested new value already exists in PATH, do nothing
			return
		}
		if pvar == oldVal+"\\bin" {
			pvars = append(pvars[:i], pvars[i+1:]...)
		}
	}
	val = strings.Join(pvars, ";")
	val = val + ";" + newVal + "\\bin"
	err = key.SetStringValue("PATH", val)
	if err != nil {
		fmt.Println("error", err)
	}
}

func getDirVersion(dir string) string {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.Name() == "VERSION" {
			dat, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
			if err != nil {
				return "Error reading file."
			}
			return string(dat)
		}
	}
	return ""
}

func help() {
	fmt.Println("\nRunning version " + GvmVersion + ".")
	fmt.Println("\nUsage:")
	fmt.Println(" ")
	fmt.Println("  gvm arch                     : Show architecture of OS.")
	fmt.Println("  gvm install <version>        : The version must be a version of Go.")
	fmt.Println("  gvm goroot [path]            : Sets/appends GOROOT/PATH. Without the extra arg just shows current GOROOT.")
	fmt.Println("  gvm list                     : List the Go installations at or adjacent to GOROOT. Aliased as ls.")
	fmt.Println("  gvm uninstall <version>      : Uninstall specified version of Go. If it was your GOROOT/PATH, make sure to set a new one after.")
	fmt.Println("  gvm use <version>            : Switch to use the specified version. This will set your GOROOT and PATH.")
	fmt.Println("  gvm version                  : Displays the current running version of gvm for Windows. Aliased as v.")
}
