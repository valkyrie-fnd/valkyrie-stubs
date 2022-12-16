package genericpam

// Very thin wrapper around http client doing requests to the valkyrie, to stimulate tests

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
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
		PlayerId:       player,
		Provider:       "Red Tiger",
		SessionKey:     "ABC",
		Casino:         "",
		Currency:       "",
		ExternalGameId: game,
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
		return nil, utils.Stack(err, fmt.Errorf("evo/%s request failed: %s", path, b))
	}

	return &resp, nil
}
