package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/dto"

	"github.com/sharovik/wt/services"
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

var vcs services.VcsInterface

func main() {
	workingBranch := flag.String("workingBranch", "", "Working branch which will be compared with the destination branch.")
	destinationBranch := flag.String("destinationBranch", defaultDestinationBranch, "Destination branch with which we will compare selected working branch.")
	vcsType := flag.String("vcs", "git", "The type of vcs which will be used for retrieving diff information.")
	path := flag.String("path", ".", "The type of vcs which will be used for retrieving diff information.")
	ext := flag.String("fileExt", "", "The type of extension of the diff which we need to check.")
	pathToIgnoreFile := flag.String("pathToIgnoreFile", defaultIgnorePath, fmt.Sprintf("The path to file, where line-by-line written the list of paths which should be ignored. Default it's: %s", defaultIgnorePath))
	displayTemplate := flag.String("displayTemplate", displayFeatures, fmt.Sprintf("The view which will be used for display results. Default is: %s", displayFeatures))
	ignoreFromAnalysis := flag.String("ignoreFromAnalysis", defaultIgnoredPaths, fmt.Sprintf("The list of folders/files separated by comma, which will be ignored during the files analysis. Default is: %s", defaultIgnoredPaths))
	maxAnalysisDepth := flag.Int("maxAnalysisDepth", analysis.DefaultMaxDeepLevel, fmt.Sprintf("The maximum analysis code depth will be used during the code usage analysing. Default is: %d", analysis.DefaultMaxDeepLevel))
	withToBeChecked := flag.Bool("withToBeChecked", false, fmt.Sprintf("Display or not the files which should be covered by features annotation. Default is: %v", false))
	version := flag.Bool("version", false, "Shows the app version.")
	cpuProfile := flag.String("cpuProfile", "", "write cpu profile to `file`")
	memProfile := flag.String("memProfile", "", "write memory profile to `file`")

	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *version {
		fmt.Println(fmt.Sprintf(versionTemplate, appVersion))
		return
	}

	if *vcsType == "" {
		log.Fatal(fmt.Errorf("The vcs should not be empty "))
	}

	if *path == "" {
		*path = "."
	}

	if *pathToIgnoreFile == "" {
		*pathToIgnoreFile = defaultIgnorePath
	}

	if *displayTemplate == "" {
		*displayTemplate = displayFeatures
	}

	if *workingBranch == "" || *destinationBranch == "" {
		log.Fatal(fmt.Errorf("Working branch and destination branch should not be empty. Please make sure you define them. "))
	}

	analysis.MaxDeepLevel = *maxAnalysisDepth

	loadVcs(*vcsType)
	analysis.InitAnalysisService(*ext)

	absolutePath, err := filepath.Abs(*path)
	if err != nil {
		return
	}

	paths, err := services.GetIgnoredFilePaths(*pathToIgnoreFile, absolutePath)
	if err != nil {
		log.Fatal(err)
	}

	if *ignoreFromAnalysis != "" {
		for _, path := range strings.Split(*ignoreFromAnalysis, ",") {
			paths = append(paths, fmt.Sprintf("%s/%s", absolutePath, path))
		}
	}

	fmt.Println(fmt.Sprintf("Start analysing the code in path: `%s`", absolutePath))
	index, pathIndex, importsIndex, err := services.AnalyseTheCode(absolutePath, *ext, paths)
	if err != nil {
		log.Fatal(err)
	}

	analysis.AnalysedEntrypointsIndex = index
	analysis.AnalysedPathsIndex = pathIndex
	analysis.AnalysedImportsIndex = importsIndex

	diff, err := vcs.Diff(absolutePath, *workingBranch, *destinationBranch)
	if err != nil {
		log.Fatal(err)
	}

	totalFeaturesTouched, toBeChecked := services.FindFeaturesInIndex(diff, absolutePath)

	resultString := fmt.Sprintf(versionTemplate, appVersion)
	if len(diff) > 0 {
		if len(toBeChecked) > 0 {
			resultString += fmt.Sprintf("Your changes can potentially touch the functionality in the `%d` files.", len(toBeChecked))
			if *withToBeChecked {
				resultString += fmt.Sprintf("\nPlease check the following files:\n\n")
				resultString += fmt.Sprintf("%s\n\n", printToBeChecked(toBeChecked))
			} else {
				resultString += fmt.Sprintf("\nThese files does not have `%s` annotation.\nRun comman with `-withToBeChecked=true` flag for more details.", services.FeatureAlias)
			}

			resultString += "\n"
		}

		res := ""
		switch *displayTemplate {
		case displayFull:
			res = printFull(totalFeaturesTouched, absolutePath)
			break
		case displayFeatures:
			res = printFeatures(totalFeaturesTouched, absolutePath)
			break
		}

		if res == "" {
			resultString += "I found no features touched by these changes.\nPlease, make sure you define the features in these diff"
		} else {
			resultString += fmt.Sprintf("%s\nPlease make sure you test these features before you merge `%s` branch into `%s`.", res, *workingBranch, *destinationBranch)
		}
	} else {
		resultString += fmt.Sprintf("\nThere is no diff between `%s` and `%s`", *workingBranch, *destinationBranch)
	}

	fmt.Println(resultString)

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
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

func printFull(files map[string][]dto.Feature, absolutePath string) string {
	if len(files) == 0 {
		return ""
	}

	resultString := "Below you can see the list of touched features:\n"
	for file, features := range files {
		if len(features) == 0 {
			continue
		}

		file = strings.ReplaceAll(file, absolutePath+"/", "")
		for _, feature := range features {
			resultString += "------------------\n"
			resultString += fmt.Sprintf("Feature: %s\n", feature.Name)
			resultString += fmt.Sprintf("Code path: %s:%d\n", feature.FilePath, feature.Line)
		}
		resultString += "------------------\n"
	}

	return resultString
}

func printFeatures(files map[string][]dto.Feature, absolutePath string) string {
	if len(files) == 0 {
		return "No features found.\n"
	}

	resultString := "Below you can see the list of touched features:\n"
	for file, features := range files {
		if len(features) == 0 {
			continue
		}

		file = strings.ReplaceAll(file, absolutePath+"/", "")
		resultString += fmt.Sprintf("File: %s\n", file)
		for _, feature := range features {
			resultString += fmt.Sprintf("* %s\n", feature.Name)
		}
	}

	return resultString
}

func loadVcs(vcsType string) {
	switch vcsType {
	case "git":
		vcs = services.Git{}
		break
	}
}
