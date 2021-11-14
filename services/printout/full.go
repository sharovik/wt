package printout

import (
	"fmt"
	"github.com/sharovik/wt/dto"
	"strings"
)

//FullPrintout the full printout struct
type FullPrintout struct {
	AbsolutePath string
	PrintToBeCheckedDetails bool
	TotalFeaturesTouched map[string][]dto.Feature
	ToBeChecked map[string]dto.IndexedFile
}

func (s FullPrintout) SetAbsolutePath(path string) {
	s.AbsolutePath = path
}

func (s FullPrintout) SetTotalFeaturesTouched(features map[string][]dto.Feature) {
	s.TotalFeaturesTouched = features
}

func (s FullPrintout) SetToBeChecked(files map[string]dto.IndexedFile) {
	s.ToBeChecked = files
}

func (s FullPrintout) GetToBeChecked() map[string]dto.IndexedFile {
	return s.ToBeChecked
}

func (s FullPrintout) WithToBeCheckedDetails() {
	s.PrintToBeCheckedDetails = true
}

func (s FullPrintout) IsToBeCheckedDetailsEnabled() bool {
	return s.PrintToBeCheckedDetails
}

func (s FullPrintout) ToBeCheckedText() string {
	return generateToBeCheckedText(s)
}

func (s FullPrintout) Text() string {
	if len(s.TotalFeaturesTouched) == 0 {
		return ""
	}

	resultString := "Below you can see the list of touched features:\n\n"
	for file, features := range s.TotalFeaturesTouched {
		if len(features) == 0 {
			continue
		}

		file = strings.ReplaceAll(file, s.AbsolutePath+"/", "")
		for _, feature := range features {
			resultString += "------------------\n"
			resultString += fmt.Sprintf("Feature: %s\n", feature.Name)
			resultString += fmt.Sprintf("Code path: %s:%d\n", feature.FilePath, feature.Line)
		}
		resultString += "------------------\n"
	}

	return resultString
}