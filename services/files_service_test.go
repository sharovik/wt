package services

import (
	"testing"

	"github.com/sharovik/wt/analysis"
	"github.com/sharovik/wt/dto"

	"github.com/stretchr/testify/assert"
)

func TestFindFeaturesInFile(t *testing.T) {
	analysis.InitAnalysisService(".php")
	_, err := AnalyseFile("", "test/file_not_exists.php")
	assert.Error(t, err)

	indexedFile, err := AnalyseFile("", "../test/test_functions.php")
	assert.NoError(t, err)

	assert.Equal(t, 3, len(indexedFile.Features))

	expectFeatures := []dto.Feature{
		{
			Name:     "test functionality",
			FilePath: "../test/test_functions.php",
			Line:     4,
		},
		{
			Name:     "some other functionality",
			FilePath: "../test/test_functions.php",
			Line:     12,
		},
		{
			Name:     "some other functionality",
			FilePath: "../test/test_functions.php",
			Line:     19,
		},
	}

	assert.Equal(t, expectFeatures, indexedFile.Features)
}

func TestFindAvailableFeaturesInDir(t *testing.T) {
	analysis.InitAnalysisService(".php")
	_, _, _, err := AnalyseTheCode("../some/dir", ".php", []string{""})
	assert.Error(t, err)

	foundResults, pathIndex, importsIndex, err := AnalyseTheCode("../test", ".php", []string{""})
	assert.NoError(t, err)
	assert.NotEmpty(t, foundResults)
	assert.NotEmpty(t, pathIndex)
	assert.NotEmpty(t, importsIndex)
	assert.Equal(t, 1, len(foundResults))
}
