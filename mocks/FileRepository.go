// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// FileRepository is an autogenerated mock type for the FileRepository type
type FileRepository struct {
	mock.Mock
}

// GetMigrationFilePaths provides a mock function with given fields: migrationsDirectoryAbsolutePath
func (_m *FileRepository) GetMigrationFilePaths(migrationsDirectoryAbsolutePath string) ([]string, error) {
	ret := _m.Called(migrationsDirectoryAbsolutePath)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(migrationsDirectoryAbsolutePath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(migrationsDirectoryAbsolutePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMigrationQuery provides a mock function with given fields: migrationAbsolutePath
func (_m *FileRepository) GetMigrationQuery(migrationAbsolutePath string) (string, error) {
	ret := _m.Called(migrationAbsolutePath)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(migrationAbsolutePath)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(migrationAbsolutePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}