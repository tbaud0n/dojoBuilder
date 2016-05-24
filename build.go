package dojoBuilder

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const profileTemplate = `var profile = {{.}};`

type BuildConfig struct {
	RemoveUncompressed bool `json:"removeUncompressed,omitempty"` // Remove uncompressed js files after build

	BasePath    string           `json:"basePath"`
	ReleaseDir  string           `json:"releaseDir"`
	ReleaseName string           `json:"releaseName,omitempty"`
	Action      string           `json:"action"`
	Packages    []Package        `json:"packages"`
	Layers      map[string]Layer `json:"layers"`

	LayerOptimize     string             `json:"layerOptimize,omitempty"`
	Optimize          string             `json:"optimize,omitempty"`
	CssOptimize       string             `json:"cssOptimize,omitempty"`
	Mini              bool               `json:"mini,omitempty"`
	StripConsole      string             `json:"stripConsole,omitempty"`
	SelectorEngine    string             `json:"selectorEngine,omitempty"`
	StaticHasFeatures map[string]Feature `json:"staticHasFeatures,omitempty"`
	UseSourceMaps     bool               `json:"useSourceMaps"` // Build generate source maps
}

type Package struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Layer struct {
	Boot       bool     `json:"boot"`
	CustomBase bool     `json:"customBase"`
	Include    []string `json:"include,omitempty"`
	Exclude    []string `json:"exclude,omitempty"`
}

type Feature bool

func (f Feature) MarshalJSON() ([]byte, error) {
	var v uint8 = 0
	if bool(f) {
		v = 1
	}
	return json.Marshal(v)
}

func (c *Config) generateBuildProfile(name string) (profileFullPath string, err error) {
	if bc, ok := c.BuildConfigs[name]; !ok {
		return "", errors.New("No build config found with name '" + name + "'")
	} else {
		if bc.Action == "" {
			bc.Action = "release"
		}

		profilePath := c.SrcDir + "/profiles/"
		os.MkdirAll(profilePath, 0754)

		profileFullPath = profilePath + name + ".profile.js"

		bc.BasePath = "../"

		if bc.ReleaseDir, err = filepath.Rel(c.SrcDir+`/`+bc.BasePath+"empty", c.DestDir); err != nil {
			return "", err
		}

		if j, err := json.Marshal(bc); err != nil {
			return "", err
		} else {
			if f, err := os.OpenFile(profileFullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664); err != nil {
				return "", err
			} else {
				t := template.Must(template.New("profileTemplate").Parse(profileTemplate))
				err = t.Execute(f, string(j))
			}
		}

		return profileFullPath, err
	}
}

func build(c *Config, names []string) (err error) {
	var profilePath string

	if len(names) == 0 {
		for n, _ := range c.BuildConfigs {
			names = append(names, n)
		}
	}

	for _, n := range names {
		fmt.Printf("Generating %s build\n", n)

		profilePath, err = c.generateBuildProfile(n)
		if err != nil {
			return
		}

		if err = executeBuildProfile(c, profilePath); err != nil {
			return
		}
	}

	return
}

func executeBuildProfile(c *Config, profilePath string) (err error) {
	buildScriptPath := c.SrcDir + "/util/buildscripts/build.sh"

	cmd := exec.Command(buildScriptPath, "--profile", profilePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Run()
	if err != nil {
		return errors.New("Build command failed")
	}

	return
}
