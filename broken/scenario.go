package broken

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"

	"github.com/valkyrie-fnd/valkyrie-stubs/genericpam"
)

type scenario struct {
	Title          string
	Description    string
	PathPattern    *regexp.Regexp
	RequestPattern *regexp.Regexp // consider using some jsonpath library for better accuracy
	Response       any
	HardError      string
}

var predefinedScenarios = map[string]scenario{
	"bet_fails": {
		Title:          "Withdrawal fails",
		Description:    "Next withdrawal/bet/debit is rejected",
		PathPattern:    regexp.MustCompile(`/players/.*/transactions$`),
		RequestPattern: regexp.MustCompile(`"transactionType":\s?"WITHDRAW"`),
		Response: genericpam.AddTransactionResponse{
			Status: genericpam.ERROR,
			Error: &genericpam.PamError{
				Code:    genericpam.PAMERRUNDEFINED,
				Message: "forced error",
			},
		},
	},
	"balance_fails": {
		Title:          "Balance fails",
		Description:    "Next balance request fails",
		PathPattern:    regexp.MustCompile(`/players/.*/balance$`),
		RequestPattern: regexp.MustCompile(``),
		Response: genericpam.AddTransactionResponse{
			Status: genericpam.ERROR,
			Error: &genericpam.PamError{
				Code:    genericpam.PAMERRUNDEFINED,
				Message: "forced error",
			},
		},
	},
}

func (s *scenario) match(r *fiber.Request) bool {
	return routeMatch(s.PathPattern, r.URI().Path())
}

func routeMatch(pat *regexp.Regexp, uri []byte) bool {
	return pat.Match(uri)
}

func (s *scenario) String() string {
	return fmt.Sprintf("scenario [title: %s]", s.Title)
}
