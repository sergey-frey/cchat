package userapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sergey-frey/cchat/server/auth-service/internal/domain/models"
)

var (
	ErrUserExists = fmt.Errorf("user already exists")
	ErrUserNotFound = fmt.Errorf("user not found")
)

type CreateUserResponse struct {
	Status int                `json:"status"`
	Data   models.NormalizedUser `json:"data"`
}

type EmailOfUser struct {
	Email string `json:"email"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *Client) GetUser(ctx context.Context, email string) (*models.NormalizedUser, error) {
	const op = "api.userapi.client.GetUser"

	endpoint := fmt.Sprintf("/users/%s", email)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: unexpected status code %d", op, resp.StatusCode)
	}

	defer resp.Body.Close()

	var newUser models.NormalizedUser
	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &newUser, nil
}

func (c *Client) CreateUser(ctx context.Context, email string) (*models.NormalizedUser, error) {
	const op = "api.userapi.client.CreateUser"

	requestUser := EmailOfUser{
		Email: email,
	}

	requestBody, err := json.Marshal(requestUser)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	endpoint := fmt.Sprintf("http://user-service:3040/users/")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusConflict {
			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return nil, fmt.Errorf("%s: unexpected status code %d", op, resp.StatusCode)
	}

	defer resp.Body.Close()

	var newUser CreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &newUser.Data, nil
}