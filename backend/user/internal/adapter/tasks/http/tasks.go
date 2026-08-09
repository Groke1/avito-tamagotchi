package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/adapter/httpclient"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewTasksClient(addr string, httpClient *http.Client) (*client, error) {
	addr = httpclient.NormalizeAddr(addr)
	if addr == "" {
		return nil, errors.New("tasks service address is empty")
	}

	baseURL, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse tasks service address: %w", err)
	}

	return &client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

func (c *client) GetCompletedTasks(ctx context.Context, userID string) ([]entity.TasksStat, error) {
	endpoint := c.baseURL.ResolveReference(
		&url.URL{
			Path: "/internal/tasks/" + userID,
		},
	)

	resp, err := httpclient.WithoutRequest[getCompletedTasksResponse](
		ctx, c.httpClient, endpoint.String(), http.MethodGet)
	if err != nil {
		return nil, fmt.Errorf("get tasks stat: %w", err)
	}

	return resp.Items, nil
}
