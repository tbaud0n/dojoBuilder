package dojoBuilder

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

type Config struct {
	BuildMode    bool   // Use dojo build if true
	SrcDir       string // Absolute path of the src js dir
	DestDir      string // Absolute path where the output files will be placed
	BuildConfigs map[string]BuildConfig
}

// dojoConfigRelPath is the path of the file containing the dojoConfig JSON. It is  relative to SrcDir.
const dojoConfigRelPath = "app/dojoConfig.json"

func Run(c *Config, names []string) (err error) {
	if c.DestDir == "" {
		return errors.New("No DestDir defined in config")
	}

	if _, err = os.Stat(c.DestDir); os.IsNotExist(err) {
		if err = os.MkdirAll(c.DestDir, 0754); err != nil {
			return
		}
	}

	if c.BuildMode {
		err = build(c, names)
	} else {
		err = installFiles(c)
	}

	return
}

func GetDojoConfig(c *Config) (template.JS, error) {
	dojoConfigFilePath := fmt.Sprintf("%s/%s", c.DestDir, dojoConfigRelPath)

	b, err := ioutil.ReadFile(dojoConfigFilePath)

	if err != nil {
		return "", err
	}

	return template.JS(b), err
}
