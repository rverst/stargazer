package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
)

const defaultGraphQLURL = "https://api.github.com/graphql"

// graphQLURL is a variable so that tests can point the client at a fake server.
var graphQLURL = defaultGraphQLURL

// bearerTransport adds a static access token to every request.
type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req)
}

type Star struct {
	Url           string
	Name          string
	NameWithOwner string
	Description   string
	License       string
	LicenseUrl    string
	Stars         int
	Archived      bool
	StarredAt     time.Time
}

var query struct {
	User struct {
		StarredRepositories struct {
			IsOverLimit bool
			TotalCount  int
			Edges       []struct {
				StarredAt time.Time
				Node      struct {
					Description string
					Languages   struct {
						Edges []struct {
							Node struct {
								Name string
							}
						}
					} `graphql:"languages(first: $lc, orderBy: {field: SIZE, direction: DESC})"`
					LicenseInfo struct {
						Name     string
						Nickname string
						Url      string
					}
					IsArchived     bool
					IsPrivate      bool
					Name           string
					NameWithOwner  string
					StargazerCount int
					Url            string
				}
			}
			PageInfo struct {
				EndCursor   string
				HasNextPage bool
			}
		} `graphql:"starredRepositories(first: $count, orderBy: {field: STARRED_AT, direction: DESC}, after: $cursor)"`
	} `graphql:"user(login: $login)"`
}

func fetchStars(user string, token string) (stars map[string][]Star, total int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	httpClient := &http.Client{Transport: bearerTransport{token: token, base: http.DefaultTransport}}

	client := githubv4.NewEnterpriseClient(graphQLURL, httpClient)

	vars := map[string]interface{}{
		"login":  githubv4.String(user),
		"lc":     githubv4.Int(1),
		"count":  githubv4.Int(50),
		"cursor": githubv4.String(""),
	}

	stars = make(map[string][]Star)
	total = 0
	for {
		err = client.Query(ctx, &query, vars)
		if err != nil {
			return
		}

		for _, e := range query.User.StarredRepositories.Edges {
			// skip private repos
			if e.Node.IsPrivate {
				continue
			}
			if isIgnored(e.Node.NameWithOwner) {
				continue
			}

			total++
			lng := "Unknown"
			if len(e.Node.Languages.Edges) > 0 {
				lng = e.Node.Languages.Edges[0].Node.Name
			}
			if _, ok := stars[lng]; !ok {
				stars[lng] = make([]Star, 0)
			}

			var lic string
			if e.Node.LicenseInfo.Nickname != "" {
				lic = e.Node.LicenseInfo.Nickname
			} else {
				lic = e.Node.LicenseInfo.Name
			}
			if strings.ToLower(lic) == "other" {
				lic = ""
			}

			stars[lng] = append(stars[lng], Star{
				Url:           e.Node.Url,
				Name:          e.Node.Name,
				NameWithOwner: e.Node.NameWithOwner,
				Description:   e.Node.Description,
				License:       lic,
				LicenseUrl:    e.Node.LicenseInfo.Url,
				Stars:         e.Node.StargazerCount,
				Archived:      e.Node.IsArchived,
				StarredAt:     e.StarredAt,
			})
		}

		if !query.User.StarredRepositories.PageInfo.HasNextPage {
			break
		}
		vars["cursor"] = githubv4.String(query.User.StarredRepositories.PageInfo.EndCursor)
	}

	return
}
