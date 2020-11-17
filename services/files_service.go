package services

import (
	"bufio"
	"fmt"
	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/dto"
	"os"
	"path/filepath"
	"strings"
)

func isIgnoredPath(path string, ignoredPaths []string) bool {
	for _, dirName := range ignoredPaths {
		if strings.Contains(path, dirName) {
			return true
		}
	}

	return false
}

//AnalyseTheCode method loads the available project features in the memory
func AnalyseTheCode(src string, ext string, ignored []string) (foundResults map[string]dto.IndexedFile, pathIndex map[string]dto.IndexedFile, importsIndex map[string][]string, err error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return foundResults, pathIndex, importsIndex, err
	}

	foundResults = map[string]dto.IndexedFile{}
	pathIndex = map[string]dto.IndexedFile{}
	importsIndex = map[string][]string{}

	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		if isIgnoredPath(path, ignored) {
			return nil
		}

		if ext != "" {
			if filepath.Ext(path) != ext {
				return nil
			}
		}

		indexedFile, err := AnalyseFile(src, path)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error during file `%s` parse : %s", path, err.Error()))
			return err
		}

		if indexedFile.Path == "" {
			return nil
		}

		if indexedFile.MainEntrypoint != "" {
			foundResults[indexedFile.MainEntrypoint] = indexedFile
		}

		pathIndex[indexedFile.Path] = indexedFile

		if len(indexedFile.OtherImports) > 0 && indexedFile.MainEntrypoint != "" {
			for _, importSrc := range indexedFile.OtherImports {
				importsIndex[importSrc] = append(importsIndex[importSrc], indexedFile.MainEntrypoint)
			}
		}

		return nil
	}); err != nil {
		return foundResults, pathIndex, importsIndex, err
	}

	return
}

//AnalyseFile find the features in the destination file path
//basePath - will be used for generation of relation path
//filePath - the actual absolute file path
func AnalyseFile(basePath string, filePath string) (indexedFile dto.IndexedFile, err error) {
	fsFile, err := os.Open(filePath)
	if err != nil {
		return
	}

	// Splits on newlines by default.
	scanner := bufio.NewScanner(fsFile)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	line := 1

	indexedFile.Path = filePath
	indexedFile.OtherImports = map[string]string{}
	mainEntrypointSrc := ""
	entryPointSrc := ""
	for scanner.Scan() {
		text := scanner.Text()
		//We check for features annotation
		if strings.Contains(text, FeatureAlias) {
			featureType, err := extractFeatureType(scanner.Text())
			if err != nil || featureType == "" {
				continue
			}

			relPath, err := filepath.Rel(basePath, filePath)
			if err != nil {
				return indexedFile, err
			}

			feature := dto.Feature{}
			feature.FilePath = relPath
			feature.Line = line
			feature.Name = featureType
			indexedFile.Features = append(indexedFile.Features, feature)
			setToFoundFeatures(feature)
		}

		if entryPointSrc == "" {
			if mainEntrypointSrc == "" {
				mainEntrypointSrc, err = analysis.An.ExtractMainEntrypointSrc(text)
				if err != nil {
					return indexedFile, err
				}
			}

			//If we found main entrypoint import, that means we can have an entrypoint object title.
			if mainEntrypointSrc != "" {
				mainEntrypointTitle, err := analysis.An.FindMainEntrypointTitle(text)
				if err != nil {
					return indexedFile, err
				}

				if mainEntrypointTitle != "" {
					mainEntrypointSrc = analysis.An.GenerateEntryPointSrc(mainEntrypointSrc, mainEntrypointTitle)
					indexedFile.MainEntrypoint = mainEntrypointSrc
					entryPointSrc = mainEntrypointSrc
				}
			}
		}

		importSrc, err := analysis.An.ExtractImportUsage(text)
		if err != nil {
			return indexedFile, err
		}

		if indexedFile.OtherImports[importSrc] == "" && importSrc != "" {
			indexedFile.OtherImports[importSrc] = importSrc
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		return indexedFile, err
	}

	if err := fsFile.Close(); err != nil {
		return indexedFile, err
	}

	return indexedFile, nil
}

//GetIgnoredFilePaths used for generation of the list of ignored file paths, which will be ignored during the features search
func GetIgnoredFilePaths(path string) (files []string, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	fsFile, err := os.Open(path)
	if err != nil {
		return
	}

	// Splits on newlines by default.
	scanner := bufio.NewScanner(fsFile)

	for scanner.Scan() {
		if scanner.Text() != "" {
			files = append(files, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return files, err
	}

	return
}
