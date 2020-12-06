package analysis

import "github.com/sharovik/wt/dto"

//BaseAnalysisInterface interface which will be used for objects which will help to find the features and namespaces
type BaseAnalysisInterface interface {
	GetNeedleKey() string

	//ExtractMainEntrypointSrc method retrieves the entrypoint src from selected text string. This string will be used for main entrypoint path generation.
	ExtractMainEntrypointSrc(text string) (string, error)

	//ExtractImportUsage retrieves the import-src from selected text string. This will be used in the functionality where we get information about usage of the file.
	ExtractImportUsage(text string) (string, error)

	//FindMainEntrypointTitle method returns the main entrypoint title string from selected text string
	FindMainEntrypointTitle(text string) (string, error)

	//GenerateEntryPointSrc generates the entrypoint src string based on received src and entrypoint title
	GenerateEntryPointSrc(src string, entrypointTitle string) string
}

const DefaultMaxDeepLevel = 2

var (
	//AnalysedEntrypointsIndex this index will be used for finding of the information about the main entrypoint by the entrypoint name
	AnalysedEntrypointsIndex = map[string]dto.IndexedFile{}

	//AnalysedPathsIndex this index will be used for finding of the information about the main entrypoint by the entrypoint path
	AnalysedPathsIndex = map[string]dto.IndexedFile{}

	//AnalysedImportsIndex will be used for indexing of other imports used in your object
	AnalysedImportsIndex = map[string][]string{}
	MaxDeepLevel         = DefaultMaxDeepLevel
	An                   BaseAnalysisInterface
)

//FindUsage will retrieve the list of file paths where your entrypoint is used
func FindUsage(entrypoint string, usage []string, deepLevel int) []string {
	if len(AnalysedImportsIndex[entrypoint]) == 0 {
		return usage
	}

	if exists(usage, entrypoint) {
		return usage
	}

	usage = append(usage, AnalysedImportsIndex[entrypoint]...)
	if deepLevel <= MaxDeepLevel {
		deepLevel += 1
		for _, otherImport := range AnalysedImportsIndex[entrypoint] {
			usage = FindUsage(otherImport, usage, deepLevel)
		}
	}

	return usage
}

func exists(a []string, n string) (exists bool) {
	for _, p := range a {
		if p == n {
			return true
		}
	}

	return exists
}

func InitAnalysisService(ext string) {
	switch ext {
	case ".php":
		An = PhpAnalysis{}
		break
	default:
		An = DefaultAnalysis{}
		break
	}
}
