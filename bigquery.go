// Package databigquery is a togo data backend that queries Google BigQuery via
// the REST jobs.query API. Select with `togo provider:use data bigquery`.
//
//   BIGQUERY_PROJECT   GCP project id
//   BIGQUERY_TOKEN     OAuth2 access token (secret)
//   BIGQUERY_ENDPOINT  API base (default https://bigquery.googleapis.com)
package databigquery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/togo-framework/data"
	"github.com/togo-framework/providers"
	"github.com/togo-framework/togo"
)

func init() {
	togo.RegisterProviderFunc("data-bigquery", togo.PriorityService+1, func(k *togo.Kernel) error {
		providers.Use(k, providers.CapData, "bigquery", newBQ(k), false)
		if k.Log != nil {
			k.Log.Info("plugin active", "plugin", "data-bigquery")
		}
		return nil
	})
}

type bq struct {
	endpoint string
	project  string
	token    string
	hc       *http.Client
}

func newBQ(k *togo.Kernel) *bq {
	return &bq{
		endpoint: strings.TrimRight(providers.Value(k, providers.CapData, "bigquery", "endpoint", "https://bigquery.googleapis.com", false), "/"),
		project:  providers.Value(k, providers.CapData, "bigquery", "project", "", false),
		token:    providers.Value(k, providers.CapData, "bigquery", "token", "", true),
		hc:       &http.Client{Timeout: 60 * time.Second},
	}
}

func (b *bq) Query(ctx context.Context, query string, _ ...any) ([]data.Row, error) {
	if b.project == "" || b.token == "" {
		return nil, fmt.Errorf("data-bigquery: set BIGQUERY_PROJECT and BIGQUERY_TOKEN")
	}
	body, _ := json.Marshal(map[string]any{"query": query, "useLegacySql": false})
	url := b.endpoint + "/bigquery/v2/projects/" + b.project + "/queries"
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+b.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := b.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out struct {
		Schema struct {
			Fields []struct {
				Name string `json:"name"`
			} `json:"fields"`
		} `json:"schema"`
		Rows []struct {
			F []struct {
				V any `json:"v"`
			} `json:"f"`
		} `json:"rows"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		msg := "http " + resp.Status
		if out.Error != nil {
			msg = out.Error.Message
		}
		return nil, fmt.Errorf("bigquery: %s", msg)
	}
	rows := make([]data.Row, 0, len(out.Rows))
	for _, r := range out.Rows {
		row := make(data.Row, len(out.Schema.Fields))
		for i, f := range out.Schema.Fields {
			if i < len(r.F) {
				row[f.Name] = r.F[i].V
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
