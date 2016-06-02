package dojoBuilder

import (
	"fmt"
	"io"
	"os"
)

func isStringInSlice(slice []string, s string) bool {
	for st, _ := range slice {
		if st == s {
			return true
		}
	}

	return false
}

func CopyDir(src string, dest string) (err error) {

	sfi, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if !sfi.IsDir() {
		return fmt.Errorf("CopyFile cannot copy a file %s", src)
	} else if sfi.Mode()&os.ModeSymlink != 0 {
		linkSrc, err = os.Readlink(src)
		if err != nil {
			return
		}
		return CopyDir(linkSrc, dest)
	}

	err = os.MkdirAll(dest, sfi.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(src)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		srcfilepointer := src + "/" + obj.Name()

		destfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			CopyDir(srcfilepointer, destfilepointer)
			if err != nil {
				return
			}
		} else {
			err = CopyFile(srcfilepointer, destfilepointer)
			if err != nil {
				return
			}
		}

	}
	return
}

// CopyFile copies a file from src to dest. If src and dest files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dest.
func CopyFile(src, dest string) (err error) {
	sfi, err := os.Lstat(src)
	if err != nil {
		return
	}

	if !sfi.Mode().IsRegular() {
		if sfi.IsDir() {
			return fmt.Errorf("CopyFile cannot copy a directory %s", src)
		} else if sfi.Mode()&os.ModeSymlink != 0 {
			linkSrc, err = os.Readlink(src)
			if err != nil {
				return
			}
			return CopyFile(linkSrc, dest)
		} else {
			return fmt.Errorf("CopyFile: non-regular source file %s (%q)", src, sfi.Mode().String())
		}
	}

	dfi, err := os.Lstat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dest); err == nil {
		return
	}
	err = copyFileContents(src, dest)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dest. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dest string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
