package sentry

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	apiKey          string
	apiEndpoint     string
	organizationOrg string
	apiTimeout      time.Duration
	httpClient      *http.Client
}

func NewDefaultClient(apiKey string, apiEndpoint string, organizationSlug string) *Client {
	return &Client{
		apiKey:          apiKey,
		apiEndpoint:     apiEndpoint,
		organizationOrg: organizationSlug,
		apiTimeout:      5 * time.Second,
		httpClient:      &http.Client{},
	}
}

type requestParams struct {
	method       string
	subPath      string
	queries      map[string]string
	arrayQueries map[string][]string
}

func (c *Client) doAPIRequest(ctx context.Context, params *requestParams, out interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.apiTimeout)

	defer cancel()

	req, err := c.newRequest(ctx, params)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, out)

	return err
}

func (c *Client) newRequest(ctx context.Context, params *requestParams) (*http.Request, error) {
	url, err := url.Parse(c.apiEndpoint)
	if err != nil {
		return nil, err
	}

	url.Path = path.Join(url.Path, params.subPath) + "/" // XXX: trailing slash is required

	req, err := http.NewRequestWithContext(ctx, params.method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()

	for k, v := range params.queries {
		query.Add(k, v)
	}

	for k, vs := range params.arrayQueries {
		for _, v := range vs {
			query.Add(k, v)
		}
	}

	req.URL.RawQuery = query.Encode()

	return req, nil
}
