package printout

import (
	"fmt"
	"strings"

	"github.com/sharovik/wt/configuration"
	"github.com/sharovik/wt/dto"
)

//FullPrintout the full printout struct
type FullPrintout struct {
	AbsolutePath            string
	Config                  configuration.Config
	PrintToBeCheckedDetails bool
	TotalFeaturesTouched    map[string][]dto.Feature
	ToBeChecked             map[string]dto.IndexedFile
}

func (s *FullPrintout) SetAbsolutePath(path string) {
	s.AbsolutePath = path
}

func (s *FullPrintout) SetConfig(config configuration.Config) {
	s.Config = config
}

func (s *FullPrintout) SetTotalFeaturesTouched(features map[string][]dto.Feature) {
	s.TotalFeaturesTouched = features
}

func (s *FullPrintout) SetToBeChecked(files map[string]dto.IndexedFile) {
	s.ToBeChecked = files
}

func (s FullPrintout) GetToBeChecked() map[string]dto.IndexedFile {
	return s.ToBeChecked
}

func (s *FullPrintout) WithToBeCheckedDetails() {
	s.PrintToBeCheckedDetails = true
}

func (s FullPrintout) IsToBeCheckedDetailsEnabled() bool {
	return s.PrintToBeCheckedDetails
}

func (s FullPrintout) ToBeCheckedText() string {
	return generateToBeCheckedText(&s)
}

func (s FullPrintout) Text() string {
	resultString := InfoText(fmt.Sprintf("Analysing the code in path: `%s`\n", s.AbsolutePath))

	if len(s.TotalFeaturesTouched) == 0 {
		resultString += WarningText("No features found.")
		return resultString
	}

	resultString += InfoText("Below you can see the list of touched features:\n\n")
	for file, features := range s.TotalFeaturesTouched {
		if len(features) == 0 {
			continue
		}

		file = strings.ReplaceAll(file, s.AbsolutePath+"/", "")
		for _, feature := range features {
			resultString += NormalText("------------------\n")
			resultString += NormalText(fmt.Sprintf("Feature: %s\n", feature.Name))
			resultString += NormalText(fmt.Sprintf("Code path: %s:%d\n", feature.FilePath, feature.Line))
		}
		resultString += NormalText("------------------\n")
	}

	resultString += InfoText(fmt.Sprintf("%s", s.ToBeCheckedText()))
	resultString += WarningText(fmt.Sprintf("\nPlease make sure you test these features before you merge `%s` branch into `%s`.", s.Config.WorkingBranch, s.Config.DestinationBranch))

	return resultString
}
