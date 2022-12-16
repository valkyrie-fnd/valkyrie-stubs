package memorydatastore

import (
	"os"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PamApiToken      string `yaml:"pamApiToken"`
	Providers        []datastore.Provider
	ProviderApiKeys  []datastore.ProviderApiKey `yaml:"providerApiKeys"`
	ProviderSessions []datastore.Session        `yaml:"providerSessions"`
	Games            []datastore.Game
	GameRounds       []datastore.GameRound `yaml:"gameRounds"`
	Accounts         []datastore.Account
	Players          []datastore.Player
	Sessions         []datastore.Session
	Transactions     []datastore.Transaction
	SessionTimeout   *int `yaml:"sessionTimeout,omitempty"`
}

func ReadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parse(data)
}

func parse(data []byte) (*Config, error) {
	data = []byte(os.ExpandEnv(string(data)))
	conf := Config{}
	err := yaml.Unmarshal(data, &conf)
	return &conf, err
}
