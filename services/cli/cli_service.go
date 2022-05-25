package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/configuration"
	"github.com/sharovik/wt/services/printout"
)

const (
	defaultIgnorePath        = ".gitignore"
	defaultDestinationBranch = "master"
	defaultIgnoredPaths      = "tests"

	appVersion = "v1.2.0"

	versionTemplate = "What touched by sharovik. Version: %s\n\n"
)

//Service the service struct
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
	CPUProfile         string
	MemProfile         string
}

//ParseArgs we parse the arguments from the CLI
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
	s.CPUProfile = *cpuProfile
	s.MemProfile = *memProfile

	s.loadDefaults()
	s.initConfiguration()
}

func (s Service) initConfiguration() {
	configuration.C = configuration.Config{
		WorkingBranch:      s.WorkingBranch,
		DestinationBranch:  s.DestinationBranch,
		VcsType:            s.VcsType,
		Path:               s.Path,
		Ext:                s.Ext,
		PathToIgnoreFile:   s.PathToIgnoreFile,
		DisplayTemplate:    s.DisplayTemplate,
		IgnoreFromAnalysis: s.IgnoreFromAnalysis,
		MaxAnalysisDepth:   s.MaxAnalysisDepth,
		WithToBeChecked:    s.WithToBeChecked,
		Version:            s.Version,
		CpuProfile:         s.CPUProfile,
		MemProfile:         s.MemProfile,
	}
}

func (s *Service) loadDefaults() {
	if s.Version {
		fmt.Println(fmt.Sprintf(versionTemplate, appVersion))
		os.Exit(0)
	}

	if s.VcsType == "" {
		log.Fatal(fmt.Errorf("The vcs should not be empty "))
	}

	if s.Path == "" {
		s.Path = "."
	}

	if s.PathToIgnoreFile == "" {
		s.PathToIgnoreFile = defaultIgnorePath
	}

	if s.DisplayTemplate == "" {
		s.DisplayTemplate = printout.DisplayFeatures
	}

	if s.WorkingBranch == "" || s.DestinationBranch == "" {
		log.Fatal(fmt.Errorf("Working branch and destination branch should not be empty. Please make sure you define them. "))
	}
}
