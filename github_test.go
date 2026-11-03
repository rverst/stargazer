package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestFetchStarsAuthenticatesAndPaginates(t *testing.T) {
	const pages = 3
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), "Bearer secret"; got != want {
			t.Errorf("Authorization = %q, want %q", got, want)
		}

		var body struct {
			Variables map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		wantCursor := ""
		if calls > 0 {
			wantCursor = "cursor-" + strconv.Itoa(calls)
		}
		if got := body.Variables["cursor"]; got != wantCursor {
			t.Errorf("request %d: cursor = %v, want %q", calls, got, wantCursor)
		}
		calls++

		name := "owner/repo-" + strconv.Itoa(calls)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"user": map[string]any{
			"starredRepositories": map[string]any{
				"edges": []any{map[string]any{
					"starredAt": "2026-01-02T03:04:05Z",
					"node": map[string]any{
						"nameWithOwner":  name,
						"stargazerCount": 7,
						"languages": map[string]any{"edges": []any{
							map[string]any{"node": map[string]any{"name": "Go"}},
						}},
						"licenseInfo": map[string]any{"name": "MIT License"},
					},
				}},
				"pageInfo": map[string]any{
					"endCursor":   "cursor-" + strconv.Itoa(calls),
					"hasNextPage": calls < pages,
				},
			},
		}}})
	}))
	defer srv.Close()

	graphQLURL = srv.URL
	defer func() { graphQLURL = defaultGraphQLURL }()

	stars, total, err := fetchStars("rverst", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if calls != pages {
		t.Errorf("requests = %d, want %d", calls, pages)
	}
	if total != pages || len(stars["Go"]) != pages {
		t.Errorf("total = %d, Go stars = %d, want %d", total, len(stars["Go"]), pages)
	}
	if got := stars["Go"][0]; got.License != "MIT License" || got.Stars != 7 {
		t.Errorf("unexpected star: %+v", got)
	}
}
