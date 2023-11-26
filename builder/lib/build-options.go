package lib

import (
	"github.com/evanw/esbuild/pkg/api"
	"path"
)

func GetOutputChunkDir(Package, DistDir, Namespace string) string {
	if Namespace != "" {
		return path.Join(Package, DistDir, Namespace+"-"+"chunks")
	}

	return path.Join(Package, DistDir, "chunks")
}

func GenerateBuildOptions(Package, Pages, DistDir, Namespace string, plugins []api.Plugin) api.BuildOptions {
	return api.BuildOptions{
		EntryPointsAdvanced: []api.EntryPoint{
			{
				OutputPath: "_browser.entry",
				InputPath:  path.Join(Package, Pages, "_browser.entry.tsx"),
			},
			{
				OutputPath: "_app",
				InputPath:  path.Join(Package, Pages, "_app.tsx"),
			},
			{
				OutputPath: "_auto_routes_",
				InputPath:  GetRouteEntryFilePath(Package, DistDir),
			},
		},

		Metafile:    true,
		Bundle:      true,
		Outdir:      GetOutputChunkDir(Package, DistDir, Namespace),
		Splitting:   true,
		Format:      api.FormatESModule,
		TreeShaking: api.TreeShakingTrue,
		//MinifyIdentifiers: true,
		//MinifySyntax:      true,
		//MinifyWhitespace:  true,
		Plugins: plugins,
	}
}
