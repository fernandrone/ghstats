package describe

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/shurcooL/githubv4"
)

type AuthorDescription struct {
	Login   string
	Matches int
}

type AuthorTable []AuthorDescription

// Len implements sort.Interace
func (a AuthorTable) Len() int {
	return len(a)
}

// Less implements sort.Interace
func (a AuthorTable) Less(i, j int) bool {
	return a[i].Matches <= a[j].Matches
}

// Swap implements sort.Interace
func (a AuthorTable) Swap(i, j int) {
	tmp := a[i]
	a[i] = a[j]
	a[j] = tmp
}

func Authors(client *githubv4.Client, repo string, filter string, w io.Writer) error {
	searchQuery := fmt.Sprintf("%s repo:%s", filter, repo)

	type PullRequestNode struct {
		PullRequest struct {
			Author struct {
				Login string
			}
		} `graphql:"... on PullRequest"`
	}

	var query struct {
		Search struct {
			IssueCount int
			Nodes      []PullRequestNode
			PageInfo   struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"search(query: $searchQuery, type: ISSUE, first: 100, after: $commentsCursor)"`
	}

	variables := map[string]interface{}{
		"searchQuery":    githubv4.String(searchQuery),
		"commentsCursor": (*githubv4.String)(nil), // Null after argument to get first page.
	}

	if err := client.Query(context.Background(), &query, variables); err != nil {
		return err
	}

	if query.Search.IssueCount > 1000 {
		fmt.Fprintln(os.Stderr, "Warning: there are more than 1000 results, try using filters to narrow down your search https://docs.github.com/en/rest/reference/search#about-the-search-api")
	}

	var pullRequests []PullRequestNode

	for {
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			return err
		}

		pullRequests = append(pullRequests, query.Search.Nodes...)

		if !query.Search.PageInfo.HasNextPage {
			break
		}

		variables["commentsCursor"] = githubv4.NewString(query.Search.PageInfo.EndCursor)
	}

	authors := make(map[string]int)

	// increment the counter for every match for each author
	for _, node := range pullRequests {
		if val, ok := authors[node.PullRequest.Author.Login]; !ok {
			authors[node.PullRequest.Author.Login] = 1
		} else {
			authors[node.PullRequest.Author.Login] = val + 1
		}
	}

	var authorsTable AuthorTable

	for login, matches := range authors {
		authorsTable = append(authorsTable, AuthorDescription{login, matches})
	}

	// sort the list of authors
	sort.Sort(sort.Reverse(authorsTable))

	table := tablewriter.NewWriter(w)

	var rows [][]string

	for _, author := range authorsTable {
		rows = append(rows, []string{
			author.Login,
			fmt.Sprintf("%d", author.Matches),
		})
	}

	table.SetHeader([]string{"AUTHORS", "MATCHES"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	table.AppendBulk(rows)
	table.Render()

	return nil
}
