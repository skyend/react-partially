package lib

import (
	"path"
	"regexp"
)

// PageFilePath :pagefile(index.tsx) path from page directory
type PageFilePath = string

type RoutingPath = string

var RemoveTsxExtensionRegExp = regexp.MustCompile("\\.tsx$")

func PageFilePathToRoutingPath(pageFilePathFromPageDir PageFilePath) RoutingPath {
	rawPagePath := RemoveTsxExtensionRegExp.ReplaceAll([]byte(pageFilePathFromPageDir), []byte(""))
	pagePath := RoutingPath(rawPagePath)
	if pagePath == "index" {
		pagePath = "/"
	}

	return pagePath
}

func RoutingPathToHtmlPath(pagePath RoutingPath) string {
	if pagePath == "/" {
		return "index.html"
	}

	return pagePath + ".html"
}

func RoutingPathToPageFileFullPath(packagePath, pagesPath, urlPath string) string {
	if urlPath == "/" {
		return path.Join(packagePath, pagesPath, "index.tsx")
	}

	return path.Join(packagePath, pagesPath, urlPath+".tsx")
}
