package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v72/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ListTeams creates a tool to list all teams in an organization.
func ListTeams(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_teams",
			mcp.WithDescription(t("TOOL_LIST_TEAMS_DESCRIPTION", "List all teams in a GitHub organization.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_LIST_TEAMS_USER_TITLE", "List teams"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("The organization name"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			org, err := RequiredParam[string](request, "org")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			opts := &github.ListOptions{
				PerPage: pagination.perPage,
				Page:    pagination.page,
			}
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			teams, resp, err := client.Teams.ListTeams(ctx, org, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list teams: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list teams: %s", string(body))), nil
			}
			r, err := json.Marshal(teams)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal teams: %w", err)
			}
			return mcp.NewToolResultText(string(r)), nil
		}
}

// GetTeamByName creates a tool to get a team by its slug (name) in an organization.
func GetTeamByName(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_team_by_name",
			mcp.WithDescription(t("TOOL_GET_TEAM_BY_NAME_DESCRIPTION", "Get details of a team in a GitHub organization by its name (slug).")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_TEAM_BY_NAME_USER_TITLE", "Get team by name"),
				ReadOnlyHint: ToBoolPtr(true),
			}),
			mcp.WithString("org",
				mcp.Required(),
				mcp.Description("The organization name"),
			),
			mcp.WithString("team_slug",
				mcp.Required(),
				mcp.Description("The team slug (name)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			org, err := RequiredParam[string](request, "org")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			teamSlug, err := RequiredParam[string](request, "team_slug")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			team, resp, err := client.Teams.GetTeamBySlug(ctx, org, teamSlug)
			if err != nil {
				return nil, fmt.Errorf("failed to get team: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get team: %s", string(body))), nil
			}
			r, err := json.Marshal(team)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal team: %w", err)
			}
			return mcp.NewToolResultText(string(r)), nil
		}
}
