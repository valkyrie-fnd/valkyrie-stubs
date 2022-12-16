package genericpam

type GameLaunchResponse struct {
	GameURL string `json:"gameUrl"`
}

// Operator game launch request
type GameLaunchRequest struct {
	PlayerId       string       `json:"playerId"`
	Provider       string       `json:"provider"`
	SessionKey     string       `json:"sessionKey"`
	ExternalGameId string       `json:"gameId"`
	Casino         string       `json:"casino"`
	Country        string       `json:"country"`
	Language       string       `json:"language"`
	Currency       string       `json:"currency"`
	LaunchConfig   LaunchConfig `json:"launchConfig"`
}

type LaunchConfig struct {
	Config map[string]interface{}
}
