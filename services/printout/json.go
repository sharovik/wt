package printout

import (
	"encoding/json"
	"fmt"
)

//JSONPrintout the full printout struct
type JSONPrintout struct {
	BasePrintout
}

//Text - generates the text output
func (s JSONPrintout) Text() string {
	obj := PrintObject{
		TotalFeaturesTouched: s.TotalFeaturesTouched,
		ToBeChecked:          s.ToBeChecked,
		ProjectsToCheck:      s.ProjectsToCheck,
	}

	byteStr, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	return string(byteStr)
}
