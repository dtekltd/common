package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func FileExists(file string) bool {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func EnsureFilePath(path string, file string) string {
	return filepath.Join(EnsureDirPath(path), file)
}

func EnsureDirPath(path string) string {
	var isAbsolute bool
	isWin := strings.Contains(runtime.GOOS, "window")
	if isWin && strings.Contains(path, ":") {
		isAbsolute = true
	}
	if !isWin && path[0] == '/' {
		isAbsolute = true
	}
	if !isAbsolute {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(dir, path)
	}
	if !FileExists(path) {
		err := os.Mkdir(path, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
	return path
}

func CopyFile(source string, dest string) (err error) {
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
		if sourceinfo, err := os.Stat(source); err == nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return err
}

func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := filepath.Join(source, obj.Name())
		destinationfilepointer := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
		}
	}
	return
}

func ChownRecursively(root string, uId int, gId int) error {
	// Change file ownership.
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		err = os.Chown(path, uId, gId)
		if err != nil {
			return err
		}
		return nil
	})
}
