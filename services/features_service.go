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

	//ProjectAlias - the constant for project annotation search
	ProjectAlias     = "@project"
	regexFeatureType = `(?i)(?:@featureType)(.*)`
	regexProject     = `(?i)(?:@project)(.*)`
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

func extractProjects(text string) (project []string, err error) {
	re, err := regexp.Compile(regexProject)

	if err != nil {
		return
	}

	matches := re.FindStringSubmatch(text)

	if len(matches) != 2 {
		return
	}

	projectAnnotationValue := strings.TrimSpace(matches[1])
	if "" == projectAnnotationValue {
		return
	}

	return strings.Split(projectAnnotationValue, ","), nil
}

//FindFeaturesInIndex method tries to find features in the prepared indexes
func FindFeaturesInIndex(diff []string, absolutePath string) (result dto.AnalysisResult) {
	var usage []string
	potentiallyTouched := map[string]dto.IndexedFile{}
	result.TotalFeaturesTouched = map[string][]dto.Feature{}
	result.ToBeChecked = map[string]dto.IndexedFile{}
	for _, file := range diff {
		if analysis.AnalysedPathsIndex[file].Path == "" {
			continue
		}

		relatedFile := analysis.AnalysedPathsIndex[file]
		relativePath := strings.ReplaceAll(relatedFile.Path, absolutePath+"/", "")
		result.TotalFeaturesTouched[relativePath] = relatedFile.Features
		result.AppendProjects(relatedFile.RelatedProjects)

		usage = []string{}
		usage = analysis.FindUsage(relatedFile.MainEntrypoint, usage, 0)
		if len(usage) > 0 {
			for _, usedInEntrypoint := range usage {
				usedInFile := analysis.AnalysedEntrypointsIndex[usedInEntrypoint]
				usedInFile.UsedIn = append(usedInFile.UsedIn, relatedFile)
				potentiallyTouched[usedInEntrypoint] = usedInFile
				result.AppendProjects(usedInFile.RelatedProjects)
			}
		}
	}

	for _, touchedFile := range potentiallyTouched {
		relativePath := strings.ReplaceAll(touchedFile.Path, absolutePath+"/", "")
		if relativePath == "" {
			continue
		}

		if len(touchedFile.Features) > 0 {
			result.TotalFeaturesTouched[relativePath] = touchedFile.Features
		} else {
			result.ToBeChecked[relativePath] = touchedFile
		}
	}

	return
}
