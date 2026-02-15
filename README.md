# gh-issue-tree

A GitHub CLI extension that recursively crawls sub-issues and generates a structured Markdown tree.

When managing tasks with GitHub Issues and sub-issues, it can be hard to see the full picture. `gh issue-tree` takes an issue number, recursively fetches all sub-issues via the GitHub GraphQL API, and outputs a Markdown document showing the complete hierarchy with issue bodies.

## Installation

```sh
gh extension install kspace-trk/gh-issue-tree
```

## Usage

```sh
gh issue-tree <issue-number>
```

Run this from inside a Git repository that has a GitHub remote. The owner/repo is automatically detected.

### Save to file

```sh
gh issue-tree 42 > issues.md
```

## Output example

```markdown
# Tree

- [ ] #42 User authentication system
├── [x] #43 Design auth database schema
├── [ ] #44 Implement login flow
│   ├── [x] #46 Create login API endpoint
│   └── [ ] #47 Build login form UI
└── [ ] #45 Implement OAuth integration

---

# #42 User authentication system
**Status:** open | **Labels:** epic | **Sub-issues:** #43, #44, #45

Build a complete authentication system supporting email/password and OAuth.

## #43 Design auth database schema
**Status:** closed | **Parent:** #42 | **Assignees:** @alice

Define tables for users, sessions, and OAuth tokens.

---

## #44 Implement login flow
**Status:** open | **Parent:** #42 | **Sub-issues:** #46, #47

### #46 Create login API endpoint
**Status:** closed | **Parent:** #44 | **Assignees:** @bob

POST /api/login - validate credentials and return JWT.

---

### #47 Build login form UI
**Status:** open | **Parent:** #44

Create a responsive login form with validation.

---

## #45 Implement OAuth integration
**Status:** open | **Parent:** #42 | **Labels:** feature

Add Google and GitHub OAuth providers.
```

Each issue section includes:

- **Status** (open/closed, shown as checkboxes in the tree)
- **Assignees** and **Labels** (when present)
- **Parent** / **Sub-issues** references for navigation
- Full issue body

## Requirements

- [GitHub CLI](https://cli.github.com/) (`gh`) v2.0+
- A GitHub repository with sub-issues

## License

[MIT](LICENSE)
