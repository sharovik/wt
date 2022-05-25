package dto

//AnalysisResult the analysis result struct will be used for final analysis output generation
type AnalysisResult struct {
	ToBeChecked map[string]IndexedFile
	TotalFeaturesTouched map[string][]Feature
	ProjectsToCheck map[string]string
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