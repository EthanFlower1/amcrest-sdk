package amcrest

import (
	"context"
	"testing"
)

func TestAnalytics(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasAnalytics, "Video Analytics")

	t.Run("GetCaps", func(t *testing.T) {
		body, err := c.Analytics.GetCaps(ctx, 0)
		if err != nil {
			t.Fatalf("GetCaps: %v", err)
		}
		if body == "" {
			t.Fatal("expected non-empty analytics caps")
		}
		t.Logf("Analytics caps (first 500 chars): %.500s", body)
	})

	t.Run("GetGlobalConfig", func(t *testing.T) {
		cfg, err := c.Analytics.GetGlobalConfig(ctx)
		if err != nil {
			t.Fatalf("GetGlobalConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("AnalyticsGlobal.%s = %s", k, v)
		}
	})

	t.Run("GetRuleConfig", func(t *testing.T) {
		cfg, err := c.Analytics.GetRuleConfig(ctx)
		if err != nil {
			t.Fatalf("GetRuleConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("AnalyticsRule.%s = %s", k, v)
		}
	})

	t.Run("GetSceneList", func(t *testing.T) {
		body, err := c.Analytics.GetSceneList(ctx, 0)
		if err != nil {
			t.Fatalf("GetSceneList: %v", err)
		}
		t.Logf("Scene list (first 500 chars): %.500s", body)
	})
}
