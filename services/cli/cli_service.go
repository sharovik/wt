package cli

import (
	"flag"
	"fmt"
	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/services/printout"
)

const (
	defaultIgnorePath        = ".gitignore"
	defaultDestinationBranch = "master"
	defaultIgnoredPaths      = "tests"
)

type Service struct {
	WorkingBranch      string
	DestinationBranch  string
	VcsType            string
	Path               string
	Ext                string
	PathToIgnoreFile   string
	DisplayTemplate    string
	IgnoreFromAnalysis string
	MaxAnalysisDepth   int
	WithToBeChecked    bool
	Version            bool
	CpuProfile         string
	MemProfile         string
}

func (s *Service) ParseArgs() {
	workingBranch := flag.String("workingBranch", "", "Working branch which will be compared with the destination branch.")
	destinationBranch := flag.String("destinationBranch", defaultDestinationBranch, "Destination branch with which we will compare selected working branch.")
	vcsType := flag.String("vcs", "git", "The type of vcs which will be used for retrieving diff information.")
	path := flag.String("path", ".", "The type of vcs which will be used for retrieving diff information.")
	ext := flag.String("fileExt", "", "The type of extension of the diff which we need to check.")
	pathToIgnoreFile := flag.String("pathToIgnoreFile", defaultIgnorePath, fmt.Sprintf("The path to file, where line-by-line written the list of paths which should be ignored. Default it's: %s", defaultIgnorePath))
	displayTemplate := flag.String("displayTemplate", printout.DisplayFeatures, fmt.Sprintf("The view which will be used for display results. Default is: %s", printout.DisplayFeatures))
	ignoreFromAnalysis := flag.String("ignoreFromAnalysis", defaultIgnoredPaths, fmt.Sprintf("The list of folders/files separated by comma, which will be ignored during the files analysis. Default is: %s", defaultIgnoredPaths))
	maxAnalysisDepth := flag.Int("maxAnalysisDepth", analysis.DefaultMaxDeepLevel, fmt.Sprintf("The maximum analysis code depth will be used during the code usage analysing. Default is: %d", analysis.DefaultMaxDeepLevel))
	withToBeChecked := flag.Bool("withToBeChecked", false, fmt.Sprintf("Display or not the files which should be covered by features annotation. Default is: %v", false))
	version := flag.Bool("version", false, "Shows the app version.")
	cpuProfile := flag.String("cpuProfile", "", "write cpu profile to `file`")
	memProfile := flag.String("memProfile", "", "write memory profile to `file`")

	flag.Parse()

	s.WorkingBranch = *workingBranch
	s.DestinationBranch = *destinationBranch
	s.VcsType = *vcsType
	s.Path = *path
	s.Ext = *ext
	s.PathToIgnoreFile = *pathToIgnoreFile
	s.DisplayTemplate = *displayTemplate
	s.IgnoreFromAnalysis = *ignoreFromAnalysis
	s.MaxAnalysisDepth = *maxAnalysisDepth
	s.WithToBeChecked = *withToBeChecked
	s.Version = *version
	s.CpuProfile = *cpuProfile
	s.MemProfile = *memProfile
}
