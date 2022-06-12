⛔️ DEPRECATED: This project is deprecated since I'm too lazy to support it, sorry. Feel free to fork it.

# gvm
This is a Go Version Manager for Windows, written in Go. There is a "gvm" already, but unfortunately it only supports unix systems. Basically it's the equivalent of "nvm" (Node Version Manager) but for Go. Credit to https://github.com/coreybutler/nvm-windows/. His program is also written in Go, and it was helpful looking through his code for inspiration.

I created this out of need, since I wanted the equivalent of nvm for Go on Windows, and also to familiarize myself with the Go language. The commands are currently:

```
gvm arch                     : Show architecture of OS.
gvm install <version>        : Install a valid version of Go.
gvm goroot [path]            : Set/append GOROOT/PATH. Without the extra arg just shows current GOROOT. TAKE CARE, THIS CHANGES YOUR PATH!
gvm list                     : List the Go installations at or adjacent to GOROOT. Aliased as ls.
gvm use <version>            : Switch to use the specified version. This will set your GOROOT and PATH.
gvm uninstall <version>      : Uninstall specified version of Go. If it was your GOROOT/PATH, make sure to set a new one after.
gvm version                  : Display the current running version of gvm for Windows. Aliased as v.
```

If you want to build from source,  ```git clone``` this repo into your GOPATH and run ```go install```.

If you have Go installed, the simplest way to install is probably ```go get github.com/danielkermode/gvm```.Then if GOPATH/bin is in your PATH, you can run ```gvm``` from anywhere in the command line and you're good to go.

If you just want the executable (which you can use to get Go installations if you don't have one on your machine):

amd64: https://github.com/danielkermode/gvm/releases/download/v1.0.0-amd64/gvm.exe

386: https://github.com/danielkermode/gvm/releases/download/v1.0.0-386/gvm.exe

The above are found at https://github.com/danielkermode/gvm/releases/.

**How it works, for those interested**: You should have a GOROOT environment variable set on your computer, but if you don't this program can set one for you with ```gvm goroot```. You must also be running Windows, of course! You don't actually need Go installed if you have the .exe for this program, you can use this to get the Go files from scratch. ```gvm install <version>``` will extract the files needed for the desired version and put them in a folder called "goX.X.X" (eg. go1.6.2) which will be *adjacent* to your GOROOT. Then the other gvm commands simply look through these folders and determine the version numbers with the "VERSION" file in each directory that Go comes installed with. ```gvm use <version>``` will change your GOROOT to whatever directory satisfies the version number within this GOROOT environment. It will also change your PATH variable, by appending to PATH as needed. If you have an existing installation of Go, as long as the folder is called "go" or "goX.X.X" and it is adjacent to or is your GOROOT then it will be registered as normal. If there are multiple differently named directories containing the same installation of Go, the first one found will be taken for use by gvm and the others will be ignored (ideally you shouldn't have this in the first place, but it could happen by accident if you have one folder called "go" and another called "goX.X.X" with the same version).

A note on the PATH changing: it's fairly safe, I'm appending "\\bin" to all adjustments and it's extremely unlikely to screw up your PATH unless you enter something amazingly stupid that I haven't foreseen. However just make sure that your desired GOROOT is actually where Go is installed or where you want it to be installed if you do set it with ```gvm goroot```. If you're paranoid (and I wouldn't blame you; I nearly fucked up my whole machine while testing the PATH adjustments) then just set your GOROOT and PATH yourself, and you can use ```install``` and ```use``` as normal.

If you read the source, be forgiving with the code. I'm not experienced with Go so I'm sure there are things I've overlooked or not done in the best way.

You must restart your command prompt to register the changes to GOROOT after selecting a new Go version.

The architectures supported are the Windows ones on https://golang.org/dl/, so amd64 and 386.
