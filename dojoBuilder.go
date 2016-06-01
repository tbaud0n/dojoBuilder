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
	installDir        string // Private directory to copy src for build mode
	DojoConfigRelPath string // Path (relative to SrcDir) of the file containing the dojoConfig JSON
	BuildConfigs      map[string]BuildConfig
}

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

	if err = installFiles(c); err != nil {
		return
	}

	if c.BuildMode {
		err = build(c, names)
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
