package sentry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	apiKey          string
	apiEndpoint     string
	organizationOrg string
	apiTimeout      time.Duration
}

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

func NewClient(apiKey string, apiEndpoint string, organizationSlug string) *Client {
	return &Client{
		apiKey:          apiKey,
		apiEndpoint:     apiEndpoint,
		organizationOrg: organizationSlug,
		apiTimeout:      5 * time.Second,
	}
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
	ctx, cancel := context.WithTimeout(ctx, c.apiTimeout)
	defer cancel()

	baseUrl, err := url.Parse(c.apiEndpoint)
	if err != nil {
		return nil, err
	}
	// XXX: tailing slash is required
	baseUrl.Path = baseUrl.Path + fmt.Sprintf("organizations/%s/stats_v2/", c.organizationOrg)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()

	// required from here
	for _, g := range input.GroupBy {
		query.Add("groupBy", g)
	}
	query.Add("field", input.Field)

	// non-required from here
	if input.StatsPeriod != "" {
		query.Add("statsPeriod", input.StatsPeriod)
	}
	if input.Interval != "" {
		query.Add("interval", input.Interval)
	}
	if input.Start != "" {
		query.Add("start", input.Start)
	}
	if input.End != "" {
		query.Add("end", input.End)
	}
	for _, p := range input.Project {
		query.Add("project", p)
	}
	if input.Category != "" {
		query.Add("category", input.Category)
	}
	if input.Outcome != "" {
		query.Add("outcome", input.Outcome)
	}
	if input.Reason != "" {
		query.Add("reason", input.Reason)
	}

	req.URL.RawQuery = query.Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("error body: %s\n", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var out EventCountsForOrgV2Response
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
