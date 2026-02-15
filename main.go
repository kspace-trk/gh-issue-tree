package main

import (
	"fmt"
	"os"
	"strconv"

	graphql "github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: gh issue-tree <issue-number>")
		os.Exit(1)
	}

	number, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid issue number %q\n", os.Args[1])
		os.Exit(1)
	}

	repo, err := repository.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: could not determine repository. Run this command from inside a git repository.")
		os.Exit(1)
	}

	client, err := graphql.DefaultGraphQLClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create GraphQL client: %v\n", err)
		os.Exit(1)
	}

	issue, err := FetchIssueTree(client, repo.Owner, repo.Name, number)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(RenderMarkdown(issue))
}
