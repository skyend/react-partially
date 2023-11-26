package main

import (
	"encoding/json"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/gofiber/fiber/v2"
	"nwn-build-system/builder/lib"
	"nwn-build-system/builder/plugins"
	"os"
	"path"
	"path/filepath"
	"sync"
)

const PackageDir = "apps/app-ui"
const PagesDir = "src/pages"
const StaticDir = "static"
const DistDir = "dist"

// Namespace 가 지정되면 static files 에 prefix 가 부여 되며
// Namespace 가 지정된 채로 빌드된 Page 는 해당 Namespace 의 static file 만 참조 한다
// Namespace 를 지정 함으로써 특정한 페이지들을 독립적으로 빌드 하고 해당 페이지들이 참조하는 Static 파일들을 독립적으로 실행 할 수 있다.
const Namespace = "default"

func main() {
	relativeDistPath := path.Join(PackageDir, DistDir)
	relativePagesDirPath := path.Join(PackageDir, PagesDir)

	os.MkdirAll(relativeDistPath, os.ModePerm)
	lib.GenerateRouteEntryFile(PackageDir, DistDir)

	incRouteManager := lib.IncrementalRouteManager{Routes: map[string]lib.Route{}}

	pageWatcher := lib.PageWatcher{}
	pageWatcher.Init()
	pageWatcher.SetPageDir(PackageDir, PagesDir)
	go pageWatcher.Run()
	defer pageWatcher.Close()

	buildWait := sync.WaitGroup{}
	notifyPlugin := plugins.CreateEsbuildPluginBuildingStatusHook(
		func() {
			buildWait.Add(1)
		},
		func(result *api.BuildResult) {
			for _, file := range result.OutputFiles {
				fmt.Println("Write file... ", file.Path)
				dir := path.Dir(file.Path)
				os.MkdirAll(dir, os.ModePerm)
				os.WriteFile(file.Path, file.Contents, os.ModePerm)
			}

			// Reflect build meta to struct
			meta := lib.EsbuildBuiltMeta{}
			err := json.Unmarshal([]byte(result.Metafile), &meta)
			if err != nil {
				panic(err)
			}

			metaDataPath := path.Join(PackageDir, DistDir, "meta.json")
			os.WriteFile(metaDataPath, []byte(result.Metafile), os.ModePerm)

			lib.PostBuildOutputProcessor(relativePagesDirPath, Namespace, meta)

			buildWait.Done()
		},
	)

	esbuildCtx, err := api.Context(lib.GenerateBuildOptions(
		PackageDir,
		PagesDir,
		DistDir,
		Namespace,
		[]api.Plugin{notifyPlugin}),
	)
	if err != nil {
		panic(err)
	}

	esbuildWatchError := esbuildCtx.Watch(api.WatchOptions{})
	if esbuildWatchError != nil {
		panic(esbuildWatchError)
	}

	app := fiber.New()

	app.Static("/", path.Join(PackageDir, StaticDir))

	app.Use(func(c *fiber.Ctx) error {
		urlPath := c.Path()

		ext := filepath.Ext(urlPath)
		if ext != "" {
			return c.Next()
		}

		tsxPathFromCwd := lib.RoutingPathToPageFileFullPath(PackageDir, PagesDir, urlPath)

		// dist/ 디렉토리에서 [Page].tsx 파일 까지의 상대 경로
		pageFilePathFromDist, err := filepath.Rel(relativeDistPath, tsxPathFromCwd)
		if err != nil {
			return c.Next()
		}

		pageFilePathFromPageDir, err := filepath.Rel(relativePagesDirPath, tsxPathFromCwd)
		if err != nil {
			return c.Next()
		}

		routingPath := lib.PageFilePathToRoutingPath(pageFilePathFromPageDir)

		stat, err := os.Stat(tsxPathFromCwd)
		if err != nil || stat.IsDir() {
			fmt.Println(fmt.Errorf("failed to build page %s", tsxPathFromCwd).Error())
			return c.Status(500).Send([]byte("failed to build"))
		}

		route := lib.Route{
			RouteFilePathFromDist: pageFilePathFromDist,
			RouteFilePath:         pageFilePathFromPageDir,
			RoutePath:             routingPath,
		}

		// 라우트 목록에 포함되어 있으면,
		// 이미 Esbuild 에 의해 watching 되고 있으므로 추가적으로 auto generated routeEntry 를 업데이트 할 필요가 없다
		if !incRouteManager.HasRoute(route) {
			incRouteManager.IncludeRoute(route)
			lib.UpdateRouteEntryFile(PackageDir, DistDir, incRouteManager.Routes)
		}

		fmt.Println("path:", tsxPathFromCwd, pageFilePathFromDist, routingPath)

		buildWait.Wait()

		htmlFilePath := path.Join(PackageDir, DistDir, lib.RoutingPathToHtmlPath(routingPath))

		return c.SendFile(htmlFilePath)
	})

	app.Listen(":3000")
}
