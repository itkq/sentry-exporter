package sentry

import (
	"context"
	"fmt"
)

type Project struct {
	Id   string `json:"id"`
	Slug string `json:"slug"`
}

// https://docs.sentry.io/api/organizations/list-an-organizations-projects/
func (c *Client) ListOrganizationProjects(ctx context.Context) ([]Project, error) {
	params := requestParams{
		method:  "GET",
		subPath: fmt.Sprintf("organizations/%s/projects/", c.organizationOrg),
	}

	var out []Project
	ctx, cancel := context.WithTimeout(ctx, c.apiTimeout)
	defer cancel()

	if err := c.doAPIRequest(ctx, &params, &out); err != nil {
		return nil, err
	}

	return out, nil
}
