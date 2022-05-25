package printout

import (
	"fmt"

	"github.com/sharovik/wt/configuration"
	"github.com/sharovik/wt/dto"
	"github.com/sharovik/wt/services"
)

const (
	//DisplayFull - type when we display full list of features
	DisplayFull = "full"
	//DisplayFeatures - type when we display the features
	DisplayFeatures = "features"
	//DisplayJSON - we output the json object
	DisplayJSON = "json"
)

//BasePrintoutInterface the base interface for the printout object
type BasePrintoutInterface interface {
	SetTotalFeaturesTouched(totalFeaturesTouched map[string][]dto.Feature)
	SetAbsolutePath(path string)
	SetConfig(config configuration.Config)
	SetToBeChecked(files map[string]dto.IndexedFile)
	SetProjectsToCheck(projects map[string]string)
	GetToBeChecked() map[string]dto.IndexedFile
	GetProjectsToCheck() map[string]string
	WithToBeCheckedDetails()
	IsToBeCheckedDetailsEnabled() bool
	Text() string
	ToBeCheckedText() string
}

//BasePrintout the base struct for printout
type BasePrintout struct {
	AbsolutePath            string
	Config                  configuration.Config
	PrintToBeCheckedDetails bool
	TotalFeaturesTouched    map[string][]dto.Feature
	ToBeChecked             map[string]dto.IndexedFile
	ProjectsToCheck         map[string]string
}

//PrintObject will be used for printout object generation
type PrintObject struct {
	TotalFeaturesTouched map[string][]dto.Feature   `json:"total_features_touched"`
	ToBeChecked          map[string]dto.IndexedFile `json:"to_be_checked"`
	ProjectsToCheck      map[string]string          `json:"projects_to_check"`
}

//SetAbsolutePath - setter for absolute path
func (s *BasePrintout) SetAbsolutePath(path string) {
	s.AbsolutePath = path
}

//SetConfig - setter for configuration
func (s *BasePrintout) SetConfig(config configuration.Config) {
	s.Config = config
}

//SetTotalFeaturesTouched - setter for total features touched objects map
func (s *BasePrintout) SetTotalFeaturesTouched(features map[string][]dto.Feature) {
	s.TotalFeaturesTouched = features
}

//SetToBeChecked - setter for to be checked files map
func (s *BasePrintout) SetToBeChecked(files map[string]dto.IndexedFile) {
	s.ToBeChecked = files
}

//GetToBeChecked - getter for to be checked map
func (s BasePrintout) GetToBeChecked() map[string]dto.IndexedFile {
	return s.ToBeChecked
}

//WithToBeCheckedDetails - sets the print to be checked flag
func (s *BasePrintout) WithToBeCheckedDetails() {
	s.PrintToBeCheckedDetails = true
}

//IsToBeCheckedDetailsEnabled - return the current state of PrintToBeCheckedDetails flag
func (s BasePrintout) IsToBeCheckedDetailsEnabled() bool {
	return s.PrintToBeCheckedDetails
}

//ToBeCheckedText - generates the "to be checked" text
func (s BasePrintout) ToBeCheckedText() string {
	return generateToBeCheckedText(&s)
}

//SetProjectsToCheck - setter for ProjectsToCheck attribute
func (s *BasePrintout) SetProjectsToCheck(projects map[string]string) {
	s.ProjectsToCheck = projects
}

//GetProjectsToCheck - getter for ProjectsToCheck attribute
func (s BasePrintout) GetProjectsToCheck() map[string]string {
	return s.ProjectsToCheck
}

//Text - main method for text output generation
func (s BasePrintout) Text() string {
	return ""
}

func generateToBeCheckedText(obj BasePrintoutInterface) string {
	if len(obj.GetToBeChecked()) == 0 {
		return ""
	}

	resultString := fmt.Sprintf("\n\nYour changes can potentially touch the functionality in the `%d` files.", len(obj.GetToBeChecked()))
	if obj.IsToBeCheckedDetailsEnabled() {
		resultString += fmt.Sprintf("\nPlease check the following files:\n\n")
		resultString += fmt.Sprintf("%s\n\n", generateToBeCheckedDetails(obj.GetToBeChecked()))
	} else {
		resultString += fmt.Sprintf("\nThese files does not have `%s` annotation.\nRun comman with `-withToBeChecked=true` flag for more details.\n", services.FeatureAlias)
	}

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

//FromType generates the printout object based on selected display type
func FromType(displayType string) (BasePrintoutInterface, error) {
	switch displayType {
	case DisplayFull:
		return &FullPrintout{}, nil
	case DisplayFeatures:
		return &FeaturesPrintout{}, nil
	case DisplayJSON:
		return &JSONPrintout{}, nil
	}

	return nil, fmt.Errorf("Failed to generate printout obj. ")
}

//InfoText - prints the text with INFO color
func InfoText(text string) string {
	return fmt.Sprintf("\033[1;32m%s\033[0m", text)
}

//NormalText - prints the text with Normal color
func NormalText(text string) string {
	return fmt.Sprintf("\033[1;36m%s\033[0m", text)
}

//WarningText - prints the text with Warning color
func WarningText(text string) string {
	return fmt.Sprintf("\033[1;33m%s\033[0m", text)
}
