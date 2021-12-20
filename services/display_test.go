package services_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/models"
	"github.com/jimenezmaximiliano/migrations/services"
)

func TestDisplayError(test *testing.T) {
	test.Parallel()

	var result string
	printer := &printLogger{
		Log: &result,
	}
	service := services.NewDisplayService(printer)

	service.DisplayError(errors.New("oops"))

	assert.Contains(test, result, "oops")
}

func TestDisplayErrorWithMessage(test *testing.T) {
	test.Parallel()

	var result string
	printer := &printLogger{
		Log: &result,
	}
	service := services.NewDisplayService(printer)

	service.DisplayErrorWithMessage(errors.New("oops"), "message")

	assert.Contains(test, result, "oops")
	assert.Contains(test, result, "message")
}

func TestDisplayHelp(test *testing.T) {
	test.Parallel()

	var result string
	printer := &printLogger{
		Log: &result,
	}
	service := services.NewDisplayService(printer)

	service.DisplayHelp()

	assert.Contains(test, result, "Usage")
	assert.Contains(test, result, "Examples")
	assert.Contains(test, result, "Documentation")
}

func TestDisplaySetupError(test *testing.T) {
	test.Parallel()

	var result string
	printer := &printLogger{
		Log: &result,
	}
	service := services.NewDisplayService(printer)

	service.DisplaySetupError(errors.New("oops"))

	assert.Contains(test, result, "oops")
}

func TestDisplayGeneralError(test *testing.T) {
	test.Parallel()

	var result string
	printer := &printLogger{
		Log: &result,
	}
	service := services.NewDisplayService(printer)

	service.DisplayGeneralError(errors.New("oops"))

	assert.Contains(test, result, "oops")
}

func TestDisplayingRunMigrationsWithoutMigrations(test *testing.T) {
	test.Parallel()

	var result string
	printer := &printLogger{
		Log: &result,
	}
	service := services.NewDisplayService(printer)

	service.DisplayRunMigrations(models.Collection{})

	assert.Contains(test, result, "No migrations to run")
}

func TestDisplayingRunMigrationsWithMigrations(test *testing.T) {
	test.Parallel()

	var result string
	printer := &printLogger{
		Log: &result,
	}
	service := services.NewDisplayService(printer)

	migrations := models.Collection{}

	migration1, err := models.NewMigration("/tmp/1_gophers.sql", "SELECT 1;", models.StatusSuccessful)
	require.Nil(test, err)
	err = migrations.Add(migration1)
	require.Nil(test, err)

	migration2, err := models.NewMigration("/tmp/2_fusilli_jerry.sql", "SELECT 1;", models.StatusFailed)
	require.Nil(test, err)
	err = migrations.Add(migration2)
	require.Nil(test, err)

	migration3, err := models.NewMigration("/tmp/3_walrus.sql", "SELECT 1;", models.StatusNotRun)
	require.Nil(test, err)
	err = migrations.Add(migration3)
	require.Nil(test, err)

	service.DisplayRunMigrations(migrations)

	assert.Contains(test, result, "OK")
	assert.Contains(test, result, "1_gophers.sql")
	assert.Contains(test, result, "FAIL")
	assert.Contains(test, result, "2_fusilli_jerry.sql")
	assert.Contains(test, result, "INFO")
	assert.Contains(test, result, "3_walrus.sql")
}

type printLogger struct {
	Log *string
}

func (logger printLogger) Print(writer io.Writer, format string, a ...interface{}) error {
	*logger.Log += fmt.Sprintf(format, a...)

	return nil
}
