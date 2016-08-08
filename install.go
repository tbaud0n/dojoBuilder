package dojoBuilder

import (
	"os"
	"path/filepath"
	"syscall"
)

var (
	installExcludeFunc ExcludeFunc = func(path string, f os.FileInfo) (bool, error) {
		return false, nil
	}

	// DefaultInstallExcludeFunc skips .git folder, .gitignore  and .gitattributes files
	DefaultInstallExcludeFunc = func(path string, f os.FileInfo) (bool, error) {
		var skippedFilesPatterns []string = []string{`\.gitignore`, `\.gitattributes`}
		var skippedDirsPatterns []string = []string{`\.git`}

		if f.IsDir() {
			return IsMatchSliceMember(skippedDirsPatterns, path)
		}

		return IsMatchSliceMember(skippedFilesPatterns, path)
	}
)

func SetInstallExcludeFunc(exFunc ExcludeFunc) { installExcludeFunc = exFunc }

func (c *Config) installFiles() (err error) {

	// Delete obsolete symlink and folders
	err = filepath.Walk(c.DestDir, func(path string, f os.FileInfo, err error) (_err error) {
		var forceRemove bool

		if path == c.DestDir {
			return nil
		}

		srcPath := c.SrcDir + path[len(c.DestDir):]

		if forceRemove, _err = installExcludeFunc(srcPath, f); _err != nil {
			return _err
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

		if skip, _err := installExcludeFunc(path, f); _err != nil {
			return err
		} else if skip {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		} else if isDir {
			if _err = os.Mkdir(newPath, 0754); _err != nil {
				return _err
			}
		} else if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			origPath, _err := filepath.EvalSymlinks(path)
			if _err != nil {
				return _err
			}

			// fmt.Printf("Path : %s\nPoints to : %s\n\n", path, origPath)
			if _err = os.Symlink(origPath, c.DestDir+path[len(c.SrcDir):]); _err != nil {
				return _err
			}
		} else {
			if _err = os.Link(path, c.DestDir+path[len(c.SrcDir):]); _err != nil {
				return _err
			}
		}

		st := f.Sys().(*syscall.Stat_t)

		os.Chown(newPath, int(st.Uid), int(st.Gid))

		return
	})

	return
}
