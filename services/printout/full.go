package printout

import (
	"fmt"
)

//FullPrintout the full printout struct
type FullPrintout struct {
	BasePrintout
}

//Text - main method for text output generation
func (s FullPrintout) Text() string {
	resultString := InfoText(fmt.Sprintf("Analysing the code in path: `%s`\n", s.AbsolutePath))

	if len(s.TotalFeaturesTouched) == 0 {
		resultString += WarningText("No features found.")
		return resultString
	}

	if len(s.ProjectsToCheck) > 0 {
		resultString += WarningText("You might need to implement fixes for the next dependencies:\n")

		for _, p := range s.ProjectsToCheck {
			resultString += NormalText(fmt.Sprintf("* %s\n", p))
		}

		resultString += "\n\n"
	}

	resultString += InfoText("Below you can see the list of touched features:\n\n")
	for _, touchedFeature := range s.TotalFeaturesTouched {
		for _, feature := range touchedFeature.Features {
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
