package count

import (
	"context"
	"fmt"
	"io"

	"github.com/shurcooL/githubv4"
)

func PullRequests(client *githubv4.Client, repo string, filter string, w io.Writer) error {
	searchQuery := fmt.Sprintf("%s is:pr", filter)

	if repo != "" {
		searchQuery = fmt.Sprintf("%s repo:%s", searchQuery, repo)
	}

	variables := map[string]interface{}{
		"searchQuery": githubv4.String(searchQuery),
	}

	var query struct {
		Search struct {
			IssueCount int
		} `graphql:"search(query: $searchQuery, type: ISSUE, first: 1)"`
	}

	if err := client.Query(context.Background(), &query, variables); err != nil {
		return err
	}

	fmt.Fprintln(w, query.Search.IssueCount)
	return nil
}

func Issues(client *githubv4.Client, repo string, filter string, w io.Writer) error {
	searchQuery := fmt.Sprintf("%s is:issue", filter)

	if repo != "" {
		searchQuery = fmt.Sprintf("%s repo:%s", searchQuery, repo)
	}

	variables := map[string]interface{}{
		"searchQuery": githubv4.String(searchQuery),
	}

	var query struct {
		Search struct {
			IssueCount int
		} `graphql:"search(query: $searchQuery, type: ISSUE, first: 1)"`
	}

	if err := client.Query(context.Background(), &query, variables); err != nil {
		return err
	}

	fmt.Fprintln(w, query.Search.IssueCount)
	return nil
}
