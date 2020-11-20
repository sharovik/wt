package services

import (
	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/dto"
	"regexp"
	"strings"
)

const (
	FeatureAlias     = "@featureType"
	regexFeatureType = `(?i)(?:@featureType)(.*)`
)

var (
	//PF object of ProjectFeatures configuration
	PF = ProjectFeatures{
		FoundFeatures:       map[string][]dto.Feature{},
		FoundFeaturesByFile: map[string][]dto.Feature{},
	}
)

//ProjectFeatures this is main struct which is used for generation of WT configuration
type ProjectFeatures struct {
	FoundFeatures       map[string][]dto.Feature `json:"found_features"`
	FoundFeaturesByFile map[string][]dto.Feature `json:"found_features_by_file"`
}

//CleanUpGlobalVars method which cleanup the project features object attributes
func CleanUpGlobalVars() {
	PF.FoundFeatures = map[string][]dto.Feature{}
	PF.FoundFeaturesByFile = map[string][]dto.Feature{}
}

func setToFoundFeatures(feature dto.Feature) {
	PF.FoundFeatures[feature.Name] = append(PF.FoundFeatures[feature.Name], feature)
	doesNotExists := true
	for _, item := range PF.FoundFeaturesByFile[feature.FilePath] {
		if item == feature {
			doesNotExists = false
		}
	}

	if doesNotExists {
		PF.FoundFeaturesByFile[feature.FilePath] = append(PF.FoundFeaturesByFile[feature.FilePath], feature)
	}
}

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
		relativePath := strings.ReplaceAll(relatedFile.Path, absolutePath + "/", "")
		totalFeaturesTouched[relativePath] = relatedFile.Features

		for _, feature := range relatedFile.Features {
			setToFoundFeatures(feature)
		}

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
		relativePath := strings.ReplaceAll(touchedFile.Path, absolutePath + "/", "")
		if relativePath == "" {
			continue
		}

		if len(touchedFile.Features) > 0 {
			totalFeaturesTouched[relativePath] = touchedFile.Features
			for _, feature := range touchedFile.Features {
				setToFoundFeatures(feature)
			}
		} else {
			toBeChecked[relativePath] = touchedFile
		}
	}

	return
}