package printout

import (
	"fmt"
	"github.com/sharovik/wt/dto"
	"strings"
)

//FeaturesPrintout the full printout struct
type FeaturesPrintout struct {
	AbsolutePath string
	PrintToBeCheckedDetails bool
	TotalFeaturesTouched map[string][]dto.Feature
	ToBeChecked map[string]dto.IndexedFile
}

func (s FeaturesPrintout) SetAbsolutePath(path string) {
	s.AbsolutePath = path
}

func (s FeaturesPrintout) SetTotalFeaturesTouched(features map[string][]dto.Feature) {
	s.TotalFeaturesTouched = features
}

func (s FeaturesPrintout) SetToBeChecked(files map[string]dto.IndexedFile) {
	s.ToBeChecked = files
}

func (s FeaturesPrintout) GetToBeChecked() map[string]dto.IndexedFile {
	return s.ToBeChecked
}

func (s FeaturesPrintout) ToBeCheckedText() string {
	return generateToBeCheckedText(s)
}

func (s FeaturesPrintout) WithToBeCheckedDetails() {
	s.PrintToBeCheckedDetails = true
}

func (s FeaturesPrintout) IsToBeCheckedDetailsEnabled() bool {
	return s.PrintToBeCheckedDetails
}

func (s FeaturesPrintout) Text() string {
	if len(s.TotalFeaturesTouched) == 0 {
		return "No features found.\n"
	}

	resultString := "Below you can see the list of touched features:\n"
	for file, features := range s.TotalFeaturesTouched {
		if len(features) == 0 {
			continue
		}

		file = strings.ReplaceAll(file, s.AbsolutePath+"/", "")
		resultString += fmt.Sprintf("File: %s\n", file)
		for _, feature := range features {
			resultString += fmt.Sprintf("* %s\n", feature.Name)
		}
	}

	return resultString
}