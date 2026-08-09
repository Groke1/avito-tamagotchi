package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxErrorBodySize = 4 << 10

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http status %d: %s", e.StatusCode, e.Body)
}

func doJSON[Resp any](ctx context.Context, client *http.Client,
	method string, endpoint string, body io.Reader) (*Resp, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))

		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(data)),
		}
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil //nolint:nilnil // A successful 204 response intentionally has no response value
	}

	var result Resp

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

func WithoutRequest[Resp any](ctx context.Context, client *http.Client, endpoint, method string) (*Resp, error) {
	return doJSON[Resp](ctx, client, method, endpoint, nil)
}

func WithRequest[Req any, Resp any](
	ctx context.Context, client *http.Client,
	endpoint string, method string, request Req) (*Resp, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	return doJSON[Resp](ctx, client, method, endpoint, bytes.NewReader(data))
}

func NormalizeAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	addr = strings.TrimRight(addr, "/")

	if addr == "" {
		return ""
	}

	if strings.HasPrefix(addr, "http://") ||
		strings.HasPrefix(addr, "https://") {
		return addr
	}

	return "http://" + addr
}
