package services

import (
	"regexp"
	"strings"
)

const (
	featureAlias     = "@featureType"
	regexFeatureType = `(?i)(?:@featureType)(.*)`
)

var (
	//PF object of ProjectFeatures configuration
	PF = ProjectFeatures{
		FoundFeatures:       map[string][]Feature{},
		FoundFeaturesByFile: map[string][]Feature{},
	}
)

//ProjectFeatures this is main struct which is used for generation of WT configuration
type ProjectFeatures struct {
	FoundFeatures       map[string][]Feature `json:"found_features"`
	FoundFeaturesByFile map[string][]Feature `json:"found_features_by_file"`
}

//Feature the struct for unique feature
type Feature struct {
	Name     string
	FilePath string
	Line     int
}

//CleanUpGlobalVars method which cleanup the project features object attributes
func CleanUpGlobalVars() {
	PF.FoundFeatures = map[string][]Feature{}
	PF.FoundFeaturesByFile = map[string][]Feature{}
}

func setToFoundFeatures(feature Feature) {
	PF.FoundFeatures[feature.Name] = append(PF.FoundFeatures[feature.Name], feature)
	PF.FoundFeaturesByFile[feature.FilePath] = append(PF.FoundFeaturesByFile[feature.FilePath], feature)
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
