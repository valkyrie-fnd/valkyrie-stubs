package broken

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scenarios(t *testing.T) {
	tests := []struct {
		name    string
		uri     []byte
		request string
		want    bool
	}{
		{
			name:    "bet_fails",
			uri:     []byte("/players/5000001/balance?provider=evolution&currency=EUR"),
			request: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := predefinedScenarios[tt.name]
			assert.Equal(t, tt.want, routeMatch(s.PathPattern, tt.uri))
		})
	}
}
