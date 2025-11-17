package userapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/sergey-frey/cchat/server/chat-service/internal/domain/models"
)

var (
	ErrUserExists = fmt.Errorf("user already exists")
	ErrUserNotFound = fmt.Errorf("user not found")
)

type NormalizedUserResponse struct {
	Status int                `json:"status"`
	Data   models.NormalizedUser `json:"data"`
}

type ChecUsersRequest struct {
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

func (c *Client) CheckUser(ctx context.Context, ids []uuid.UUID) error {
	const op = "api.userapi.client.CheckUser"

	log := c.log.With(
		slog.String("op", op),
	)

	param := url.Values{}
	stringsIds := make([]string, len(ids))
	for i, id := range ids {
		stringsIds[i] = id.String()
	}

	stringIds := strings.Join(stringsIds, ",")

	param.Add("email", stringIds)

	log.Info("email:")

    usersServiceURL := "http://users-service-2:8080"
    checkURL := fmt.Sprintf("%s/internal/users/batch-check?uuids=%s", 
        usersServiceURL, 
        param.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, checkURL, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusConflict {
			return fmt.Errorf("%s: %s", op, "user not found")
		}

		return fmt.Errorf("%s: unexpected status code %d", op, resp.StatusCode)
	}

	defer resp.Body.Close()

	var newUser NormalizedUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}