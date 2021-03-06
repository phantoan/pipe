// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"

	"github.com/pipe-cd/pipe/pkg/model"
)

// OAuthClient is a oauth client for github.
type OAuthClient struct {
	*github.Client

	projectID  string
	adminTeam  string
	editorTeam string
	viewerTeam string
}

// NewOAuthClient creates a new oauth client for GitHub.
func NewOAuthClient(ctx context.Context,
	sso *model.ProjectSSOConfig_GitHub,
	rbac *model.ProjectRBACConfig,
	projectID, code string,
	enterprise bool,
) (*OAuthClient, error) {
	c := &OAuthClient{
		projectID:  projectID,
		adminTeam:  rbac.Admin,
		editorTeam: rbac.Editor,
		viewerTeam: rbac.Viewer,
	}
	cfg := oauth2.Config{
		ClientID:     sso.ClientId,
		ClientSecret: sso.ClientSecret,
		Endpoint:     oauth2github.Endpoint,
	}
	if enterprise {
		return newGHEOAuthClient(ctx, sso, code, c, cfg)
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	c.Client = github.NewClient(cfg.Client(ctx, token))
	return c, nil
}

func newGHEOAuthClient(ctx context.Context,
	sso *model.ProjectSSOConfig_GitHub,
	code string,
	c *OAuthClient,
	cfg oauth2.Config,
) (*OAuthClient, error) {
	baseURL, err := url.Parse(sso.BaseUrl)
	if err != nil {
		return nil, err
	}
	cfg.Endpoint.TokenURL = fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, "/login/oauth/access_token")
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	cli, err := github.NewEnterpriseClient(sso.BaseUrl, sso.UploadUrl, cfg.Client(ctx, token))
	if err != nil {
		return nil, err
	}
	c.Client = cli
	return c, nil
}

// GetUser returns a user model.
func (c *OAuthClient) GetUser(ctx context.Context) (*model.User, error) {
	user, _, err := c.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	teams, _, err := c.Teams.ListUserTeams(ctx, nil)
	if err != nil {
		return nil, err
	}
	role, err := c.decideRole(user.GetLogin(), teams)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Username:  user.GetLogin(),
		AvatarUrl: user.GetAvatarURL(),
		Role: &model.Role{
			ProjectId:   c.projectID,
			ProjectRole: *role,
		},
	}, nil
}

func (c *OAuthClient) decideRole(user string, teams []*github.Team) (*model.Role_ProjectRole, error) {
	var viewer, editor bool
	for _, team := range teams {
		slug := team.GetSlug()
		org := team.Organization.GetLogin()
		if org == "" || slug == "" {
			continue
		}
		switch fmt.Sprintf("%s/%s", org, slug) {
		case c.adminTeam:
			r := model.Role_ADMIN
			return &r, nil
		case c.editorTeam:
			editor = true
		case c.viewerTeam:
			viewer = true
		}
	}
	if editor {
		r := model.Role_EDITOR
		return &r, nil
	}
	if viewer {
		r := model.Role_VIEWER
		return &r, nil
	}
	return nil, fmt.Errorf("user (%s) not found in any of the %d project teams", user, len(teams))
}
