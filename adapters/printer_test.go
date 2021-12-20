package adapters_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jimenezmaximiliano/migrations/adapters"
)

func TestPrinting(test *testing.T) {
	test.Parallel()

	buffer := bytes.Buffer{}

	printer := adapters.PrinterAdapter{}
	err := printer.Print(&buffer, "%d", 1)
	require.Nil(test, err)

	assert.Equal(test, "1", buffer.String())
}
