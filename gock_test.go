package api

import (
	"net/http/httputil"
	"testing"

	"github.com/h2non/gock"
)

func assertGockNoUnmatchedRequests(t *testing.T) {
	t.Helper()
	if !gock.HasUnmatchedRequest() {
		return
	}
	t.Errorf("Assertion failed: gock.HasUnmatchedRequest() = true")

	for i, req := range gock.GetUnmatchedRequests() {

		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			t.Errorf("Failed to dump request %d: %v", i, err)
		} else {
			t.Logf("Unmatched request %d:\n%s", i, dump)
		}
	}
}
