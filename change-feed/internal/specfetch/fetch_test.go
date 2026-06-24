package specfetch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchReturnsDataAndETag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"abc"`)
		_, _ = w.Write([]byte(`{"openapi":"3.0.3"}`))
	}))
	defer srv.Close()

	res, err := New(srv.Client()).Fetch(context.Background(), srv.URL, "")
	if err != nil {
		t.Fatal(err)
	}
	if res.Unchanged || res.ETag != `"abc"` || len(res.Data) == 0 {
		t.Fatalf("got %+v", res)
	}
}
