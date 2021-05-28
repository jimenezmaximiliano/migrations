package helpers

import (
	"testing"
)

func TestAddingATrailingSlashOnAPathWithoutIt(test *testing.T) {
	const path = "/tmp"
	const expectedResult = "/tmp/"
	result := AddTrailingSlashToPathIfNeeded(path)
	if result != expectedResult {
		test.Errorf("expected %s but got %s", expectedResult, result)
	}
}

func TestAddingATrailingSlashOnAPathWithIt(test *testing.T) {
	const path = "/tmp/"
	const expectedResult = "/tmp/"
	result := AddTrailingSlashToPathIfNeeded(path)
	if result != expectedResult {
		test.Errorf("expected %s but got %s", expectedResult, result)
	}
}

func TestAddingATrailingSlashOnAnEmptyPath(test *testing.T) {
	const path = ""
	const expectedResult = "/"
	result := AddTrailingSlashToPathIfNeeded(path)
	if result != expectedResult {
		test.Errorf("expected %s but got %s", expectedResult, result)
	}
}
