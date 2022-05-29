package dto

//AnalysisResult the analysis result struct will be used for final analysis output generation
type AnalysisResult struct {
	ToBeChecked          []ToCheck
	TotalFeaturesTouched []FeatureTouched
	ProjectsToCheck      map[string]string
}

//FeatureTouched the touched feature struct
type FeatureTouched struct {
	FilePath string
	Features []Feature
}

//ToCheck the file to check
type ToCheck struct {
	FilePath string
	IndexedFile IndexedFile
}

//AppendProjects method for sync of the projects for the indexed file
func (r *AnalysisResult) AppendProjects(projects map[string]string) {
	if nil == r.ProjectsToCheck {
		r.ProjectsToCheck = map[string]string{}
	}

	if len(projects) == 0 {
		return
	}

	for _, p := range projects {
		if r.ProjectsToCheck[p] != "" {
			continue
		}

		r.ProjectsToCheck[p] = p
	}
}
