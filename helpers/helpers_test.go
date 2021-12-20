package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jimenezmaximiliano/migrations/helpers"
)

func TestAddingATrailingSlashOnAPathWithoutIt(test *testing.T) {
	const path = "/tmp"
	const expectedResult = "/tmp/"
	result := helpers.AddTrailingSlashToPathIfNeeded(path)

	assert.Equal(test, expectedResult, result)
}

func TestAddingATrailingSlashOnAPathWithIt(test *testing.T) {
	const path = "/tmp/"
	const expectedResult = "/tmp/"
	result := helpers.AddTrailingSlashToPathIfNeeded(path)

	assert.Equal(test, expectedResult, result)
}

func TestAddingATrailingSlashOnAnEmptyPath(test *testing.T) {
	const path = ""
	const expectedResult = "/"
	result := helpers.AddTrailingSlashToPathIfNeeded(path)

	assert.Equal(test, expectedResult, result)
}
