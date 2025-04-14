package handlers

import (
	"net/http"
	"your-app/graphql"

	"github.com/labstack/echo/v4"
	"github.com/shurcooL/githubv4"
)

// GetAllTasks returns all tasks in a project
func GetAllTasks(c echo.Context) error {
	projectID := c.Param("project_id")
	client := getGitHubV4Client(c)

	var q graphql.GetProjectItemsQuery
	variables := map[string]interface{}{
		"projectId": githubv4.ID(projectID),
		"cursor":    (*githubv4.String)(nil),
	}

	err := client.Query(c.Request().Context(), &q, variables)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert to simplified task structure
	var tasks []map[string]interface{}
	for _, item := range q.Node.ProjectV2.Items.Nodes {
		task := map[string]interface{}{
			"id":    item.ID,
			"title": getItemTitle(item),
			"type":  item.Content.TypeName,
		}
		tasks = append(tasks, task)
	}

	return c.JSON(http.StatusOK, tasks)
}

// GetSingleTask returns a single task
func GetSingleTask(c echo.Context) error {
	itemID := c.Param("item_id")
	client := getGitHubV4Client(c)

	var q graphql.GetSingleItemQuery
	variables := map[string]interface{}{
		"itemId": githubv4.ID(itemID),
	}

	err := client.Query(c.Request().Context(), &q, variables)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, q.Node.ProjectItem)
}

// CreateTask creates a new draft issue task
func CreateTask(c echo.Context) error {
	projectID := c.Param("project_id")
	client := getGitHubV4Client(c)

	type TaskRequest struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	var req TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	var m graphql.AddDraftIssueMutation
	input := githubv4.AddProjectV2DraftIssueInput{
		ProjectID: githubv4.ID(projectID),
		Title:     githubv4.String(req.Title),
		Body:      githubv4.String(req.Body),
	}

	err := client.Mutate(c.Request().Context(), &m, input, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"id":    string(m.AddProjectV2DraftIssue.ProjectItem.ID),
		"title": req.Title,
	})
}

// DeleteTask deletes a project item
func DeleteTask(c echo.Context) error {
	itemID := c.Param("item_id")
	client := getGitHubV4Client(c)

	var m graphql.DeleteItemMutation
	input := githubv4.DeleteProjectV2ItemInput{
		ItemID: githubv4.ID(itemID),
	}

	err := client.Mutate(c.Request().Context(), &m, input, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"deletedItemId": string(m.DeleteProjectV2Item.DeletedItemID),
	})
}

// Helper function to get item title based on type
func getItemTitle(item graphql.ProjectItem) string {
	switch item.Content.TypeName {
	case "Issue":
		return string(item.Content.Issue.Title)
	case "DraftIssue":
		return string(item.Content.DraftIssue.Title)
	case "PullRequest":
		return string(item.Content.PullRequest.Title)
	default:
		return "Untitled"
	}
}
