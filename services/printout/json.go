package printout

import (
	"encoding/json"
	"fmt"
	"github.com/sharovik/wt/configuration"
	"github.com/sharovik/wt/dto"
)

//JsonPrintout the full printout struct
type JsonPrintout struct {
	AbsolutePath            string
	Config                  configuration.Config
	PrintToBeCheckedDetails bool
	TotalFeaturesTouched    map[string][]dto.Feature
	ToBeChecked             map[string]dto.IndexedFile
}

//PrintObject will be used for printout object generation
type PrintObject struct {
	TotalFeaturesTouched    map[string][]dto.Feature `json:"total_features_touched"`
	ToBeChecked             map[string]dto.IndexedFile `json:"to_be_checked"`
}

func (s *JsonPrintout) SetAbsolutePath(path string) {
	s.AbsolutePath = path
}

func (s *JsonPrintout) SetConfig(config configuration.Config) {
	s.Config = config
}

func (s *JsonPrintout) SetTotalFeaturesTouched(features map[string][]dto.Feature) {
	s.TotalFeaturesTouched = features
}

func (s *JsonPrintout) SetToBeChecked(files map[string]dto.IndexedFile) {
	s.ToBeChecked = files
}

func (s JsonPrintout) GetToBeChecked() map[string]dto.IndexedFile {
	return s.ToBeChecked
}

func (s *JsonPrintout) WithToBeCheckedDetails() {
	s.PrintToBeCheckedDetails = true
}

func (s JsonPrintout) IsToBeCheckedDetailsEnabled() bool {
	return s.PrintToBeCheckedDetails
}

func (s JsonPrintout) ToBeCheckedText() string {
	return generateToBeCheckedText(&s)
}

func (s JsonPrintout) Text() string {
	obj := PrintObject{
		TotalFeaturesTouched: s.TotalFeaturesTouched,
		ToBeChecked:          s.ToBeChecked,
	}

	outputStr := ""
	byteStr, err := json.Marshal(obj)
	if err != nil {
		outputStr += fmt.Sprintf("Error: %s", err)
	}

	outputStr = string(byteStr)
	return outputStr
}
