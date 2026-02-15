package main

import (
	"fmt"
	"strings"
)

func RenderMarkdown(issue Issue) string {
	var sb strings.Builder

	// Tree overview
	sb.WriteString("# Tree\n\n")
	renderTreeOverview(&sb, issue, "", true)
	sb.WriteString("\n---\n\n")

	// Detailed sections
	renderIssue(&sb, issue, 0, 1)
	return sb.String()
}

func stateCheckbox(state string) string {
	if strings.EqualFold(state, "closed") {
		return "[x]"
	}
	return "[ ]"
}

func renderTreeOverview(sb *strings.Builder, issue Issue, prefix string, isRoot bool) {
	if isRoot {
		fmt.Fprintf(sb, "- %s #%d %s\n", stateCheckbox(issue.State), issue.Number, issue.Title)
	}
	for i, sub := range issue.SubIssues {
		isLast := i == len(issue.SubIssues)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		fmt.Fprintf(sb, "%s%s%s #%d %s\n", prefix, connector, stateCheckbox(sub.State), sub.Number, sub.Title)

		childPrefix := prefix + "│   "
		if isLast {
			childPrefix = prefix + "    "
		}
		renderTreeOverview(sb, sub, childPrefix, false)
	}
}

func renderIssue(sb *strings.Builder, issue Issue, parentNumber int, depth int) {
	// Heading
	hashes := strings.Repeat("#", depth)
	fmt.Fprintf(sb, "%s #%d %s\n", hashes, issue.Number, issue.Title)

	// Metadata line
	var meta []string
	meta = append(meta, fmt.Sprintf("**Status:** %s", strings.ToLower(issue.State)))
	if len(issue.Assignees) > 0 {
		mentions := make([]string, len(issue.Assignees))
		for i, a := range issue.Assignees {
			mentions[i] = "@" + a
		}
		meta = append(meta, fmt.Sprintf("**Assignees:** %s", strings.Join(mentions, ", ")))
	}
	if len(issue.Labels) > 0 {
		meta = append(meta, fmt.Sprintf("**Labels:** %s", strings.Join(issue.Labels, ", ")))
	}
	if parentNumber > 0 {
		meta = append(meta, fmt.Sprintf("**Parent:** #%d", parentNumber))
	}
	if len(issue.SubIssues) > 0 {
		refs := make([]string, len(issue.SubIssues))
		for i, sub := range issue.SubIssues {
			refs[i] = fmt.Sprintf("#%d", sub.Number)
		}
		meta = append(meta, fmt.Sprintf("**Sub-issues:** %s", strings.Join(refs, ", ")))
	}
	fmt.Fprintf(sb, "%s\n", strings.Join(meta, " | "))

	// Body
	body := strings.TrimSpace(issue.Body)
	if body != "" {
		fmt.Fprintf(sb, "\n%s\n", body)
	}

	// Sub-issues
	for i, sub := range issue.SubIssues {
		sb.WriteString("\n")
		if i > 0 {
			sb.WriteString("---\n\n")
		}
		renderIssue(sb, sub, issue.Number, depth+1)
	}
}
