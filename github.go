package main

import (
	"context"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"time"
)

type Star struct {
	Url           string
	Name          string
	NameWithOwner string
	Description   string
	License       string
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
						Name string
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

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)

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
			total++
			lng := "Unknown"
			if len(e.Node.Languages.Edges) > 0 {
				lng = e.Node.Languages.Edges[0].Node.Name
			}
			if _, ok := stars[lng]; !ok {
				stars[lng] = make([]Star, 0)
			}
			stars[lng] = append(stars[lng], Star{
				Url:           e.Node.Url,
				Name:          e.Node.Name,
				NameWithOwner: e.Node.NameWithOwner,
				Description:   e.Node.Description,
				License:       e.Node.LicenseInfo.Name,
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
