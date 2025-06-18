package github

import (
	"context"
	"encoding/json"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"net/http"
	"testing"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v72/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetTeamByName(t *testing.T) {
	mockTeam := &github.Team{
		ID:   github.Ptr(int64(123)),
		Name: github.Ptr("devs"),
		Slug: github.Ptr("devs"),
	}

	tests := []struct {
		name           string
		mockedClient   *http.Client
		requestArgs    map[string]interface{}
		expectError    bool
		expectedTeam   *github.Team
		expectedErrMsg string
	}{
		{
			name: "successful get team by name",
			mockedClient: mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetOrgsTeamsByOrgByTeamSlug,
					mockTeam,
				),
			),
			requestArgs: map[string]interface{}{
				"org":       "myorg",
				"team_slug": "devs",
			},
			expectError:  false,
			expectedTeam: mockTeam,
		},
		{
			name: "team not found",
			mockedClient: mock.NewMockedHTTPClient(
				mock.WithRequestMatchHandler(
					mock.GetOrgsTeamsByOrgByTeamSlug,
					mockResponse(t, http.StatusNotFound, `{"message": "Team not found"}`),
				),
			),
			requestArgs: map[string]interface{}{
				"org":       "myorg",
				"team_slug": "missing",
			},
			expectError:    true,
			expectedErrMsg: "failed to get team",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := github.NewClient(tc.mockedClient)
			_, handler := GetTeamByName(stubGetClientFn(client), translations.NullTranslationHelper)

			request := createMCPRequest(tc.requestArgs)
			result, err := handler(context.Background(), request)

			if tc.expectError {
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg)
				} else {
					assert.NotNil(t, result)
					textContent := getTextResult(t, result)
					assert.Contains(t, textContent.Text, tc.expectedErrMsg)
				}
				return
			}

			require.NoError(t, err)
			textContent := getTextResult(t, result)

			var returnedTeam github.Team
			err = json.Unmarshal([]byte(textContent.Text), &returnedTeam)
			require.NoError(t, err)
			assert.Equal(t, *tc.expectedTeam.ID, *returnedTeam.ID)
			assert.Equal(t, *tc.expectedTeam.Name, *returnedTeam.Name)
			assert.Equal(t, *tc.expectedTeam.Slug, *returnedTeam.Slug)
		})
	}
}

func Test_ListTeams(t *testing.T) {
	mockTeams := []*github.Team{
		{ID: github.Ptr(int64(1)), Name: github.Ptr("devs"), Slug: github.Ptr("devs")},
		{ID: github.Ptr(int64(2)), Name: github.Ptr("ops"), Slug: github.Ptr("ops")},
	}

	tests := []struct {
		name           string
		mockedClient   *http.Client
		requestArgs    map[string]interface{}
		expectError    bool
		expectedTeams  []*github.Team
		expectedErrMsg string
	}{
		{
			name: "successful list teams",
			mockedClient: mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetOrgsTeamsByOrg,
					mockTeams,
				),
			),
			requestArgs: map[string]interface{}{
				"org": "myorg",
			},
			expectError:   false,
			expectedTeams: mockTeams,
		},
		{
			name: "org not found",
			mockedClient: mock.NewMockedHTTPClient(
				mock.WithRequestMatchHandler(
					mock.GetOrgsTeamsByOrg,
					mockResponse(t, http.StatusNotFound, `{"message": "Org not found"}`),
				),
			),
			requestArgs: map[string]interface{}{
				"org": "missing",
			},
			expectError:    true,
			expectedErrMsg: "failed to list teams",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := github.NewClient(tc.mockedClient)
			_, handler := ListTeams(stubGetClientFn(client), translations.NullTranslationHelper)

			request := createMCPRequest(tc.requestArgs)
			result, err := handler(context.Background(), request)

			if tc.expectError {
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg)
				} else {
					assert.NotNil(t, result)
					textContent := getTextResult(t, result)
					assert.Contains(t, textContent.Text, tc.expectedErrMsg)
				}
				return
			}

			require.NoError(t, err)
			textContent := getTextResult(t, result)

			var returnedTeams []*github.Team
			err = json.Unmarshal([]byte(textContent.Text), &returnedTeams)
			require.NoError(t, err)
			assert.Len(t, returnedTeams, len(tc.expectedTeams))
			for i, team := range returnedTeams {
				assert.Equal(t, *tc.expectedTeams[i].ID, *team.ID)
				assert.Equal(t, *tc.expectedTeams[i].Name, *team.Name)
				assert.Equal(t, *tc.expectedTeams[i].Slug, *team.Slug)
			}
		})
	}
}
