package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"flag"

	"github.com/tbaud0n/dojoBuilder"
)

const pageTemplate = `<!doctype html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <link rel="stylesheet" href="pkg/dijit/themes/claro/claro.css">
</head>
<body class="claro">
  <script type="text/javascript">var dojoConfig = {{getDojoConfig}};</script>
{{if .buildMode}}
  <script type="text/javascript" src="pkg/app/main.js"></script>
{{else}}
  <script type="text/javascript" src="pkg/dojo/dojo.js"></script>
{{end}}
  <script type="text/javascript">require(['app/main']);</script>
</body>
</html>`

var builderConfig *dojoBuilder.Config

func getDojoConfig() template.JS {
	dc, err := dojoBuilder.GetDojoConfig(builderConfig)
	if err != nil {
		log.Fatal(err)
	}

	return dc
}

func handler(w http.ResponseWriter, r *http.Request) {
	templateFuncs := template.FuncMap{
		"getDojoConfig": getDojoConfig,
	}
	t := template.Must(template.New("page").Funcs(templateFuncs).Parse(pageTemplate))

	t.Execute(w, nil)
}

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	buildMode := flag.Bool("buildMode", false, "Use build mode of dojoBuilder")

	flag.Parse()

	builderConfig = &dojoBuilder.Config{
		BuildMode:         *buildMode,
		SrcDir:            dir + "/client",
		DestDir:           dir + "/pkg",
		DojoConfigRelPath: "app/dojoConfig.json",
		BuildConfigs: map[string]dojoBuilder.BuildConfig{
			"default": dojoBuilder.BuildConfig{
				RemoveUncompressed:    true,
				RemoveConsoleStripped: true,
				Packages: []dojoBuilder.Package{
					dojoBuilder.Package{Name: "dojo", Location: "dojo"},
					dojoBuilder.Package{Name: "dijit", Location: "dijit"},
					dojoBuilder.Package{Name: "dojox", Location: "dojox"},
					dojoBuilder.Package{Name: "app", Location: "app"},
				},
				Layers: map[string]dojoBuilder.Layer{
					"app/main": dojoBuilder.Layer{
						Include:    []string{"dojo/dojo", "dijit/dijit", "app/main"},
						CustomBase: true,
						Boot:       true,
					},
				},
				LayerOptimize:     "closure",
				CssOptimize:       "comments",
				Mini:              true,
				StripConsole:      "warn",
				SelectorEngine:    "lite",
				StaticHasFeatures: map[string]dojoBuilder.Feature{
				// "dojo-trace-api":        false,
				// "dojo-log-api":          false,
				// "dojo-publish-privates": false,
				// "dojo-sync-loader":      false,
				// "dojo-xhr-factory":      false,
				// "dojo-test-sniff":       false,
				},
				UseSourceMaps: false,
			},
		},
	}
}

func main() {

	if err := dojoBuilder.Run(builderConfig, nil, true); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handler)
	http.Handle("/pkg/", http.StripPrefix("/pkg/", http.FileServer(http.Dir(builderConfig.DestDir))))
	http.ListenAndServe(":8080", nil)
}
