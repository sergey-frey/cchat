package messageapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/sergey-frey/cchat/server/chat-service/internal/domain/models"
)

var (
	ErrUserExists = fmt.Errorf("user already exists")
	ErrUserNotFound = fmt.Errorf("user not found")
)

type NormalizedMessageResponse struct {
	Status int                `json:"status"`
	Data   models.NormalizedUser `json:"data"`
}

type EmailOfUser struct {
	Email string `json:"email"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	log *slog.Logger
}

func NewClient(httpClient *http.Client, baseURL string, log *slog.Logger) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		log: log,
	}
}

func (c *Client) GetUser(ctx context.Context, email string) (*models.NormalizedUser, error) {
	const op = "api.userapi.client.GetUser"

	log := c.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	param := url.Values{}

	param.Add("email", email)

	log.Info("email:")

	endpoint := fmt.Sprintf("http://users-service-1:8080/users/?%s", param.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: unexpected status code %d", op, resp.StatusCode)
	}

	defer resp.Body.Close()

	var user NormalizedMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user.Data, nil
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

	endpoint := fmt.Sprintf("http://users-service-1:8080/users/")

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

	var newUser NormalizedMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &newUser.Data, nil
}