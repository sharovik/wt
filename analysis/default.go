package analysis

//DefaultAnalysis the default analysis struct
type DefaultAnalysis struct {
}

func (a DefaultAnalysis) GetNeedleKey() string {
	return ""
}

func (a DefaultAnalysis) ExtractMainEntrypointSrc(text string) (res string, err error) {
	return
}

func (a DefaultAnalysis) ExtractImportUsage(text string) (res string, err error) {
	return
}

func (a DefaultAnalysis) FindMainEntrypointTitle(text string) (res string, err error) {
	return
}

func (a DefaultAnalysis) GenerateEntryPointSrc(src string, entrypointTitle string) string {
	return ""
}
