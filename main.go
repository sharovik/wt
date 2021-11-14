package main

import (
	"fmt"
	"github.com/sharovik/wt/services/printout"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/dto"

	"github.com/sharovik/wt/services"
	"github.com/sharovik/wt/services/cli"
)

const (
	defaultIgnorePath        = ".gitignore"
	defaultDestinationBranch = "master"
	defaultIgnoredPaths      = "tests"

	displayFull     = "full"
	displayFeatures = "features"

	appVersion = "v1.1.4"

	versionTemplate = "What touched by sharovik. Version: %s\n\n"
)

var (
	vcs        services.VcsInterface
	cliService = cli.Service{}
)

func main() {
	cliService.ParseArgs()

	if cliService.CpuProfile != "" {
		f, err := os.Create(cliService.CpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if cliService.Version {
		fmt.Println(fmt.Sprintf(versionTemplate, appVersion))
		return
	}

	if cliService.VcsType == "" {
		log.Fatal(fmt.Errorf("The vcs should not be empty "))
	}

	if cliService.Path == "" {
		cliService.Path = "."
	}

	if cliService.PathToIgnoreFile == "" {
		cliService.PathToIgnoreFile = defaultIgnorePath
	}

	if cliService.DisplayTemplate == "" {
		cliService.DisplayTemplate = displayFeatures
	}

	if cliService.WorkingBranch == "" || cliService.DestinationBranch == "" {
		log.Fatal(fmt.Errorf("Working branch and destination branch should not be empty. Please make sure you define them. "))
	}

	analysis.MaxDeepLevel = cliService.MaxAnalysisDepth

	loadVcs(cliService.VcsType)
	analysis.InitAnalysisService(cliService.Ext)

	absolutePath, err := filepath.Abs(cliService.Path)
	if err != nil {
		return
	}

	paths, err := services.GetIgnoredFilePaths(cliService.PathToIgnoreFile, absolutePath)
	if err != nil {
		log.Fatal(err)
	}

	if cliService.IgnoreFromAnalysis != "" {
		for _, path := range strings.Split(cliService.IgnoreFromAnalysis, ",") {
			paths = append(paths, fmt.Sprintf("%s/%s", absolutePath, path))
		}
	}

	fmt.Println(fmt.Sprintf("Start analysing the code in path: `%s`", absolutePath))
	index, pathIndex, importsIndex, err := services.AnalyseTheCode(absolutePath, cliService.Ext, paths)
	if err != nil {
		log.Fatal(err)
	}

	analysis.AnalysedEntrypointsIndex = index
	analysis.AnalysedPathsIndex = pathIndex
	analysis.AnalysedImportsIndex = importsIndex

	diff, err := vcs.Diff(absolutePath, cliService.WorkingBranch, cliService.DestinationBranch)
	if err != nil {
		log.Fatal(err)
	}

	totalFeaturesTouched, toBeChecked := services.FindFeaturesInIndex(diff, absolutePath)

	resultString := fmt.Sprintf(versionTemplate, appVersion)
	if len(diff) == 0 {
		displayObj, err := printout.GeneratePrintoutObj(cliService.DisplayTemplate, totalFeaturesTouched, absolutePath, toBeChecked)
		if err != nil {
			log.Fatal(err)
		}

		if len(toBeChecked) > 0 {
			resultString += fmt.Sprintf("Your changes can potentially touch the functionality in the `%d` files.", len(toBeChecked))
			if cliService.WithToBeChecked {
				resultString += fmt.Sprintf("\nPlease check the following files:\n\n")
				resultString += fmt.Sprintf("%s\n\n", printToBeChecked(toBeChecked))
			} else {
				resultString += fmt.Sprintf("\nThese files does not have `%s` annotation.\nRun comman with `-withToBeChecked=true` flag for more details.", services.FeatureAlias)
			}

			resultString += "\n\n"
		}

		res := displayObj.Text()

		if res == "" {
			resultString += "I found no features touched by these changes.\nPlease, make sure you define the features in these diff"
		} else {
			resultString += fmt.Sprintf("%s\nPlease make sure you test these features before you merge `%s` branch into `%s`.", res, cliService.WorkingBranch, cliService.DestinationBranch)
		}
	} else {
		resultString += fmt.Sprintf("\nThere is no diff between `%s` and `%s`", cliService.WorkingBranch, cliService.DestinationBranch)
	}

	fmt.Println(resultString)

	if cliService.MemProfile != "" {
		f, err := os.Create(cliService.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func printToBeChecked(toBeChecked map[string]dto.IndexedFile) (resultString string) {
	if len(toBeChecked) == 0 {
		return ""
	}

	for relativePath, file := range toBeChecked {
		resultString += fmt.Sprintf("- `%s`", relativePath)
		if len(file.UsedIn) > 0 {
			resultString += " touched by ["
			for _, usedInFile := range file.UsedIn {
				resultString += fmt.Sprintf("%s,", usedInFile.MainEntrypoint)
			}
			resultString += "]"
		}

		resultString += "\n"
	}

	resultString += fmt.Sprintf("\n !!!Warning!!! Please make sure you add `%s` annotation to these files.", services.FeatureAlias)
	return
}

func loadVcs(vcsType string) {
	switch vcsType {
	case "git":
		vcs = services.Git{}
		break
	}
}
