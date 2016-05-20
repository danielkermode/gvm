package web

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var goBaseAddress = "https://storage.googleapis.com/golang/go"

func Download(version string, arch string, root string) bool {
	url := goBaseAddress + version + "." + arch + ".zip"
	filedir := root + "\\" + "go" + version + ".zip"
	fmt.Println("Downloading Go v" + version + "... Please wait...")
	// create a file to store downloaded data
	out, err := os.Create(filedir)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
	}

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
	}

	if response.Status[0:3] != "200" {
		fmt.Println("Download failed for url " + url + ". You can check the url manually. Rolling Back.")
		//remove the zip file after closing it.
		defer os.Remove(filedir)
		defer out.Close()
		return false
	}
	//this is the folder with goX.X.X. Initially it has a folder called "go" with the required files.
	//The aim is to move the files from "go" into the parent goX.X.X.
	dest := filepath.Join(root, "go"+version)
	//this is the "go" folder inside the new folder containing the files for running Go.
	godest := filepath.Join(dest, "go")
	fmt.Println("Unzipping files...")
	ziperr := unzip(filedir, dest, version)
	if ziperr != nil {
		fmt.Println("Error while unzipping", url, "-", ziperr)
		return false
	}
	//remove the zip file after closing it.
	defer os.Remove(filedir)
	defer out.Close()
	//delete contents of go folder (happens after below defer)
	defer os.RemoveAll(godest)
	//copy contents of go folder (before deleting it)
	defer copyDir(godest, dest)
	return true
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()
	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return
}

func copyDir(source string, dest string) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		fmt.Println("Error while adjusting new directory", err)
	}
	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		fmt.Println("Error while adjusting new directory", err)
	}
	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			// create sub-directories - recursively
			copyDir(sourcefilepointer, destinationfilepointer)
		} else {
			// perform copy
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func unzip(src string, dest string, version string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	os.MkdirAll(dest, 0755)
	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			//make any needed folders for the file in question
			os.MkdirAll(filepath.Clean(path+"\\.."), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}
