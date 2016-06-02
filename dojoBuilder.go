package dojoBuilder

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	BuildMode         bool   // Use dojo build if true
	SrcDir            string // Absolute path of the src js dir
	DestDir           string // Absolute path where the output files will be placed
	Bin               string // Name of the bin used to build dojo (optional) [node, node-debug, java]
	DojoConfigRelPath string // Path (relative to SrcDir) of the file containing the dojoConfig JSON
	BuildConfigs      map[string]BuildConfig
}

// ExcludeFunc is called
// It allows ignore some folder when linking source files to DestDir
type ExcludeFunc func(path string, f os.FileInfo) (bool, error)

func Run(c *Config, names []string, reset bool) (err error) {
	if c.DestDir == "" {
		return errors.New("No DestDir defined in config")
	}

	if _, err = os.Stat(c.DestDir); os.IsNotExist(err) {
		if err = os.MkdirAll(c.DestDir, 0754); err != nil {
			return
		}
	}

	if reset {
		filepath.Walk(c.DestDir, func(path string, f os.FileInfo, err error) (_err error) {
			if path != c.DestDir {
				_err = os.RemoveAll(path)
			}
			return
		})
	}

	if c.BuildMode {
		err = c.build(names)
	} else {
		err = c.installFiles()
	}

	return
}

func GetDojoConfig(c *Config) (template.JS, error) {
	dojoConfigFilePath := fmt.Sprintf("%s/%s", c.DestDir, c.DojoConfigRelPath)

	b, err := ioutil.ReadFile(dojoConfigFilePath)

	if err != nil {
		return "", err
	}

	return template.JS(b), err
}
