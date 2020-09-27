package main

import (
	"flag"
	"fmt"
	"github.com/sharovik/wt/services"
	"log"
)

const defaultIgnorePath = ".gitignore"
const defaultDestinationBranch = "master"

var vcs services.VcsInterface

func main() {
	workingBranch := flag.String("workingBranch", "", "Working branch which will be compared with the destination branch.")
	destinationBranch := flag.String("destinationBranch", defaultDestinationBranch, "Destination branch with which we will compare selected working branch.")
	vcsType := flag.String("vcs", "git", "The type of vcs which will be used for retrieving diff information.")
	path := flag.String("path", ".", "The type of vcs which will be used for retrieving diff information.")
	ext := flag.String("fileExt", "", "The type of extension of the files which we need to check.")
	pathToIgnoreFile := flag.String("pathToIgnoreFile", defaultIgnorePath, fmt.Sprintf("The path to file, where line-by-line written the list of paths which should be ignored. By default it's: %s", defaultIgnorePath))

	flag.Parse()

	if *vcsType == "" {
		log.Fatal(fmt.Errorf("The vcs should not be empty "))
	}

	if *path == "" {
		*path = "."
	}

	if *pathToIgnoreFile == "" {
		*pathToIgnoreFile = defaultIgnorePath
	}

	if *workingBranch == "" || *destinationBranch == "" {
		log.Fatal(fmt.Errorf("Working branch and destination branch should not be empty. Please make sure you define them. "))
	}

	_, err := services.LoadAvailableFeaturesInDir(*path, *ext, *pathToIgnoreFile)
	if err != nil {
		log.Fatal(err)
	}

	loadVcs(*vcsType)

	resultString := ""
	if *workingBranch != "" && *destinationBranch != "" {
		files, err := vcs.Diff(*path, *workingBranch, *destinationBranch)
		if err != nil {
			log.Fatal(err)
		}

		if len(files) > 0 {
			for _, file := range files {
				if len(services.PF.FoundFeaturesByFile[file]) > 0 {
					resultString += fmt.Sprintf("Changes in file: '%s' can potentially touch next features:\n", file)
					for _, feature := range services.PF.FoundFeaturesByFile[file] {
						resultString += "------------------\n"
						resultString += fmt.Sprintf("Feature: %s\n", feature.Name)
						resultString += fmt.Sprintf("Code path: %s:%d\n", feature.FilePath, feature.Line)
					}
					resultString += "------------------\n"
				}
			}

			if resultString == "" {
				resultString = "I found no features touched by these changes.\nPlease, make sure you define the features in these files"
			} else {
				resultString += fmt.Sprintf("Please make sure you test these features before you merge `%s` branch into `%s`.", *workingBranch, *destinationBranch)
			}
		}
	}

	fmt.Println(resultString)
}

func loadVcs(vcsType string) {
	switch vcsType {
	case "git":
		vcs = services.Git{}
		break
	}
}