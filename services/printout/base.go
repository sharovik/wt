package printout

import (
	"fmt"
	"github.com/sharovik/wt/dto"
	"github.com/sharovik/wt/services"
)

const (
	DisplayFull     = "full"
	DisplayFeatures = "features"
)

type BasePrintoutInterface interface {
	SetTotalFeaturesTouched(totalFeaturesTouched map[string][]dto.Feature)
	SetAbsolutePath(path string)
	SetToBeChecked(files map[string]dto.IndexedFile)
	GetToBeChecked() map[string]dto.IndexedFile
	WithToBeCheckedDetails()
	IsToBeCheckedDetailsEnabled() bool
	Text() string
	ToBeCheckedText() string
}

func generateToBeCheckedText(obj BasePrintoutInterface) string {
	if len(obj.GetToBeChecked()) == 0 {
		return ""
	}

	resultString := fmt.Sprintf("Your changes can potentially touch the functionality in the `%d` files.", len(obj.GetToBeChecked()))
	if obj.IsToBeCheckedDetailsEnabled() {
		resultString += fmt.Sprintf("\nPlease check the following files:\n\n")
		resultString += fmt.Sprintf("%s\n\n", generateToBeCheckedDetails(obj.GetToBeChecked()))
	} else {
		resultString += fmt.Sprintf("\nThese files does not have `%s` annotation.\nRun comman with `-withToBeChecked=true` flag for more details.", services.FeatureAlias)
	}

	resultString += "\n\n"

	return resultString
}

func generateToBeCheckedDetails(toBeChecked map[string]dto.IndexedFile) (resultString string) {
	if len(toBeChecked) == 0 {
		return ""
	}

	for relativePath, file := range toBeChecked {
		resultString += fmt.Sprintf("- `%s`", relativePath)
		if len(file.UsedIn) > 0 {
			resultString += " touched by ["
			for _, usedInFile := range file.UsedIn {
				resultString += fmt.Sprintf("%s,", usedInFile.MainEntrypoint)
			}
			resultString += "]"
		}

		resultString += "\n"
	}

	resultString += fmt.Sprintf("\n !!!Warning!!! Please make sure you add `%s` annotation to these files.", services.FeatureAlias)
	return
}

func GeneratePrintoutObj(displayType string, totalFeaturesTouched map[string][]dto.Feature, path string, toBeChecked map[string]dto.IndexedFile) (BasePrintoutInterface, error) {
	switch displayType {
	case DisplayFull:
		obj := new(FullPrintout)
		obj.SetAbsolutePath(path)
		obj.SetToBeChecked(toBeChecked)
		obj.SetTotalFeaturesTouched(totalFeaturesTouched)

		return obj, nil
	case DisplayFeatures:
		obj := new(FeaturesPrintout)
		obj.SetAbsolutePath(path)
		obj.SetToBeChecked(toBeChecked)
		obj.SetTotalFeaturesTouched(totalFeaturesTouched)

		return obj, nil
	}

	return nil, fmt.Errorf("Failed to generate printout obj. ")
}