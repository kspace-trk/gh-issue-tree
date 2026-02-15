package main

// Issue represents a GitHub issue with its sub-issues forming a tree.
type Issue struct {
	Number    int
	Title     string
	State     string
	URL       string
	Body      string
	Assignees []string
	Labels    []string
	SubIssues []Issue
}
