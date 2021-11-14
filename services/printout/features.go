package printout

import (
	"fmt"
	"github.com/sharovik/wt/configuration"
	"github.com/sharovik/wt/dto"
	"strings"
)

//FeaturesPrintout the full printout struct
type FeaturesPrintout struct {
	AbsolutePath string
	Config configuration.Config
	PrintToBeCheckedDetails bool
	TotalFeaturesTouched map[string][]dto.Feature
	ToBeChecked map[string]dto.IndexedFile
}

func (s *FeaturesPrintout) SetAbsolutePath(path string) {
	s.AbsolutePath = path
}

func (s *FeaturesPrintout) SetConfig(config configuration.Config) {
	s.Config = config
}

func (s *FeaturesPrintout) SetTotalFeaturesTouched(features map[string][]dto.Feature) {
	s.TotalFeaturesTouched = features
}

func (s *FeaturesPrintout) SetToBeChecked(files map[string]dto.IndexedFile) {
	s.ToBeChecked = files
}

func (s FeaturesPrintout) GetToBeChecked() map[string]dto.IndexedFile {
	return s.ToBeChecked
}

func (s FeaturesPrintout) ToBeCheckedText() string {
	return generateToBeCheckedText(&s)
}

func (s *FeaturesPrintout) WithToBeCheckedDetails() {
	s.PrintToBeCheckedDetails = true
}

func (s FeaturesPrintout) IsToBeCheckedDetailsEnabled() bool {
	return s.PrintToBeCheckedDetails
}

func (s FeaturesPrintout) Text() string {
	resultString := InfoText(fmt.Sprintf("Analysing the code in path: `%s`\n", s.AbsolutePath))

	if len(s.TotalFeaturesTouched) == 0 {
		resultString += WarningText("No features found.")
		return resultString
	}

	resultString += InfoText("Below you can see the list of touched features:\n\n")
	alreadyAdded := map[string]string{}
	for file, features := range s.TotalFeaturesTouched {
		if len(features) == 0 {
			continue
		}

		file = strings.ReplaceAll(file, s.AbsolutePath+"/", "")
		for _, feature := range features {
			if alreadyAdded[feature.Name] != "" {
				continue
			}

			resultString += NormalText(fmt.Sprintf("* %s (file: %s)\n", feature.Name, file))
			alreadyAdded[feature.Name] = feature.Name
		}
	}

	resultString += InfoText(fmt.Sprintf("%s", s.ToBeCheckedText()))
	resultString += WarningText(fmt.Sprintf("\nPlease make sure you test these features before you merge `%s` branch into `%s`.", s.Config.WorkingBranch, s.Config.DestinationBranch))

	return resultString
}