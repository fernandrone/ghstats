package list

import (
	"context"
	"fmt"
	"io"

	"github.com/shurcooL/githubv4"
)

func Repositories(client *githubv4.Client, org string, user string, filter string, w io.Writer) error {
	searchQuery := filter

	if org != "" {
		searchQuery = fmt.Sprintf("%s org:%s", searchQuery, org)
	}

	if user != "" {
		searchQuery = fmt.Sprintf("%s user:%s", searchQuery, user)
	}

	if org == "" && user == "" {
		login, err := login(client)

		if err != nil {
			return err
		}

		searchQuery = fmt.Sprintf("%s user:%s", searchQuery, login)
	}

	type RepositoryNode struct {
		Repository struct {
			NameWithOwner string
		} `graphql:"... on Repository"`
	}

	var query struct {
		Search struct {
			Nodes    []RepositoryNode
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"search(query: $searchQuery, type: REPOSITORY, first: 100, after: $commentsCursor)"`
	}

	variables := map[string]interface{}{
		"searchQuery":    githubv4.String(searchQuery),
		"commentsCursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	if err := client.Query(context.Background(), &query, variables); err != nil {
		return err
	}

	var repositories []RepositoryNode

	for {
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			return err
		}

		repositories = append(repositories, query.Search.Nodes...)

		if !query.Search.PageInfo.HasNextPage {
			break
		}

		variables["commentsCursor"] = githubv4.NewString(query.Search.PageInfo.EndCursor)
	}

	for _, repo := range repositories {
		fmt.Fprintln(w, repo.Repository.NameWithOwner)
	}

	// table := tablewriter.NewWriter(w)

	// var rows [][]string

	// for _, v := range allRepos {
	// 	rows = append(rows, []string{
	// 		v.Repository.NameWithOwner,
	// 	})
	// }

	// table.SetAutoWrapText(false)
	// table.SetAutoFormatHeaders(true)
	// table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	// table.SetAlignment(tablewriter.ALIGN_LEFT)
	// table.SetCenterSeparator("")
	// table.SetColumnSeparator("")
	// table.SetRowSeparator("")
	// table.SetHeaderLine(false)
	// table.SetBorder(false)
	// table.SetTablePadding("\t")
	// table.SetNoWhiteSpace(true)

	// table.AppendBulk(rows)
	// table.Render()

	return nil
}

func login(client *githubv4.Client) (string, error) {
	var query struct {
		Viewer struct {
			Login string
		}
	}

	if err := client.Query(context.Background(), &query, nil); err != nil {
		return "", err
	}

	return query.Viewer.Login, nil
}
