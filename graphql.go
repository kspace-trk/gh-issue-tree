package main

import (
	"fmt"

	graphql "github.com/cli/go-gh/v2/pkg/api"
)

const maxDepth = 8

type issueResponse struct {
	Repository struct {
		Issue issueNode `json:"issue"`
	} `json:"repository"`
}

type issueNode struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	State     string `json:"state"`
	URL       string `json:"url"`
	Body      string `json:"body"`
	Assignees struct {
		Nodes []struct {
			Login string `json:"login"`
		} `json:"nodes"`
	} `json:"assignees"`
	Labels struct {
		Nodes []struct {
			Name string `json:"name"`
		} `json:"nodes"`
	} `json:"labels"`
	SubIssues struct {
		Nodes []struct {
			Number    int    `json:"number"`
			Title     string `json:"title"`
			State     string `json:"state"`
			URL       string `json:"url"`
			Body      string `json:"body"`
			Assignees struct {
				Nodes []struct {
					Login string `json:"login"`
				} `json:"nodes"`
			} `json:"assignees"`
			Labels struct {
				Nodes []struct {
					Name string `json:"name"`
				} `json:"nodes"`
			} `json:"labels"`
			SubIssuesSummary struct {
				Total int `json:"total"`
			} `json:"subIssuesSummary"`
		} `json:"nodes"`
		PageInfo struct {
			HasNextPage bool   `json:"hasNextPage"`
			EndCursor   string `json:"endCursor"`
		} `json:"pageInfo"`
	} `json:"subIssues"`
}

const issueQuery = `
query($owner: String!, $repo: String!, $number: Int!) {
  repository(owner: $owner, name: $repo) {
    issue(number: $number) {
      number title state url body
      assignees(first: 10) { nodes { login } }
      labels(first: 10) { nodes { name } }
      subIssues(first: 50) {
        nodes {
          number title state url body
          assignees(first: 10) { nodes { login } }
          labels(first: 10) { nodes { name } }
          subIssuesSummary { total }
        }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
}
`

const issueQueryWithCursor = `
query($owner: String!, $repo: String!, $number: Int!, $cursor: String!) {
  repository(owner: $owner, name: $repo) {
    issue(number: $number) {
      subIssues(first: 50, after: $cursor) {
        nodes {
          number title state url body
          assignees(first: 10) { nodes { login } }
          labels(first: 10) { nodes { name } }
          subIssuesSummary { total }
        }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
}
`

func nodeToIssue(n issueNode) Issue {
	issue := Issue{
		Number: n.Number,
		Title:  n.Title,
		State:  n.State,
		URL:    n.URL,
		Body:   n.Body,
	}
	for _, a := range n.Assignees.Nodes {
		issue.Assignees = append(issue.Assignees, a.Login)
	}
	for _, l := range n.Labels.Nodes {
		issue.Labels = append(issue.Labels, l.Name)
	}
	return issue
}

func FetchIssueTree(client *graphql.GraphQLClient, owner, repo string, number int) (Issue, error) {
	return fetchIssueRecursive(client, owner, repo, number, 0)
}

func fetchIssueRecursive(client *graphql.GraphQLClient, owner, repo string, number, depth int) (Issue, error) {
	if depth > maxDepth {
		return Issue{}, fmt.Errorf("maximum depth %d exceeded", maxDepth)
	}

	var resp issueResponse
	variables := map[string]interface{}{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	}

	err := client.Do(issueQuery, variables, &resp)
	if err != nil {
		return Issue{}, fmt.Errorf("fetching issue #%d: %w", number, err)
	}

	node := resp.Repository.Issue
	issue := nodeToIssue(node)

	// Collect sub-issue nodes from first page
	type subIssueNode struct {
		Number           int
		Title            string
		State            string
		URL              string
		Body             string
		Assignees        []string
		Labels           []string
		SubIssuesTotal   int
	}

	var subNodes []subIssueNode
	for _, sn := range node.SubIssues.Nodes {
		si := subIssueNode{
			Number: sn.Number,
			Title:  sn.Title,
			State:  sn.State,
			URL:    sn.URL,
			Body:   sn.Body,
			SubIssuesTotal: sn.SubIssuesSummary.Total,
		}
		for _, a := range sn.Assignees.Nodes {
			si.Assignees = append(si.Assignees, a.Login)
		}
		for _, l := range sn.Labels.Nodes {
			si.Labels = append(si.Labels, l.Name)
		}
		subNodes = append(subNodes, si)
	}

	// Paginate if needed
	pageInfo := node.SubIssues.PageInfo
	for pageInfo.HasNextPage {
		var pageResp issueResponse
		pageVars := map[string]interface{}{
			"owner":  owner,
			"repo":   repo,
			"number": number,
			"cursor": pageInfo.EndCursor,
		}
		err := client.Do(issueQueryWithCursor, pageVars, &pageResp)
		if err != nil {
			return Issue{}, fmt.Errorf("fetching sub-issues page for #%d: %w", number, err)
		}
		for _, sn := range pageResp.Repository.Issue.SubIssues.Nodes {
			si := subIssueNode{
				Number: sn.Number,
				Title:  sn.Title,
				State:  sn.State,
				URL:    sn.URL,
				Body:   sn.Body,
				SubIssuesTotal: sn.SubIssuesSummary.Total,
			}
			for _, a := range sn.Assignees.Nodes {
				si.Assignees = append(si.Assignees, a.Login)
			}
			for _, l := range sn.Labels.Nodes {
				si.Labels = append(si.Labels, l.Name)
			}
			subNodes = append(subNodes, si)
		}
		pageInfo = pageResp.Repository.Issue.SubIssues.PageInfo
	}

	// Recursively fetch sub-issues that have their own sub-issues
	for _, sn := range subNodes {
		if sn.SubIssuesTotal > 0 {
			child, err := fetchIssueRecursive(client, owner, repo, sn.Number, depth+1)
			if err != nil {
				return Issue{}, err
			}
			issue.SubIssues = append(issue.SubIssues, child)
		} else {
			issue.SubIssues = append(issue.SubIssues, Issue{
				Number:    sn.Number,
				Title:     sn.Title,
				State:     sn.State,
				URL:       sn.URL,
				Body:      sn.Body,
				Assignees: sn.Assignees,
				Labels:    sn.Labels,
			})
		}
	}

	return issue, nil
}
