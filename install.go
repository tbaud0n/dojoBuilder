package dojoBuilder

import (
	"os"
	"path/filepath"
	"syscall"
)

// DefaultInstallExcludeFunc skips .git folder, .gitignore  and .gitattributes files
var DefaultInstallExcludeFunc = func(path string, f os.FileInfo) bool {
	var skippedFiles []string = []string{".gitignore", ".gitattributes"}
	var skippedDirs []string = []string{".git"}

	var skipped []string

	if f.IsDir() {
		skipped = skippedDirs
	} else {
		skipped = skippedFiles
	}

	for _, skippedFile := range skippedFiles {
		if skippedFile == f.Name() {
			return true
		}
	}

	return false
}

func (c *Config) installFiles() (err error) {

	// Delete obsolete symlink and folders
	err = filepath.Walk(c.DestDir, func(path string, f os.FileInfo, err error) (_err error) {
		var forceRemove bool

		if path == c.DestDir {
			return nil
		}

		srcPath := c.SrcDir + path[len(c.DestDir):]

		if excludeFunc(srcPath, f) {
			forceRemove = true
		}

		if _, err = os.Stat(srcPath); os.IsNotExist(err) || forceRemove {
			if err = os.RemoveAll(path); err != nil {
				return err
			}
		}

		return
	})

	if err != nil {
		return
	}

	// Create new symlinks and folders
	err = filepath.Walk(c.SrcDir, func(path string, f os.FileInfo, err error) (_err error) {
		if path == c.SrcDir {
			return nil
		}

		newPath := c.DestDir + path[len(c.SrcDir):]
		if _, err = os.Stat(newPath); err == nil {
			return
		}

		isDir := f.IsDir()

		if excludeFunc(path, f) {
			if isDir {
				return filepath.SkipDir
			}
			return
		} else if isDir {
			if _err = os.Mkdir(newPath, 0754); _err != nil {
				return
			}
		} else if _err = os.Link(path, c.DestDir+path[len(c.SrcDir):]); _err != nil {
			return
		}

		st := f.Sys().(*syscall.Stat_t)

		os.Chown(newPath, int(st.Uid), int(st.Gid))

		return
	})

	return
}
