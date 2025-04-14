package graphql

import "github.com/shurcooL/githubv4"

// AddDraftIssueMutation creates a new draft issue (task)
type AddDraftIssueMutation struct {
	AddProjectV2DraftIssue struct {
		ProjectItem struct {
			ID githubv4.String
		}
	} `graphql:"addProjectV2DraftIssue(input: $input)"`
}

// DeleteItemMutation deletes a project item
type DeleteItemMutation struct {
	DeleteProjectV2Item struct {
		DeletedItemID githubv4.String
	} `graphql:"deleteProjectV2Item(input: $input)"`
}
