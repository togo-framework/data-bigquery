package databigquery

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryParsesRows(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			w.WriteHeader(401)
			return
		}
		_, _ = w.Write([]byte(`{"schema":{"fields":[{"name":"n"},{"name":"city"}]},"rows":[{"f":[{"v":"42"},{"v":"riyadh"}]}]}`))
	}))
	defer srv.Close()
	b := &bq{endpoint: srv.URL, project: "p", token: "tok", hc: srv.Client()}
	rows, err := b.Query(context.Background(), "SELECT 42 n, 'riyadh' city")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0]["n"] != "42" || rows[0]["city"] != "riyadh" {
		t.Fatalf("rows = %v", rows)
	}
}

func TestQueryRequiresConfig(t *testing.T) {
	if _, err := (&bq{}).Query(context.Background(), "SELECT 1"); err == nil {
		t.Fatal("expected error without project/token")
	}
}
