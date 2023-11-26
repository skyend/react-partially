package plugins

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/fatih/color"
	"os"
)

func CreateEsbuildPluginBuildingStatusHook(onStart func(), onEnd func(result *api.BuildResult)) api.Plugin {

	red := color.New(color.FgRed).Add(color.Underline)
	green := color.New(color.FgHiGreen).Add(color.Underline)

	return api.Plugin{
		Name: "BuildingStatusHook",
		Setup: func(build api.PluginBuild) {
			build.OnStart(func() (api.OnStartResult, error) {
				fmt.Println("build started")
				onStart()
				return api.OnStartResult{}, nil
			})

			build.OnEnd(func(result *api.BuildResult) (api.OnEndResult, error) {
				green.Fprintf(os.Stderr, "build ended with %d errors\n", len(result.Errors))

				for _, err := range result.Errors {
					red.Println(err.Text)
				}

				onEnd(result)
				return api.OnEndResult{}, nil
			})

		},
	}

}
