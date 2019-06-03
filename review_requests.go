package main

import (
	"context"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type edge struct {
	Node struct {
		PullRequest struct {
			Repository struct {
				NameWithOwner githubv4.String
			}
			Number githubv4.Int
			URL    githubv4.String
			Title  githubv4.String
		} `graphql:"... on PullRequest"`
	}
}

type query struct {
	Search struct {
		IssueCount int
		Edges      []edge
	} `graphql:"search(query: \"type:pr state:open team-review-requested:reload/developers\", type: ISSUE, first: 100)"`
}

func reviewRequests() ([]edge, int, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var query query

	err := client.Query(context.Background(), &query, nil)

	if err != nil {
		return []edge{}, 0, err
	}

	return query.Search.Edges, query.Search.IssueCount, nil
}
