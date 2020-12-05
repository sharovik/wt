package dto

//IndexedFile the indexed file object struct, which will be used for retrieving of details about the file
type IndexedFile struct {
	Path           string            `json:"path"`
	Features       []Feature         `json:"features"`
	MainEntrypoint string            `json:"main_entrypoint"`
	OtherImports   map[string]string `json:"other_imports"`
	UsedIn         []IndexedFile
}
