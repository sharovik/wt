package dto

//IndexedFile the indexed file object struct, which will be used for retrieving of details about the file
type IndexedFile struct {
	Path            string    `json:"path"`
	Features        []Feature `json:"features"`
	RelatedProjects []string `json:"-"`
	MainEntrypoint  string `json:"main_entrypoint"`
	OtherImports    map[string]string `json:"-"`
	UsedIn          []IndexedFile `json:"used_in"`
}

//SyncProjects method for sync of the projects for the indexed file
func (receiver *IndexedFile) SyncProjects(projects []string) {
	for _, p := range projects {
		var isProjectExists = false
		for _, existingProject := range receiver.RelatedProjects {
			if existingProject == p {
				isProjectExists = true
				break
			}
		}

		if isProjectExists {
			continue
		}

		receiver.RelatedProjects = append(receiver.RelatedProjects, p)
	}
}
