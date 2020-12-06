package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhpAnalysis_GetNeedleKey(t *testing.T) {
	var php = PhpAnalysis{}

	assert.Equal(t, regexNeedlePhp, php.GetNeedleKey())
}

func TestPhpAnalysis_ExtractImport(t *testing.T) {
	var (
		php   = PhpAnalysis{}
		cases = map[string]string{
			`namespace App\Test\Test\Test;`:    `App\Test\Test\Test`,
			`namespace App\\Test\\Test\\Test;`: `App\\Test\\Test\\Test`,
		}
	)

	for text, expected := range cases {
		res, err := php.ExtractMainEntrypointSrc(text)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	}
}

func TestPhpAnalysis_ExtractImportUsage(t *testing.T) {
	var (
		php   = PhpAnalysis{}
		cases = map[string]string{
			`namespace App\Test\Test\Test`:    `App\Test\Test\Test`,
			`namespace App\\Test\\Test\\Test`: `App\\Test\\Test\\Test`,
			`App\\Test\\Test\\Test $test`:     `App\\Test\\Test\\Test`,
			`App\Test\Test\Test $test`:        `App\Test\Test\Test`,
			`App\Test::$test`:                 `App\Test`,
			`App\\Test::$test`:                `App\\Test`,
			`App\\Test::method()`:             `App\\Test`,
			`App\Test::method()`:              `App\Test`,
		}
	)

	for text, expected := range cases {
		res, err := php.ExtractImportUsage(text)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	}
}

func TestPhpAnalysis_FindMainEntrypoint(t *testing.T) {
	var (
		php   = PhpAnalysis{}
		cases = map[string]string{
			`class TestClass `:                         `TestClass`,
			`class TestClass extends OtherTestClass {`: `TestClass`,
			`class TestClass implements SomeInterface`: `TestClass`,
			`abstract class TestClass `:                `TestClass`,
			`abstract class TestClass`:                 `TestClass`,
			`interface TestInterface`:                  `TestInterface`,
			`interface TestInterface `:                 `TestInterface`,
			`trait TestTrait `:                         `TestTrait`,
			`trait TestTrait`:                          `TestTrait`,
		}
	)

	for useCase, expectedValue := range cases {
		res, err := php.FindMainEntrypointTitle(useCase)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, res)
	}
}
