package genericpam

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	XPlayerToken  = "X-Player-Token"
	QueryProvider = "provider"
)

func (c Controller) registerMiddlewares(app fiber.Router, pamAPIKey string, providerTokens map[string]string) {
	app.Use(c.getCheckPamAPIToken(pamAPIKey))
	app.Use(c.getCheckPlayerToken(providerTokens))
	// for `ctx.Params("playerId")` used in c.checkPlayerId to work, the middleware has to be registered on a
	// prefix containing the path param
	app.Use("/players/:playerId/+", c.checkPlayerID)
	app.Use(c.checkProvider)
}

func (c Controller) getCheckPamAPIToken(pamAPIKey string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		apiKey := ctx.GetReqHeaders()[fiber.HeaderAuthorization]
		if apiKey == "" || strings.TrimPrefix(apiKey, "Bearer ") != pamAPIKey {
			return ctx.Status(http.StatusUnauthorized).JSON(BaseResponse{
				Error: &PamError{
					Code:    PAMERRAPITOKEN,
					Message: "not authorized",
				},
				Status: ERROR,
			})
		}
		return ctx.Next()
	}
}

func (c Controller) getCheckPlayerToken(providerTokens map[string]string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		sessionToken := ctx.GetReqHeaders()[XPlayerToken]
		session, err := c.dataStore.GetSession(ctx.UserContext(), sessionToken)
		if err != nil {
			// Last chance, check if token is a provider specific token (reconciliation)
			return c.checkProviderToken(ctx, sessionToken, providerTokens)
		}

		// Verify that provider matches session
		provider := ctx.Query(QueryProvider)
		if session.Provider != provider {
			return ctx.Status(http.StatusUnauthorized).JSON(BaseResponse{
				Error: &PamError{
					Code:    PAMERRSESSIONNOTFOUND,
					Message: "session not found",
				},
				Status: ERROR,
			})
		}

		// Store the session in context
		ctx.Context().SetUserValue("session", session)
		return ctx.Next()
	}
}

func (c Controller) checkProviderToken(ctx *fiber.Ctx, sessionToken string, providerTokens map[string]string) error {
	provider := ctx.Query(QueryProvider)

	if pt, found := providerTokens[provider]; !found || sessionToken != pt {
		return ctx.Status(http.StatusUnauthorized).JSON(BaseResponse{
			Error: &PamError{
				Code:    PAMERRSESSIONNOTFOUND,
				Message: "session not found",
			},
			Status: ERROR,
		})
	}
	return ctx.Next()
}

func (c Controller) checkPlayerID(ctx *fiber.Ctx) error {
	uriPlayerID := ctx.Params("playerId")
	if uriPlayerID != "" {
		_, err := c.dataStore.GetPlayer(ctx.UserContext(), uriPlayerID)
		if err != nil {
			return ctx.Status(http.StatusOK).JSON(BaseResponse{
				Error: &PamError{
					Code:    PAMERRPLAYERNOTFOUND,
					Message: "User not found",
				},
				Status: ERROR,
			})
		}
	}
	return ctx.Next()
}

func (c Controller) checkProvider(ctx *fiber.Ctx) error {
	provider := ctx.Query(QueryProvider)
	if provider != "" {
		_, err := c.dataStore.GetProvider(ctx.UserContext(), provider)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(BaseResponse{
				Error: &PamError{
					Code:    PAMERRMISSINGPROVIDER,
					Message: "missing provider",
				},
				Status: ERROR,
			})
		}
	}
	return ctx.Next()
}
