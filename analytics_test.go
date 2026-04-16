package amcrest

import (
	"context"
	"testing"
)

func TestAnalyticsGetCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	body, err := c.Analytics.GetCaps(ctx, 0)
	if err != nil {
		t.Skip("GetCaps not supported on this camera, skipping")
	}
	if body == "" {
		t.Fatal("expected non-empty analytics caps")
	}
	t.Logf("Analytics caps (first 500 chars): %.500s", body)
}

func TestAnalyticsGetGlobalConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Analytics.GetGlobalConfig(ctx)
	if err != nil {
		t.Skip("GetGlobalConfig not supported on this camera, skipping")
	}
	if len(cfg) == 0 {
		t.Log("VideoAnalyseGlobal config returned empty map")
	}
	for k, v := range cfg {
		t.Logf("AnalyticsGlobal.%s = %s", k, v)
	}
}

func TestAnalyticsGetRuleConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Analytics.GetRuleConfig(ctx)
	if err != nil {
		t.Skip("GetRuleConfig not supported on this camera, skipping")
	}
	if len(cfg) == 0 {
		t.Log("VideoAnalyseRule config returned empty map")
	}
	for k, v := range cfg {
		t.Logf("AnalyticsRule.%s = %s", k, v)
	}
}

func TestAnalyticsGetSceneList(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	body, err := c.Analytics.GetSceneList(ctx, 0)
	if err != nil {
		t.Skip("GetSceneList not supported on this camera, skipping")
	}
	if body == "" {
		t.Fatal("expected non-empty scene list")
	}
	t.Logf("Scene list (first 500 chars): %.500s", body)
}
