package services

import (
	"context"
	"your-app/pkg/auth"
	"your-app/pkg/graphql"

	"github.com/shurcooL/githubv4"
)

type TaskService struct {
	client *auth.GitHubClient
}

func NewTaskService(client *auth.GitHubClient) *TaskService {
	return &TaskService{client: client}
}

func (s *TaskService) GetAllTasks(ctx context.Context, projectID string) ([]*graphql.ProjectItem, error) {
	var q graphql.GetProjectItemsQuery
	variables := map[string]interface{}{
		"projectId": githubv4.ID(projectID),
		"cursor":    (*githubv4.String)(nil),
	}

	err := s.client.Query(ctx, &q, variables)
	if err != nil {
		return nil, err
	}

	return q.Node.ProjectV2.Items.Nodes, nil
}

func (s *TaskService) GetTask(ctx context.Context, itemID string) (*graphql.ProjectItem, error) {
	var q graphql.GetSingleItemQuery
	variables := map[string]interface{}{
		"itemId": githubv4.ID(itemID),
	}

	err := s.client.Query(ctx, &q, variables)
	if err != nil {
		return nil, err
	}

	return &q.Node.ProjectItem.ProjectItem, nil
}

func (s *TaskService) CreateTask(ctx context.Context, projectID, title, body string) (string, error) {
	var m graphql.AddDraftIssueMutation
	input := githubv4.AddProjectV2DraftIssueInput{
		ProjectID: githubv4.ID(projectID),
		Title:     githubv4.String(title),
		Body:      githubv4.String(body),
	}

	err := s.client.Mutate(ctx, &m, input, nil)
	if err != nil {
		return "", err
	}

	return string(m.AddProjectV2DraftIssue.ProjectItem.ID), nil
}

func (s *TaskService) DeleteTask(ctx context.Context, itemID string) error {
	var m graphql.DeleteItemMutation
	input := githubv4.DeleteProjectV2ItemInput{
		ItemID: githubv4.ID(itemID),
	}

	return s.client.Mutate(ctx, &m, input, nil)
}
