package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMiddleware(t *testing.T) {

	require.IsType(t, &EmptyMiddleware{}, CreateMiddleware())
}
