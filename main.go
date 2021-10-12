package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/fernandrone/ghstats/pkg/count"
	"github.com/fernandrone/ghstats/pkg/describe"
	"github.com/fernandrone/ghstats/pkg/list"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	app := &cli.App{
		Name:  "ghstats",
		Usage: "github stats command line client",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				EnvVars: []string{"GITHUB_TOKEN"},
				Usage:   "GitHub Token",
			},
		},

		Commands: []*cli.Command{
			{
				Name:  "count",
				Usage: "Count GitHub objects",
				Subcommands: []*cli.Command{
					{
						Name:  "prs",
						Usage: "Count total number of pull requests in a repository",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "repo",
								Aliases:  []string{"r"},
								Usage:    "GitHub Repo (i.e. \"github/octocat\")",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "filter",
								Aliases: []string{"f"},
								Usage:   "Optional filters (i.e. \"is:merged merged:>=2020-10-08\")",
							},
						},

						Action: func(c *cli.Context) error {
							client := githubv4.NewClient(
								oauth2.NewClient(
									context.Background(),
									oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("token")}),
								),
							)

							return count.PullRequests(client, c.String("repo"), c.String("filter"), os.Stdout)
						},
					},
					{
						Name:  "issues",
						Usage: "Count total number of issues in a repository",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "repo",
								Aliases: []string{"r"},
								Usage:   "GitHub Repo (i.e. \"github/octocat\")",
							},
							&cli.StringFlag{
								Name:    "filter",
								Aliases: []string{"f"},
								Usage:   "Optional filters (i.e. \"is:closed closed:>=2020-10-08\")",
							},
						},

						Action: func(c *cli.Context) error {
							client := githubv4.NewClient(
								oauth2.NewClient(
									context.Background(),
									oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("token")}),
								),
							)

							return count.Issues(client, c.String("repo"), c.String("filter"), os.Stdout)
						},
					},
				},
			},
			{
				Name:  "list",
				Usage: "List GitHub objects",
				Subcommands: []*cli.Command{
					{
						Name:  "repo",
						Usage: "List repositories. Defaults to listing the client's repositories unless the 'user' or 'org' filters are used.",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "org",
								Aliases: []string{"o"},
								Usage:   "GitHub Organization (i.e. \"github\")",
							},
							&cli.StringFlag{
								Name:    "user",
								Aliases: []string{"u"},
								Usage:   "GitHub User (i.e. \"octocat\")",
							},
							&cli.StringFlag{
								Name:    "filter",
								Aliases: []string{"f"},
								Usage:   "Optional filters (i.e. \"stars:500\" matches repositories with exactly 500 stars)",
							},
						},

						Action: func(c *cli.Context) error {
							client := githubv4.NewClient(
								oauth2.NewClient(
									context.Background(),
									oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("token")}),
								),
							)

							return list.Repositories(client, c.String("org"), c.String("user"), c.String("filter"), os.Stdout)
						},
					},
				},
			},
			{
				Name:  "describe",
				Usage: "Describe GitHub objects",
				Subcommands: []*cli.Command{
					{
						Name:  "authors",
						Usage: "Describe authors of pull requests and/or issues in a repository",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "repo",
								Aliases:  []string{"r"},
								Usage:    "GitHub Repo (i.e. \"github/octocat\")",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "filter",
								Aliases: []string{"f"},
								Usage:   "Optional filters (i.e. \"is:pr is:merged merged:>=2020-10-08\")",
							},
						},

						Action: func(c *cli.Context) error {
							client := githubv4.NewClient(
								oauth2.NewClient(
									context.Background(),
									oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("token")}),
								),
							)

							return describe.Authors(client, c.String("repo"), c.String("filter"), os.Stdout)
						},
					},
					{
						Name:  "issues",
						Usage: "Count total number of issues in a repository",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "repo",
								Aliases: []string{"r"},
								Usage:   "GitHub Repo (i.e. \"github/octocat\")",
							},
							&cli.StringFlag{
								Name:    "filter",
								Aliases: []string{"f"},
								Usage:   "Optional filters (i.e. \"is:closed closed:>=2020-10-08\")",
							},
						},

						Action: func(c *cli.Context) error {
							client := githubv4.NewClient(
								oauth2.NewClient(
									context.Background(),
									oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("token")}),
								),
							)

							return count.Issues(client, c.String("repo"), c.String("filter"), os.Stdout)
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// func countPR(client *githubv4.Client, repo string, filter string, w io.Writer) error {
// 	searchQuery := fmt.Sprintf("%s is:pr", filter)

// 	if repo != "" {
// 		searchQuery = fmt.Sprintf("%s repo:%s", searchQuery, repo)
// 	}

// 	// append every 'is:' condition
// 	for _, cond := range is {
// 		searchQuery = fmt.Sprintf("%s is:%s", searchQuery, cond)
// 	}

// 	variables := map[string]interface{}{
// 		"searchQuery":    githubv4.String(searchQuery),
// 		"commentsCursor": (*githubv4.String)(nil), // Null after argument to get first page.
// 	}

// 	// type PullRequestNode struct {
// 	// 	PullRequest struct {
// 	// 		Repository struct {
// 	// 			NameWithOwner string
// 	// 		}
// 	// 	} `graphql:"... on PullRequest"`
// 	// }

// 	var query struct {
// 		Search struct {
// 			IssueCount int
// 			// Nodes      []PullRequestNode
// 			// PageInfo   struct {
// 			// 	EndCursor   githubv4.String
// 			// 	HasNextPage bool
// 			// }
// 		} `graphql:"search(query: $searchQuery, type: ISSUE, first: 1, after: $commentsCursor)"`
// 	}

// 	err := client.Query(context.Background(), &query, variables)
// 	if err != nil {
// 		return err
// 	}

// 	// var allPRs []PullRequestNode

// 	table := tablewriter.NewWriter(w)
// 	var rows [][]string

// 	// if repo == "" {
// 	// 	for {

// 	// 		err := client.Query(context.Background(), &query, variables)
// 	// 		if err != nil {
// 	// 			return err
// 	// 		}

// 	// 		allPRs = append(allPRs, query.Search.Nodes...)

// 	// 		if !query.Search.PageInfo.HasNextPage {
// 	// 			break
// 	// 		}

// 	// 		variables["commentsCursor"] = githubv4.NewString(query.Search.PageInfo.EndCursor)
// 	// 	}

// 	// 	prMap := make(map[string]int)

// 	// 	for _, node := range allPRs {
// 	// 		if val, ok := prMap[node.PullRequest.Repository.NameWithOwner]; !ok {
// 	// 			prMap[node.PullRequest.Repository.NameWithOwner] = 1
// 	// 		} else {
// 	// 			prMap[node.PullRequest.Repository.NameWithOwner] = val + 1
// 	// 		}
// 	// 	}

// 	// 	for k, v := range prMap {
// 	// 		rows = append(rows, []string{
// 	// 			k,
// 	// 			fmt.Sprintf("%5d", v),
// 	// 		})
// 	// 	}
// 	// } else {
// 	rows = append(rows, []string{
// 		repo,
// 		fmt.Sprintf("%5d", query.Search.IssueCount),
// 	})
// 	// }

// 	table.SetAutoWrapText(false)
// 	table.SetAutoFormatHeaders(true)
// 	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
// 	table.SetAlignment(tablewriter.ALIGN_LEFT)
// 	table.SetCenterSeparator("")
// 	table.SetColumnSeparator("")
// 	table.SetRowSeparator("")
// 	table.SetHeaderLine(false)
// 	table.SetBorder(false)
// 	table.SetTablePadding("\t")
// 	table.SetNoWhiteSpace(true)

// 	table.AppendBulk(rows)
// 	table.Render()

// 	return nil
// }
