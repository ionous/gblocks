package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNames(t *testing.T) {
	require.Equal(t, "Number", underscoreToPascal("NUMBER"), "CAPITALS -> Capitals")
	require.Equal(t, "PascalCased", underscoreToPascal("PASCAL_CASED"), "CAPITAL_NAME -> CapitalName")
}
