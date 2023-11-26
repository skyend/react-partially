package lib

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func GetRouteEntryFilePath(packageDir, distDir string) string {
	return path.Join(packageDir, distDir, RouteEntryFileName)
}

func GenerateRouteEntryFile(packageDir, distDir string) error {
	return os.WriteFile(GetRouteEntryFilePath(packageDir, distDir), []byte("s"), os.ModePerm)
}

func UpdateRouteEntryFile(packageDir, distDir string, routes map[string]Route) error {
	importString := "\n"
	for _, route := range routes {

		importString += fmt.Sprintf(
			"export const %s = { path: \"%s\", module: import(\"%s\") };",
			strings.ReplaceAll(route.RoutePath, "/", "_"),
			route.RoutePath,
			route.RouteFilePathFromDist,
		)
	}
	importString += "\n"
	return os.WriteFile(GetRouteEntryFilePath(packageDir, distDir), []byte(importString), os.ModePerm)
}
