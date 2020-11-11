package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ignoredPaths []string

func isIgnoredDir(path string) bool {
	for _, dirName := range ignoredPaths {
		if strings.Contains(path, dirName) {
			return true
		}
	}

	return false
}

//LoadAvailableFeaturesInDir method loads the available project features in the memory
func LoadAvailableFeaturesInDir(src string, ext string, pathToIgnoreFile string) (foundResults map[string][]Feature, err error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil, err
	}

	foundResults = map[string][]Feature{}
	ignoredPaths, err = GetIgnoredFilePaths(pathToIgnoreFile)
	if err != nil {
		return foundResults, err
	}

	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		if isIgnoredDir(path) {
			return nil
		}

		if ext != "" {
			if filepath.Ext(path) != ext {
				return nil
			}
		}

		features, err := FindFeaturesInFile(src, path)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error during file `%s` parse : %s", path, err.Error()))
			return err
		}

		if len(features) > 0 {
			foundResults[path] = features
		}

		return nil
	}); err != nil {
		return foundResults, err
	}

	return
}

//FindFeaturesInFile find the features in the destination file path
//basePath - will be used for generation of relation path
//filePath - the actual absolute file path
func FindFeaturesInFile(basePath string, filePath string) (features []Feature, err error) {
	fsFile, err := os.Open(filePath)
	if err != nil {
		return
	}

	// Splits on newlines by default.
	scanner := bufio.NewScanner(fsFile)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), featureAlias) {
			featureType, err := extractFeatureType(scanner.Text())
			if err != nil || featureType == "" {
				continue
			}

			relPath, err := filepath.Rel(basePath, filePath)
			if err != nil {
				return features, err
			}

			feature := Feature{}
			feature.FilePath = relPath
			feature.Line = line
			feature.Name = featureType
			features = append(features, feature)
			setToFoundFeatures(feature)
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		return features, err
	}

	if err := fsFile.Close(); err != nil {
		return features, err
	}

	return features, nil
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
		files = append(files, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return files, err
	}

	return
}
