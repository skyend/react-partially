package lib

import (
	"fmt"
	"path/filepath"
)

func PostBuildOutputProcessor(relativePagesDirPath string, namespace string, meta EsbuildBuiltMeta) {
	// @Todo generate pages by meta
	for chunkFilename, output := range meta.Outputs {
		if output.EntryPoint != "" {
			matched, err := filepath.Match(relativePagesDirPath+"/*", output.EntryPoint)
			if err != nil || !matched {
				continue
			}
			pageFilePath, err := filepath.Rel(relativePagesDirPath, output.EntryPoint)
			if err != nil {
				continue
			}
			if pageFilePath == "app.tsx" || pageFilePath == "browser.entry.tsx" {
				continue
			}

			fmt.Println("page", pageFilePath, "-->", chunkFilename)
			for _, importInfo := range output.Imports {
				fmt.Println("importing -> ", importInfo.Path)
			}
			//output.

			//					htmlFilePath := path.Join(PackageDir, DistDir, PagePathToHtmlPath(pagePath))
			//					os.WriteFile(htmlFilePath, []byte(`
			//<html>
			//	<head>
			//
			//	</head>
			//	<body>
			//		hello
			//		<script>
			//			window.__ROUTE={}
			//		</script>
			//	</body>
			//</html>
			//`), os.ModePerm)

		}
	}
}
