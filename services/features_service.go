package services

import (
	"regexp"
	"strings"

	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/dto"
)

const (
	//FeatureAlias - the constant which identifies the feature alias which our parser will try to find in the strings
	FeatureAlias     = "@featureType"
	regexFeatureType = `(?i)(?:@featureType)(.*)`
)

func extractFeatureType(text string) (featureType string, err error) {
	re, err := regexp.Compile(regexFeatureType)

	if err != nil {
		return
	}

	matches := re.FindStringSubmatch(text)

	if len(matches) != 2 {
		return
	}

	return strings.TrimSpace(matches[1]), nil
}

//FindFeaturesInIndex method tries to find features in the prepared indexes
func FindFeaturesInIndex(diff []string, absolutePath string) (totalFeaturesTouched map[string][]dto.Feature, toBeChecked map[string]dto.IndexedFile) {
	var usage []string
	potentiallyTouched := map[string]dto.IndexedFile{}
	totalFeaturesTouched = map[string][]dto.Feature{}
	toBeChecked = map[string]dto.IndexedFile{}
	for _, file := range diff {
		if analysis.AnalysedPathsIndex[file].Path == "" {
			continue
		}

		relatedFile := analysis.AnalysedPathsIndex[file]
		relativePath := strings.ReplaceAll(relatedFile.Path, absolutePath+"/", "")
		totalFeaturesTouched[relativePath] = relatedFile.Features

		usage = []string{}
		usage = analysis.FindUsage(relatedFile.MainEntrypoint, usage, 0)
		if len(usage) > 0 {
			for _, usedInEntrypoint := range usage {
				usedInFile := analysis.AnalysedEntrypointsIndex[usedInEntrypoint]
				usedInFile.UsedIn = append(usedInFile.UsedIn, relatedFile)
				potentiallyTouched[usedInEntrypoint] = usedInFile
			}
		}
	}

	for _, touchedFile := range potentiallyTouched {
		relativePath := strings.ReplaceAll(touchedFile.Path, absolutePath+"/", "")
		if relativePath == "" {
			continue
		}

		if len(touchedFile.Features) > 0 {
			totalFeaturesTouched[relativePath] = touchedFile.Features
		} else {
			toBeChecked[relativePath] = touchedFile
		}
	}

	return
}
