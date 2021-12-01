package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddingATrailingSlashOnAPathWithoutIt(test *testing.T) {
	const path = "/tmp"
	const expectedResult = "/tmp/"
	result := AddTrailingSlashToPathIfNeeded(path)

	assert.Equal(test, expectedResult, result)
}

func TestAddingATrailingSlashOnAPathWithIt(test *testing.T) {
	const path = "/tmp/"
	const expectedResult = "/tmp/"
	result := AddTrailingSlashToPathIfNeeded(path)

	assert.Equal(test, expectedResult, result)
}

func TestAddingATrailingSlashOnAnEmptyPath(test *testing.T) {
	const path = ""
	const expectedResult = "/"
	result := AddTrailingSlashToPathIfNeeded(path)

	assert.Equal(test, expectedResult, result)
}
