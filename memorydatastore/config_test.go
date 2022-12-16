package memorydatastore

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

var expectedTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
var expectedConfig = &Config{
	PamApiToken: "pam-api-token",
	Providers: []datastore.Provider{
		{
			Provider:   "Evolution",
			ProviderId: 3,
		},
	},
	ProviderApiKeys: []datastore.ProviderApiKey{
		{
			Provider: "Evolution",
			ApiKey:   "evo-api-key",
		},
	},
	ProviderSessions: []datastore.Session{
		{
			Provider: "Red Tiger",
			Key:      "RECON_TOKEN",
		},
	},
	Games: []datastore.Game{
		{
			ProviderGameId: "vctlz20yfnmp1ylr",
		},
	},
	GameRounds: []datastore.GameRound{
		{
			ProviderGameId:  "vctlz20yfnmp1ylr",
			ProviderRoundId: "vVJBwYIUc5",
			StartTime:       expectedTime,
			ProviderName:    "Evolution",
			PlayerId:        "2000001",
		},
	},
	Accounts: []datastore.Account{
		{
			Id:               3,
			PlayerIdentifier: "2000001",
			Currency:         "SEK",
			Country:          "SE",
			Language:         "sv",
			CashAmount:       100,
			BonusAmount:      10,
			PromoAmount:      1,
		},
		{
			Id:               10,
			PlayerIdentifier: "5000001",
			Currency:         "EUR",
			Country:          "SE",
			Language:         "sv",
			CashAmount:       100,
			BonusAmount:      10,
			PromoAmount:      1,
		},
	},
	Players: []datastore.Player{
		{
			Id:               2000001,
			PlayerIdentifier: "2000001",
		},
		{
			Id:               5000001,
			PlayerIdentifier: "5000001",
		},
	},
	Sessions: []datastore.Session{
		{
			Key:              "A7eK4bOmC1Ux-hbvdr4bRckEqBPDAGj06aO3bLyAR_g",
			PlayerIdentifier: "2000001",
			Provider:         "Evolution",
		},
	},
	Transactions: []datastore.Transaction{
		{
			PlayerIdentifier:      "2000001",
			CashAmount:            100,
			BonusAmount:           10,
			PromoAmount:           1,
			Currency:              "SEK",
			TransactionType:       "DEPOSIT",
			ProviderTransactionId: "123",
			ProviderBetRef:        utils.Ptr("321"),
			ProviderGameId:        "vctlz20yfnmp1ylr",
			ProviderName:          "Evolution",
			ProviderRoundId:       utils.Ptr("vVJBwYIUc5"),
			TransactionDateTime:   expectedTime,
		},
	},
}

func TestReadConfig(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "read plain yaml from config.test.yaml",
			args: args{
				file: "testdata/config.test.yaml",
			},
			want:    expectedConfig,
			wantErr: false,
		},
		{
			name: "read yaml with anchors from anchors.test.yaml",
			args: args{
				file: "testdata/anchors.test.yaml",
			},
			want:    expectedConfig,
			wantErr: false,
		},
		{
			name: "read missing file missing.test.yaml fails",
			args: args{
				file: "testdata/missing.test.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "read yaml with environment substitution env.test.yaml",
			args: args{
				file: "testdata/env.test.yaml",
			},
			want:    &Config{PamApiToken: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runWithEnvAndReset("testdata/.env.testing", func() {
				got, err := ReadConfig(tt.args.file)
				if (err != nil) != tt.wantErr {
					t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parse() got = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

// Just runs the func with env vars set and then clears the vars
func runWithEnvAndReset(file string, fn func()) {
	if file == "" {
		fn()
		return
	}
	vars, _ := godotenv.Read(file)
	_ = godotenv.Overload(file)

	fn()

	for k := range vars {
		_ = os.Unsetenv(k)
	}
}
