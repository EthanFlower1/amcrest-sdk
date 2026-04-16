package amcrest

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func loadEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			if os.Getenv(parts[0]) == "" {
				os.Setenv(parts[0], parts[1])
			}
		}
	}
}

func testClient(t *testing.T) *Client {
	t.Helper()
	loadEnv()
	host := os.Getenv("AMCREST_HOST")
	user := os.Getenv("AMCREST_USERNAME")
	pass := os.Getenv("AMCREST_PASSWORD")
	if host == "" || user == "" || pass == "" {
		t.Skip("AMCREST_HOST, AMCREST_USERNAME, AMCREST_PASSWORD must be set")
	}
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}
