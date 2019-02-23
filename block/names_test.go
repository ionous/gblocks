package block

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNames(t *testing.T) {
	require.Equal(t, "Number", underscoreToPascal("NUMBER"), "CAPITALS -> Capitals")
	require.Equal(t, "PascalCased", underscoreToPascal("PASCAL_CASED"), "CAPITAL_NAME -> CapitalName")
	require.Equal(t, "", underscoreToPascal(""), "blank")
	require.Equal(t, "friendly", Item("FRIENDLY").Friendly())
}
