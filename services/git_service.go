package services

import (
	"fmt"
	"os/exec"
	"regexp"
)

type VcsInterface interface {
	Diff(path string, branch1 string, branch2 string) (files []string, err error)
	ParseDiffOutput(output string) (files []string, err error)
}

type Git struct {

}

const regexpParseFileNames = `(?im)(?:[A-Z]\s+)(.*)`

func (g Git) Diff(path string, branch1 string, branch2 string) (files []string, err error) {
	cmd := exec.Command("git", "diff", "--name-status", fmt.Sprintf("%s..%s", branch1, branch2))
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Received git error: %s", err.Error()))
		return nil, err
	}

	files, err = g.ParseDiffOutput(string(output))
	if err != nil {
		return nil, err
	}

	return
}

func (Git) ParseDiffOutput(output string) (files []string, err error) {
	re, err := regexp.Compile(regexpParseFileNames)

	if err != nil {
		return
	}

	matches := re.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		files = append(files, match[1])
	}

	return
}
