# gvm
This is a Go Version Manager for Windows, written in Go. There is a "gvm" already, but unfortunately it only supports unix systems. Basically it's the equivalent of "nvm" (Node Version Manager) but for Go. Credit to https://github.com/coreybutler/nvm-windows/. His program is also written in Go, and it was helpful looking through his code for inspiration.

I created this out of need, since I wanted the equivalent of nvm for Go on Windows, and also to familiarize myself with the Go language. The commands are currently:

```
gvm arch                     : Show architecture of OS.
gvm install <version>        : The version must be a version of Go.
gvm goroot [path]            : This shows any GOROOT. Optional path arg to set GOROOT to that.
gvm list                     : List the Go installations at or adjacent to GOROOT. Aliased as ls.
gvm use <version>            : Switch to use the specified version. This will set your GOROOT.
gvm version                  : Displays the current running version of gvm for Windows. Aliased as v.
```

I'll serve the binaries for this program at some point, but you can simply build from source by ```git clone```ing this repo into your GOPATH and running ```go install```.

**How it works, for those interested**: You must have a GOROOT environment variable set on your computer. You must also be running Windows, of course! You don't actually need Go installed if you have the .exe for this program, you can use this to get the Go files from scratch. ```gvm install <version>``` will extract the files needed for the desired version and put them in a folder called "goX.X.X" (eg. go1.6.2) which will be *adjacent* to your GOROOT. Then the other gvm commands simply look through these folders and determine the version numbers with the "VERSION" file in each directory that Go comes installed with. ```gvm use <version>``` will change your GOROOT to whatever directory satisfies the version number within this GOROOT environment. If you have an existing installation of Go, as long as the folder is called "go" or "goX.X.X" and it is adjacent to or is your GOROOT then it will be registered as normal. If there are multiple differently named directories containing the same installation of Go, the first one found will be taken for use by gvm and the others will be ignored (ideally you shouldn't have this in the first place, but it could happen by accident if you have one folder called "go" and another called "goX.X.X" with the same version).

You must restart your command prompt to register the changes to GOROOT after selecting a new Go version.

The architectures supported are the Windows ones on https://golang.org/dl/, so amd64 and 386.