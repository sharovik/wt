package printout

import (
	"fmt"
	"strings"
)

//FeaturesPrintout the full printout struct
type FeaturesPrintout struct {
	BasePrintout
}

//Text - main method for text output generation
func (s FeaturesPrintout) Text() string {
	resultString := InfoText(fmt.Sprintf("Analysing the code in path: `%s`\n", s.AbsolutePath))

	if len(s.TotalFeaturesTouched) == 0 {
		resultString += WarningText("No features found.")
		return resultString
	}

	if len(s.ProjectsToCheck) > 0 {
		resultString += WarningText("\nYou might need to implement fixes for the next dependencies:\n")

		for _, p := range s.ProjectsToCheck {
			resultString += NormalText(fmt.Sprintf("* %s\n", p))
		}

		resultString += "\n"
	}

	resultString += InfoText("Below you can see the list of touched features:\n\n")
	alreadyAdded := map[string]string{}
	for _, feature := range s.TotalFeaturesTouched {
		file := strings.ReplaceAll(feature.FilePath, s.AbsolutePath+"/", "")
		for _, feature := range feature.Features {
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
