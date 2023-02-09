package genericpam

type GameLaunchResponse struct {
	GameURL string `json:"gameUrl"`
}

// Operator game launch request
type GameLaunchRequest struct {
	LaunchConfig   LaunchConfig `json:"launchConfig"`
	PlayerID       string       `json:"playerId"`
	Provider       string       `json:"provider"`
	SessionKey     string       `json:"sessionKey"`
	ExternalGameID string       `json:"gameId"`
	Casino         string       `json:"casino"`
	Country        string       `json:"country"`
	Language       string       `json:"language"`
	Currency       string       `json:"currency"`
}

type LaunchConfig struct {
	Config map[string]interface{}
}
