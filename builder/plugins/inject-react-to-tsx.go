package plugins

import (
	"encoding/json"
	"github.com/evanw/esbuild/pkg/api"
	"io/ioutil"
	"strings"
)

var ExampleOnLoadPlugin = api.Plugin{
	Name: "example",
	Setup: func(build api.PluginBuild) {
		// Load ".txt" files and return an array of words
		build.OnLoad(api.OnLoadOptions{Filter: `\.tsx$`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				text, err := ioutil.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}
				bytes, err := json.Marshal(strings.Fields(string(text)))
				if err != nil {
					return api.OnLoadResult{}, err
				}
				contents := string(bytes)
				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderJSON,
				}, nil
			})
	},
}
