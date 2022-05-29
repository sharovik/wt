package dto

//AnalysisResult the analysis result struct will be used for final analysis output generation
type AnalysisResult struct {
	ToBeChecked          []ToCheck
	TotalFeaturesTouched []FeatureTouched
	ProjectsToCheck      []string
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
func (r *AnalysisResult) AppendProjects(projects []string) {
	for _, p := range projects {
		var isProjectExists = false
		for _, existingProject := range r.ProjectsToCheck {
			if existingProject == p {
				isProjectExists = true
				break
			}
		}

		if isProjectExists {
			continue
		}

		r.ProjectsToCheck = append(r.ProjectsToCheck, p)
	}
}
