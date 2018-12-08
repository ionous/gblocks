package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNames(t *testing.T) {
	require.Equal(t, "Number", toFieldName("NUMBER"), "CAPITALS -> Capitals")
	require.Equal(t, "PascalCased", toFieldName("PASCAL_CASED"), "CAPITAL_NAME -> CapitalName")
}
