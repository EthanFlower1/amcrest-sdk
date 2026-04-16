package amcrest

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func parseKV(body string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		result[key] = val
	}
	return result
}

func parseKVWithPrefix(body, prefix string) map[string]string {
	raw := parseKV(body)
	result := make(map[string]string)
	for k, v := range raw {
		if strings.HasPrefix(k, prefix) {
			result[strings.TrimPrefix(k, prefix)] = v
		} else {
			result[k] = v
		}
	}
	return result
}

func readBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("amcrest: reading response body: %w", err)
	}
	body := string(data)

	if resp.StatusCode >= 400 {
		return "", &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(body),
		}
	}

	trimmed := strings.TrimSpace(body)
	if strings.HasPrefix(trimmed, "Error") {
		return "", &APIError{
			StatusCode: resp.StatusCode,
			Message:    trimmed,
		}
	}

	return body, nil
}

func checkOK(resp *http.Response) error {
	body, err := readBody(resp)
	if err != nil {
		return err
	}
	if strings.TrimSpace(body) != "OK" {
		return fmt.Errorf("amcrest: expected OK, got: %s", body)
	}
	return nil
}
