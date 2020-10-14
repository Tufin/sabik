package api_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/sabik/api"
)

func TestHTTPLog_GetQueryString(t *testing.T) {

	parsedURL, err := url.Parse("https://securecloud.tufin.io/auth/realms/generic-bank/protocol/openid-connect/auth?client_id=express&state=7ba3a73c-9dc8-45d4-a77c-cc74e632303b")
	require.NoError(t, err)
	query, err := (&api.HTTPLog{QueryString: parsedURL.RawQuery}).GetQueryString()
	require.NoError(t, err)
	require.Len(t, query, 2)
	curr := query["client_id"]
	require.Len(t, curr, 1)
	require.Equal(t, "express", curr[0])
	curr = query["state"]
	require.Len(t, curr, 1)
	require.Equal(t, "7ba3a73c-9dc8-45d4-a77c-cc74e632303b", curr[0])
}
