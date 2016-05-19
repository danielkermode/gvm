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

func main() {
	args := os.Args
	osArch := strings.ToLower(os.Getenv("PROCESSOR_ARCHITECTURE"))
	detail := ""

	if len(args) < 2 {
		help()
		return
	}
	if len(args) > 2 {
		detail = strings.ToLower(args[2])
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
		fmt.Println("No GOROOT set. Set a GOROOT for go installations with gvm goroot <path>.")
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
	fmt.Println("Set the GOROOT to " + newpath)
	fmt.Println("Note: You'll have to start another prompt to see the changes.")
}

func listGos() {
	if os.Getenv("GOROOT") == "" {
		fmt.Println("No GOROOT set. Set a GOROOT for go installations with gvm goroot <path>.")
		return
	}
	validDir := regexp.MustCompile(`go(\d\.\d\.\d){0,1}`)
	gorootroot := filepath.Clean(os.Getenv("GOROOT") + "\\..")
	files, _ := ioutil.ReadDir(gorootroot)
	goVers := make(map[string]bool)
	fmt.Println("")
	for _, f := range files {
		if f.IsDir() && validDir.MatchString(f.Name()) {
			goDir := filepath.Join(gorootroot, f.Name())
			version := getDirVersion(goDir)
			//check if the version already exists (different named dirs with the same go version can exist)
			_, exists := goVers[version]
			if exists {
				continue
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
	}
}

func useGo(newVer string) {
	if os.Getenv("GOROOT") == "" {
		fmt.Println("No GOROOT set. Set a GOROOT for go installations with gvm goroot <path>.")
		return
	}
	if newVer == "" {
		fmt.Println("A new version must be specified.")
	}
	validDir := regexp.MustCompile(`go(\d\.\d\.\d){0,1}`)
	gorootroot := filepath.Clean(os.Getenv("GOROOT") + "\\..")
	files, _ := ioutil.ReadDir(gorootroot)
	fmt.Println("")
	for _, f := range files {
		if f.IsDir() && validDir.MatchString(f.Name()) {
			goDir := filepath.Join(gorootroot, f.Name())
			version := getDirVersion(goDir)
			if version == "go"+newVer {
				//permanently set env var for user and local machine
				//The path should be the same for all windows OSes.
				machineEnvPath := "SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment"
				userEnvPath := "Environment"
				setEnvVar("GOROOT", goDir, machineEnvPath, true)
				setEnvVar("GOROOT", goDir, userEnvPath, false)
				fmt.Println("Now using Go version " + version[2:] + ". Set GOROOT to " + goDir)
				fmt.Println("Note: You'll have to start another prompt to see the changes.")
				return
			}
		}
	}
	fmt.Println("Couldn't use Go version " + newVer + ". Check Go versions with gvm list.")
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
	fmt.Println("  gvm goroot [path]            : This shows any GOROOT. Optional path arg to set GOROOT to that.")
	fmt.Println("  gvm list                     : List the Go installations at or adjacent to GOROOT. Aliased as ls.")
	fmt.Println("  gvm use <version>            : Switch to use the specified version. This will set your GOROOT.")
	fmt.Println("  gvm version                  : Displays the current running version of gvm for Windows. Aliased as v.")
}
