package dojoBuilder

import (
	"os"
	"path/filepath"
	"syscall"
)

// ExcludeDirFunc is called for non-built install
// It allows ignore some folder when linking source files to DestDir
type ExcludeDirFunc func(path string, f os.FileInfo) bool

// ExcludeFileFunc is called for non-built install
// It allows ignore some files when linking source files to DestDir
type ExcludeFileFunc func(path string, f os.FileInfo) bool

var (
	excludeDirFunc ExcludeDirFunc = func(path string, f os.FileInfo) bool {
		return false
	}
	excludeFileFunc ExcludeFileFunc = func(path string, f os.FileInfo) bool {
		return false
	}

	// DefaultExcludeDirFunc skips .git folder when installing dojo with non-built config
	DefaultExcludeDirFunc = func(path string, f os.FileInfo) (skip bool) {
		skip = (f.Name() == ".git")

		return
	}

	DefaultExcludeFileFunc = func(path string, f os.FileInfo) bool {
		var skippedFiles []string = []string{".gitignore", ".gitattributes"}

		for _, skippedFile := range skippedFiles {
			if skippedFile == f.Name() {
				return true
			}
		}

		return false
	}
)

func SetExcludeDirFunc(exDirFunc ExcludeDirFunc) { excludeDirFunc = exDirFunc }

func SetExcludeFileFunc(exFileFunc ExcludeFileFunc) { excludeFileFunc = exFileFunc }

func installFiles(c *Config) (err error) {
	c.installDir = c.DestDir

	if c.BuildMode {
		c.installDir = c.DestDir + "/tmp"
		if _, err = os.Stat(c.installDir); os.IsNotExist(err) {
			if err = os.MkdirAll(c.installDir, 0754); err != nil {
				return
			}
		}
	}

	// Delete obsolete symlink and folders
	err = filepath.Walk(c.installDir, func(path string, f os.FileInfo, err error) (_err error) {
		var forceRemove bool

		if path == c.installDir {
			return nil
		}

		srcPath := c.SrcDir + path[len(c.installDir):]

		if f.IsDir() {
			if excludeDirFunc(srcPath, f) {
				forceRemove = true
			}
		} else if excludeFileFunc(srcPath, f) {
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

		newPath := c.installDir + path[len(c.SrcDir):]
		if _, err = os.Stat(newPath); err == nil {
			return
		}
		if f.IsDir() {
			if excludeDirFunc(path, f) {
				return filepath.SkipDir
			} else if _err = os.Mkdir(newPath, 0754); _err != nil {
				return
			}
		} else if excludeFileFunc(path, f) {
			return
		} else if _err = os.Link(path, c.installDir+path[len(c.SrcDir):]); _err != nil {
			return
		}

		st := f.Sys().(*syscall.Stat_t)

		os.Chown(newPath, int(st.Uid), int(st.Gid))

		return
	})

	return
}
