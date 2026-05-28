package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cours-go/internal/domain"
)

type HTTPError struct {
	StatusCode int
	APIError   domain.APIError
}

func (e *HTTPError) Error() string {
	if e.APIError.Message != "" {
		return e.APIError.Message
	}

	return fmt.Sprintf("http error: %d", e.StatusCode)
}

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = strings.TrimRight(baseURL, "/")
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) Login(ctx context.Context, email, password string) (domain.LoginResponse, error) {
	return doJSON[domain.LoginResponse](ctx, c, http.MethodPost, "/auth/login", domain.LoginRequest{
		Email:    email,
		Password: password,
	})
}

func (c *Client) Me(ctx context.Context) (domain.User, error) {
	return doJSON[domain.User](ctx, c, http.MethodGet, "/auth/me", nil)
}

func (c *Client) GetRooms(ctx context.Context) ([]domain.Room, error) {
	return doJSON[[]domain.Room](ctx, c, http.MethodGet, "/rooms", nil)
}

func (c *Client) GetAvailability(ctx context.Context, roomID int64, day string) ([]domain.Slot, error) {
	path := fmt.Sprintf("/rooms/%d/availability?date=%s", roomID, url.QueryEscape(day))
	return doJSON[[]domain.Slot](ctx, c, http.MethodGet, path, nil)
}

func (c *Client) GetMyReservations(ctx context.Context) ([]domain.Reservation, error) {
	return doJSON[[]domain.Reservation](ctx, c, http.MethodGet, "/reservations/me", nil)
}

func (c *Client) CreateReservation(ctx context.Context, req domain.CreateReservationRequest) (domain.Reservation, error) {
	return doJSON[domain.Reservation](ctx, c, http.MethodPost, "/reservations", req)
}

func (c *Client) CancelReservation(ctx context.Context, reservationID int64) error {
	return doNoContent(ctx, c, http.MethodDelete, fmt.Sprintf("/reservations/%d", reservationID), nil)
}

func (c *Client) ImportConfig(ctx context.Context) (domain.ClientConfig, error) {
	return doJSON[domain.ClientConfig](ctx, c, http.MethodGet, "/me/config", nil)
}

func (c *Client) ExportConfig(ctx context.Context, cfg domain.ClientConfig) error {
	return doNoContent(ctx, c, http.MethodPut, "/me/config", cfg)
}

func (c *Client) ImportState(ctx context.Context) (domain.ClientState, error) {
	return doJSON[domain.ClientState](ctx, c, http.MethodGet, "/me/state", nil)
}

func (c *Client) ExportState(ctx context.Context, state domain.ClientState) error {
	return doNoContent(ctx, c, http.MethodPut, "/me/state", state)
}

func doJSON[T any](ctx context.Context, client *Client, method, path string, payload any) (T, error) {
	var zero T

	response, err := client.request(ctx, method, path, payload)
	if err != nil {
		return zero, err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return zero, decodeHTTPError(response)
	}

	if err := json.NewDecoder(response.Body).Decode(&zero); err != nil {
		return zero, fmt.Errorf("decode response: %w", err)
	}

	return zero, nil
}

func doNoContent(ctx context.Context, client *Client, method, path string, payload any) error {
	response, err := client.request(ctx, method, path, payload)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return decodeHTTPError(response)
	}

	return nil
}

func (c *Client) request(ctx context.Context, method, path string, payload any) (*http.Response, error) {
	if c.baseURL == "" {
		return nil, errors.New("api url is empty")
	}

	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform request: %w", err)
	}

	return resp, nil
}

func decodeHTTPError(response *http.Response) error {
	var payload domain.APIError
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return &HTTPError{
			StatusCode: response.StatusCode,
			APIError: domain.APIError{
				Code:    "http_error",
				Message: response.Status,
			},
		}
	}

	return &HTTPError{
		StatusCode: response.StatusCode,
		APIError:   payload,
	}
}
