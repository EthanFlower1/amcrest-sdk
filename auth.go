package amcrest

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// digestTransport implements http.RoundTripper with HTTP Digest Authentication.
type digestTransport struct {
	username  string
	password  string
	transport http.RoundTripper

	mu    sync.Mutex
	realm string
	nonce string
	qop   string
	nc    int
}

func newDigestTransport(username, password string, transport http.RoundTripper) *digestTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &digestTransport{
		username:  username,
		password:  password,
		transport: transport,
	}
}

func (t *digestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Try with cached credentials first
	t.mu.Lock()
	if t.nonce != "" {
		t.nc++
		authHeader := t.buildAuthHeader(req.Method, req.URL.RequestURI())
		req = req.Clone(req.Context())
		req.Header.Set("Authorization", authHeader)
		t.mu.Unlock()
		resp, err := t.transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}
		resp.Body.Close()
	} else {
		t.mu.Unlock()
	}

	// Send initial request without auth
	initialReq := req.Clone(req.Context())
	initialReq.Header.Del("Authorization")

	resp, err := t.transport.RoundTrip(initialReq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	challenge := resp.Header.Get("WWW-Authenticate")
	resp.Body.Close()

	if challenge == "" {
		return nil, fmt.Errorf("amcrest: 401 response without WWW-Authenticate header")
	}

	t.mu.Lock()
	t.parseChallenge(challenge)
	t.nc = 1
	authHeader := t.buildAuthHeader(req.Method, req.URL.RequestURI())
	t.mu.Unlock()

	authReq := req.Clone(req.Context())
	authReq.Header.Set("Authorization", authHeader)

	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("amcrest: failed to reset request body: %w", err)
		}
		authReq.Body = body
	}

	return t.transport.RoundTrip(authReq)
}

func (t *digestTransport) parseChallenge(challenge string) {
	challenge = strings.TrimPrefix(challenge, "Digest ")
	parts := splitChallenge(challenge)
	for k, v := range parts {
		switch strings.ToLower(k) {
		case "realm":
			t.realm = v
		case "nonce":
			t.nonce = v
		case "qop":
			t.qop = v
		}
	}
}

func splitChallenge(challenge string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(challenge, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		eqIdx := strings.Index(part, "=")
		if eqIdx < 0 {
			continue
		}
		key := strings.TrimSpace(part[:eqIdx])
		val := strings.TrimSpace(part[eqIdx+1:])
		val = strings.Trim(val, `"`)
		result[key] = val
	}
	return result
}

func (t *digestTransport) buildAuthHeader(method, uri string) string {
	ha1 := md5Hash(fmt.Sprintf("%s:%s:%s", t.username, t.realm, t.password))
	ha2 := md5Hash(fmt.Sprintf("%s:%s", method, uri))
	nc := fmt.Sprintf("%08x", t.nc)
	cnonce := md5Hash(fmt.Sprintf("%d", t.nc))

	var response string
	if t.qop == "auth" || t.qop == "auth-int" {
		response = md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, t.nonce, nc, cnonce, t.qop, ha2))
	} else {
		response = md5Hash(fmt.Sprintf("%s:%s:%s", ha1, t.nonce, ha2))
	}

	parts := []string{
		fmt.Sprintf(`username="%s"`, t.username),
		fmt.Sprintf(`realm="%s"`, t.realm),
		fmt.Sprintf(`nonce="%s"`, t.nonce),
		fmt.Sprintf(`uri="%s"`, uri),
		fmt.Sprintf(`response="%s"`, response),
	}
	if t.qop != "" {
		parts = append(parts,
			fmt.Sprintf(`qop=%s`, t.qop),
			fmt.Sprintf(`nc=%s`, nc),
			fmt.Sprintf(`cnonce="%s"`, cnonce),
		)
	}
	return "Digest " + strings.Join(parts, ", ")
}

func md5Hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
