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

func LoadAvailableFeaturesInDir(src string, ext string, pathToIgnoreFile string) (foundResults map[string][]Feature, err error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil, err
	}

	foundResults = map[string][]Feature{}
	ignoredPaths, err = GetGitIgnoreFilePaths(pathToIgnoreFile)
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

func FindFeaturesInFile(basePath string, filePath string) (features []Feature, err error) {
	fsFile, err := os.Open(filePath)
	if err != nil {
		return
	}

	// Splits on newlines by default.
	scanner := bufio.NewScanner(fsFile)

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

func GetGitIgnoreFilePaths(path string) (files []string, err error) {
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