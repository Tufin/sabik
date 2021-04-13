package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/sabik/api"
)

func TestOriginHasResponse(t *testing.T) {

	require.False(t, api.OriginHasResponse(api.OriginNginx))
}
