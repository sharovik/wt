package analysis

import (
	"fmt"
	"regexp"
	"strings"
)

//PhpAnalysis the default analysis struct
type PhpAnalysis struct {
}

const (
	regexNeedlePhp = `(?im)((?:\\{1,2}\w+|\w+\\{1,2})(?:\w+\\{0,2})+)`
	regexNamespace = `(?im)(?:namespace\s+)(.*)(?:[;])`
	regexEntryPoint = `(?im)(?:class\s+)(\w+)|(?:trait\s+)(\w+)|(?:interface\s+)(\w+)`
)

func (a PhpAnalysis) GetNeedleKey() string {
	return regexNeedlePhp
}

func (a PhpAnalysis) ExtractMainEntrypointSrc(text string) (res string, err error) {
	re, err := regexp.Compile(regexNamespace)

	if err != nil {
		return
	}

	matches := re.FindStringSubmatch(text)

	if len(matches) != 2 {
		return
	}

	return strings.TrimSpace(matches[1]), nil
}

func (a PhpAnalysis) ExtractImportUsage(text string) (res string, err error) {
	re, err := regexp.Compile(regexNeedlePhp)

	if err != nil {
		return
	}

	matches := re.FindStringSubmatch(text)

	if len(matches) != 2 {
		return
	}

	if strings.Contains(matches[1], "namespace") {
		return
	}

	return matches[1], nil
}

func (a PhpAnalysis) FindMainEntrypointTitle(text string) (res string, err error) {
	re, err := regexp.Compile(regexEntryPoint)
	if err != nil {
		return
	}

	matches := re.FindStringSubmatch(text)

	if len(matches) < 2 {
		return
	}

	//The class condition
	if matches[1] != "" {
		return matches[1], nil
	}

	//The trait condition
	if matches[2] != "" {
		return matches[2], nil
	}

	//The interface condition
	if matches[3] != "" {
		return matches[3], nil
	}

	return
}

func (a PhpAnalysis) GenerateEntryPointSrc(src string, entrypointTitle string) string {
	return fmt.Sprintf("%s\\%s", src, entrypointTitle)
}