package main

import (
	"context"
	"fmt"
	"io"
	"os"

	cli "github.com/urfave/cli/v2"

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
				Name:  "issue",
				Usage: "Issue statistics",
				Subcommands: []*cli.Command{
					{
						Name:  "count",
						Usage: "Count total number of issues",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "repo",
								Aliases:  []string{"r"},
								Usage:    "GitHub Repository (i.e. \"octocat/hello-world\")",
								Required: true,
							},
							&cli.StringSliceFlag{
								Name:    "is",
								Aliases: []string{"i"},
								Usage:   "Optional is filter (i.e. \"pr\")",
							},
							&cli.StringFlag{
								Name:    "params",
								Aliases: []string{"p"},
								Usage:   "Optional query parameters (i.e. \"merged:>=2020-10-08\")",
							},
						},

						Action: func(c *cli.Context) error {
							client := githubv4.NewClient(
								oauth2.NewClient(
									context.Background(),
									oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("token")}),
								),
							)

							return issueCount(client, c.String("repo"), c.StringSlice("is"), c.String("params"), os.Stdout)
						},
					},
				},
			},
			{
				Name:  "repo",
				Usage: "Repository statistics",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List all repositories",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "org",
								Aliases:  []string{"o"},
								Usage:    "GitHub Organization (i.e. \"octocat\")",
								Required: true,
							},
						},

						Action: func(c *cli.Context) error {
							client := githubv4.NewClient(
								oauth2.NewClient(
									context.Background(),
									oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("token")}),
								),
							)

							return repoList(client, c.String("org"), os.Stdout)
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

func issueCount(client *githubv4.Client, repo string, is []string, params string, w io.Writer) error {
	searchQuery := fmt.Sprintf("repo:%s %s", repo, params)

	// append every 'is:' condition
	for _, cond := range is {
		searchQuery = fmt.Sprintf("%s is:%s", searchQuery, cond)
	}

	fmt.Println(searchQuery)

	variables := map[string]interface{}{
		"searchQuery": githubv4.String(searchQuery),
	}

	var query struct {
		Search struct {
			IssueCount int
		} `graphql:"search(query: $searchQuery, type: ISSUE, first: 100)"`
	}

	err := client.Query(context.Background(), &query, variables)

	if err != nil {
		return err
	}

	fmt.Fprintln(w, query.Search.IssueCount)
	return nil
}

func repoList(client *githubv4.Client, org string, w io.Writer) error {
	var query struct {
		Search struct {
			Edges []struct {
				Node struct {
					Repository struct {
						NameWithOwner string
					} `graphql:"... on Repository"`
				}
			}
		} `graphql:"search(query: $searchQuery, type: REPOSITORY, first: 10)"`
	}

	searchQuery := githubv4.String(fmt.Sprintf("user:%s", org))

	fmt.Println(searchQuery)

	variables := map[string]interface{}{
		"searchQuery": searchQuery,
	}

	err := client.Query(context.Background(), &query, variables)

	if err != nil {
		return err
	}

	fmt.Fprintln(w, query.Search.Edges)
	return nil
}
