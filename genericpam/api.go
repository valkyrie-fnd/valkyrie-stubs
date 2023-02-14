package genericpam

// Very thin wrapper around http client doing requests to the valkyrie, to stimulate tests

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ValkClient struct {
	provider  string
	url       string
	authToken string
}

func NewValk(provider, baseURL, token string) *ValkClient {
	return &ValkClient{provider, baseURL, token}
}

func (e *ValkClient) GameLaunch(game, player string) (*GameLaunchResponse, error) {
	r := GameLaunchRequest{
		PlayerID:       player,
		Provider:       "Red Tiger",
		SessionKey:     "ABC",
		Casino:         "",
		Currency:       "",
		ExternalGameID: game,
	}
	return post(e.url, "balance", e.authToken, r)
}

func post(base, path, token string, body interface{}) (*GameLaunchResponse, error) {
	a := fiber.Post(base + path).
		QueryString(fmt.Sprintf("authToken=%s", token)).
		Timeout(5 * time.Second).
		JSON(body)

	var resp GameLaunchResponse
	status, b, err := a.Struct(&resp)
	if status != fiber.StatusOK {
		return nil, fmt.Errorf("valkyrie/%s request failed with status [%v]: %s, Error: %s", path, status, string(b), err)
	} else if err != nil {
		return nil, errors.Join(err...)
	}

	return &resp, nil
}
