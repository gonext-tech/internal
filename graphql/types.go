package graphql

import "github.com/shurcooL/githubv4"

// ProjectV2 represents a GitHub ProjectV2
type ProjectV2 struct {
	ID     githubv4.String
	Title  githubv4.String
	Number githubv4.Int
	URL    githubv4.String
}

// ProjectItem represents an item in a GitHub ProjectV2
type ProjectItem struct {
	ID      githubv4.String
	Content struct {
		TypeName    githubv4.String `graphql:"__typename"`
		Issue       IssueContent
		DraftIssue  DraftIssueContent
		PullRequest PullRequestContent
	}
	FieldValues struct {
		Nodes []FieldValue
	} `graphql:"fieldValues(first: 10)"`
}

// FieldValue represents custom field values
type FieldValue struct {
	TypeName  githubv4.String `graphql:"__typename"`
	TextValue struct {
		Text githubv4.String
	} `graphql:"... on ProjectV2ItemFieldTextValue"`
	DateValue struct {
		Date githubv4.String
	} `graphql:"... on ProjectV2ItemFieldDateValue"`
	SingleSelectValue struct {
		Name githubv4.String
	} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
}

// IssueContent represents issue content
type IssueContent struct {
	Title githubv4.String
	URL   githubv4.String
}

// DraftIssueContent represents draft issue content
type DraftIssueContent struct {
	Title githubv4.String
	Body  githubv4.String
}

// PullRequestContent represents PR content
type PullRequestContent struct {
	Title githubv4.String
	URL   githubv4.String
}
