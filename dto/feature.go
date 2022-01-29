package dto

//Feature the struct for unique feature
type Feature struct {
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
	Line     int `json:"line"`
}
