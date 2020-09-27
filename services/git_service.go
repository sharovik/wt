package services

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

//VcsInterface the vcs interface
type VcsInterface interface {
	Diff(path string, branch1 string, branch2 string) (files []string, err error)
	parseDiffOutput(output string) (files []string, err error)
}

//Git the object struct of Git vcs
type Git struct {
}

const regexpParseFileNames = `(?im)(?:[A-Z]\s+)(.*)`

//Diff the main method which is used for diff files get
func (g Git) Diff(path string, branch1 string, branch2 string) (files []string, err error) {
	cmd := exec.Command("git", "diff", "--name-status", fmt.Sprintf("%s..%s", branch1, branch2))
	cmd.Dir = path
	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(fmt.Sprintf("Received git error: %s; Stderr: %s", err.Error(), errbuf.String()))
		return nil, err
	}

	files, err = g.parseDiffOutput(string(output))
	if err != nil {
		return nil, err
	}

	return
}

func (Git) parseDiffOutput(output string) (files []string, err error) {
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
