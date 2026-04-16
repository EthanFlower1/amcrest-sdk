package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// getConfig retrieves a named config table and strips the "table.<name>." prefix from keys.
func (c *Client) getConfig(ctx context.Context, name string) (map[string]string, error) {
	body, err := c.cgiGet(ctx, "configManager.cgi", "getConfig", url.Values{
		"name": {name},
	})
	if err != nil {
		return nil, err
	}
	prefix := fmt.Sprintf("table.%s.", name)
	return parseKVWithPrefix(body, prefix), nil
}

// getConfigIndexed retrieves an indexed config table (e.g., "VideoColor[0]").
func (c *Client) getConfigIndexed(ctx context.Context, name string, index int) (map[string]string, error) {
	indexedName := fmt.Sprintf("%s[%d]", name, index)
	body, err := c.cgiGet(ctx, "configManager.cgi", "getConfig", url.Values{
		"name": {indexedName},
	})
	if err != nil {
		return nil, err
	}
	prefix := fmt.Sprintf("table.%s.", indexedName)
	return parseKVWithPrefix(body, prefix), nil
}

// getRawConfig retrieves a named config table without stripping any prefix.
func (c *Client) getRawConfig(ctx context.Context, name string) (map[string]string, error) {
	body, err := c.cgiGet(ctx, "configManager.cgi", "getConfig", url.Values{
		"name": {name},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// setConfig sets config values via configManager.cgi setConfig action.
func (c *Client) setConfig(ctx context.Context, params map[string]string) error {
	qv := url.Values{}
	for k, v := range params {
		qv.Set(k, v)
	}
	return c.cgiAction(ctx, "configManager.cgi", "setConfig", qv)
}
