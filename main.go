package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/sharovik/wt/configuration"
	"github.com/sharovik/wt/services/printout"

	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/services"
	"github.com/sharovik/wt/services/cli"
)

var (
	vcs        services.VcsInterface
	cliService = cli.Service{}
)

func main() {
	cliService.ParseArgs()

	if cliService.CPUProfile != "" {
		f, err := os.Create(cliService.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
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

	analysisResult := services.FindFeaturesInIndex(diff, absolutePath)

	displayObj, err := printout.FromType(cliService.DisplayTemplate)
	if err != nil {
		log.Fatal(err)
	}

	displayObj.SetToBeChecked(analysisResult.ToBeChecked)
	displayObj.SetTotalFeaturesTouched(analysisResult.TotalFeaturesTouched)
	displayObj.SetAbsolutePath(absolutePath)
	displayObj.SetConfig(configuration.C)
	displayObj.SetProjectsToCheck(analysisResult.ProjectsToCheck)

	if cliService.WithToBeChecked {
		displayObj.WithToBeCheckedDetails()
	}

	if len(analysisResult.TotalFeaturesTouched) == 0 {
		fmt.Println(displayObj.Text())

		if cliService.MemProfile != "" {
			f, err := os.Create(cliService.MemProfile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			defer f.Close() // error handling omitted for example
			runtime.GC()    // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}

		return
	}

	fmt.Println(displayObj.Text())

	if cliService.MemProfile != "" {
		f, err := os.Create(cliService.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func loadVcs(vcsType string) {
	switch vcsType {
	case "git":
		vcs = services.Git{}
		break
	}
}
