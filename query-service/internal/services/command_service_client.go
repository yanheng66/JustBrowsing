package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JustBrowsing/query-service/pkg/errors"
	"go.uber.org/zap"
)

// CommandServiceClient provides a client for the Command Service API
type CommandServiceClient struct {
	logger     *zap.Logger
	baseURL    string
	httpClient *http.Client
}

// NewCommandServiceClient creates a new Command Service client
func NewCommandServiceClient(logger *zap.Logger, baseURL string) *CommandServiceClient {
	return &CommandServiceClient{
		logger:  logger,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendRequest sends a request to the Command Service API
func (c *CommandServiceClient) SendRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var req *http.Request
	var err error

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "failed to marshal request body")
		}

		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return errors.Wrap(err, "failed to create request")
		}

		req.Body = nil // TODO: Set request body
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return errors.Wrap(err, "failed to create request")
		}
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		var errorResponse struct {
			Message string   `json:"message"`
			Details []string `json:"details"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return errors.Wrap(err, fmt.Sprintf("request failed with status %d", resp.StatusCode))
		}

		if len(errorResponse.Details) > 0 {
			return errors.BadRequestWithDetails(errorResponse.Message, errorResponse.Details)
		}
		return errors.New(errorResponse.Message)
	}

	// Decode response body if result is provided
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return errors.Wrap(err, "failed to decode response body")
		}
	}

	return nil
}