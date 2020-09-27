package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFeaturesInFile(t *testing.T) {
	_, err := FindFeaturesInFile("", "test/file_not_exists.php")
	assert.Error(t, err)

	features, err := FindFeaturesInFile("", "../test/test_functions.php")
	assert.NoError(t, err)

	assert.Equal(t, 3, len(features))

	expectFeatures := []Feature{
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

	assert.Equal(t, expectFeatures, features)
	assert.Equal(t, 2, len(PF.FoundFeatures))
	assert.NotEmpty(t, PF.FoundFeatures["test functionality"])
	assert.NotEmpty(t, PF.FoundFeatures["some other functionality"])
	assert.Equal(t, 2, len(PF.FoundFeatures["some other functionality"]))
	CleanUpGlobalVars()
}

func TestFindAvailableFeaturesInDir(t *testing.T) {
	_, err := LoadAvailableFeaturesInDir("../some/dir", ".php", ".gitignore")
	assert.Error(t, err)

	foundResults, err := LoadAvailableFeaturesInDir("../test", ".php", ".gitignore")
	assert.NoError(t, err)
	assert.NotEmpty(t, foundResults)
	assert.Equal(t, 1, len(foundResults))
	CleanUpGlobalVars()
}
