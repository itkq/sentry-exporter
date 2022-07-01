package sentry

import (
	"context"
	"fmt"
)

type EventCountsForOrgV2Response struct {
	Start     string                     `json:"start"`
	End       string                     `json:"end"`
	Intervals []string                   `json:"intervals"`
	Groups    []EventCountsForOrgV2Group `json:"groups"`
}

type EventCountsForOrgV2Group struct {
	By     EventCountsForOrgV2GroupBy `json:"by"`
	Totals map[string]int             `json:"totals"`
	Series map[string][]int           `json:"series"`
}

type EventCountsForOrgV2GroupBy struct {
	Outcome int `json:"outcome"`
	Project int `json:"project"`
}

type RetrieveEventCountsForOrgV2Request struct {
	StatsPeriod string
	Interval    string
	Start       string
	End         string
	GroupBy     []string
	Field       string // TODO: validate
	Project     []string
	Category    string // TODO: enum
	Outcome     string // TODO: enum
	Reason      string
}

// https://docs.sentry.io/api/organizations/retrieve-event-counts-for-an-organization-v2/
func (c *Client) RetrieveEventCountsForOrgV2(ctx context.Context, input *RetrieveEventCountsForOrgV2Request) (*EventCountsForOrgV2Response, error) {
	queries := make(map[string]string, 0)
	queries["field"] = input.Field // required

	arrayQueries := make(map[string][]string, 0)
	arrayQueries["groupBy"] = input.GroupBy // required
	arrayQueries["project"] = input.Project

	params := requestParams{
		method:       "GET",
		subPath:      fmt.Sprintf("organizations/%s/stats_v2/", c.organizationOrg),
		queries:      queries,
		arrayQueries: arrayQueries,
	}

	if input.StatsPeriod != "" {
		queries["statsPeriod"] = input.StatsPeriod
	}
	if input.Interval != "" {
		queries["interval"] = input.Interval
	}
	if input.Start != "" {
		queries["start"] = input.Start
	}
	if input.End != "" {
		queries["end"] = input.End
	}
	if input.Category != "" {
		queries["category"] = input.Category
	}
	if input.Outcome != "" {
		queries["outcome"] = input.Outcome
	}
	if input.Reason != "" {
		queries["reason"] = input.Reason
	}

	var out EventCountsForOrgV2Response
	ctx, cancel := context.WithTimeout(ctx, c.apiTimeout)
	defer cancel()

	if err := c.doAPIRequest(ctx, &params, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
