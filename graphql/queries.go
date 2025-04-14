package graphql

import "github.com/shurcooL/githubv4"

// GetProjectsQuery gets all projects for the authenticated user
type GetProjectsQuery struct {
	Viewer struct {
		ProjectsV2 struct {
			Nodes    []ProjectV2
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
		} `graphql:"projectsV2(first: 100, after: $cursor)"`
	}
}

// GetProjectItemsQuery gets all items in a project
type GetProjectItemsQuery struct {
	Node struct {
		ProjectV2 struct {
			Items struct {
				Nodes    []ProjectItem
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
			} `graphql:"items(first: 100, after: $cursor)"`
		} `graphql:"... on ProjectV2"`
	} `graphql:"node(id: $projectId)"`
}

// GetSingleItemQuery gets a single project item
type GetSingleItemQuery struct {
	Node struct {
		ProjectItem struct {
			ProjectItem
		} `graphql:"... on ProjectV2Item"`
	} `graphql:"node(id: $itemId)"`
}
