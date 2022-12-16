// Package genericpam provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package genericpam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gofiber/fiber/v2"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get session details
	// (GET /players/session)
	GetSession(c *fiber.Ctx, params GetSessionParams) error
	// Refresh player session
	// (PUT /players/session)
	RefreshSession(c *fiber.Ctx, params RefreshSessionParams) error
	// Get player account balance
	// (GET /players/{playerId}/balance)
	GetBalance(c *fiber.Ctx, playerId PlayerId, params GetBalanceParams) error
	// Get game rounds
	// (GET /players/{playerId}/gamerounds/{providerGameRoundId})
	GetGameRound(c *fiber.Ctx, playerId PlayerId, providerGameRoundId ProviderRoundId, params GetGameRoundParams) error
	// Get transactions
	// (GET /players/{playerId}/transactions)
	GetTransactions(c *fiber.Ctx, playerId PlayerId, params GetTransactionsParams) error
	// Create transaction for one player
	// (POST /players/{playerId}/transactions)
	AddTransaction(c *fiber.Ctx, playerId PlayerId, params AddTransactionParams) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type MiddlewareFunc fiber.Handler

// GetSession operation middleware
func (siw *ServerInterfaceWrapper) GetSession(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSessionParams

	query, err := url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "provider" -------------

	if paramValue := c.Query("provider"); paramValue != "" {

	} else {
		err = fmt.Errorf("Query argument provider is required, but not found")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return err
	}

	err = runtime.BindQueryParameter("form", true, true, "provider", query, &params.Provider)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter provider: %w", err).Error())
	}

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-Player-Token" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Player-Token")]; found {
		var XPlayerToken SessionToken

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Player-Token", runtime.ParamLocationHeader, value, &XPlayerToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Player-Token: %w", err).Error())
		}

		params.XPlayerToken = XPlayerToken

	} else {
		err = fmt.Errorf("Header parameter X-Player-Token is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// ------------- Required header parameter "X-Correlation-ID" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Correlation-ID")]; found {
		var XCorrelationID CorrelationId

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Correlation-ID", runtime.ParamLocationHeader, value, &XCorrelationID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Correlation-ID: %w", err).Error())
		}

		params.XCorrelationID = XCorrelationID

	} else {
		err = fmt.Errorf("Header parameter X-Correlation-ID is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.GetSession(c, params)
}

// RefreshSession operation middleware
func (siw *ServerInterfaceWrapper) RefreshSession(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params RefreshSessionParams

	query, err := url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "provider" -------------

	if paramValue := c.Query("provider"); paramValue != "" {

	} else {
		err = fmt.Errorf("Query argument provider is required, but not found")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return err
	}

	err = runtime.BindQueryParameter("form", true, true, "provider", query, &params.Provider)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter provider: %w", err).Error())
	}

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-Player-Token" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Player-Token")]; found {
		var XPlayerToken SessionToken

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Player-Token", runtime.ParamLocationHeader, value, &XPlayerToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Player-Token: %w", err).Error())
		}

		params.XPlayerToken = XPlayerToken

	} else {
		err = fmt.Errorf("Header parameter X-Player-Token is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// ------------- Required header parameter "X-Correlation-ID" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Correlation-ID")]; found {
		var XCorrelationID CorrelationId

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Correlation-ID", runtime.ParamLocationHeader, value, &XCorrelationID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Correlation-ID: %w", err).Error())
		}

		params.XCorrelationID = XCorrelationID

	} else {
		err = fmt.Errorf("Header parameter X-Correlation-ID is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.RefreshSession(c, params)
}

// GetBalance operation middleware
func (siw *ServerInterfaceWrapper) GetBalance(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "playerId" -------------
	var playerId PlayerId

	err = runtime.BindStyledParameter("simple", false, "playerId", c.Params("playerId"), &playerId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter playerId: %w", err).Error())
	}

	c.Context().SetUserValue(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetBalanceParams

	query, err := url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "provider" -------------

	if paramValue := c.Query("provider"); paramValue != "" {

	} else {
		err = fmt.Errorf("Query argument provider is required, but not found")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return err
	}

	err = runtime.BindQueryParameter("form", true, true, "provider", query, &params.Provider)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter provider: %w", err).Error())
	}

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-Player-Token" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Player-Token")]; found {
		var XPlayerToken SessionToken

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Player-Token", runtime.ParamLocationHeader, value, &XPlayerToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Player-Token: %w", err).Error())
		}

		params.XPlayerToken = XPlayerToken

	} else {
		err = fmt.Errorf("Header parameter X-Player-Token is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// ------------- Required header parameter "X-Correlation-ID" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Correlation-ID")]; found {
		var XCorrelationID CorrelationId

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Correlation-ID", runtime.ParamLocationHeader, value, &XCorrelationID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Correlation-ID: %w", err).Error())
		}

		params.XCorrelationID = XCorrelationID

	} else {
		err = fmt.Errorf("Header parameter X-Correlation-ID is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.GetBalance(c, playerId, params)
}

// GetGameRound operation middleware
func (siw *ServerInterfaceWrapper) GetGameRound(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "playerId" -------------
	var playerId PlayerId

	err = runtime.BindStyledParameter("simple", false, "playerId", c.Params("playerId"), &playerId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter playerId: %w", err).Error())
	}

	// ------------- Path parameter "providerGameRoundId" -------------
	var providerGameRoundId ProviderRoundId

	err = runtime.BindStyledParameter("simple", false, "providerGameRoundId", c.Params("providerGameRoundId"), &providerGameRoundId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter providerGameRoundId: %w", err).Error())
	}

	c.Context().SetUserValue(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetGameRoundParams

	query, err := url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "provider" -------------

	if paramValue := c.Query("provider"); paramValue != "" {

	} else {
		err = fmt.Errorf("Query argument provider is required, but not found")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return err
	}

	err = runtime.BindQueryParameter("form", true, true, "provider", query, &params.Provider)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter provider: %w", err).Error())
	}

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-Player-Token" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Player-Token")]; found {
		var XPlayerToken SessionToken

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Player-Token", runtime.ParamLocationHeader, value, &XPlayerToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Player-Token: %w", err).Error())
		}

		params.XPlayerToken = XPlayerToken

	} else {
		err = fmt.Errorf("Header parameter X-Player-Token is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// ------------- Required header parameter "X-Correlation-ID" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Correlation-ID")]; found {
		var XCorrelationID CorrelationId

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Correlation-ID", runtime.ParamLocationHeader, value, &XCorrelationID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Correlation-ID: %w", err).Error())
		}

		params.XCorrelationID = XCorrelationID

	} else {
		err = fmt.Errorf("Header parameter X-Correlation-ID is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.GetGameRound(c, playerId, providerGameRoundId, params)
}

// GetTransactions operation middleware
func (siw *ServerInterfaceWrapper) GetTransactions(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "playerId" -------------
	var playerId PlayerId

	err = runtime.BindStyledParameter("simple", false, "playerId", c.Params("playerId"), &playerId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter playerId: %w", err).Error())
	}

	c.Context().SetUserValue(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTransactionsParams

	query, err := url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "provider" -------------

	if paramValue := c.Query("provider"); paramValue != "" {

	} else {
		err = fmt.Errorf("Query argument provider is required, but not found")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return err
	}

	err = runtime.BindQueryParameter("form", true, true, "provider", query, &params.Provider)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter provider: %w", err).Error())
	}

	// ------------- Optional query parameter "providerTransactionId" -------------

	err = runtime.BindQueryParameter("form", true, false, "providerTransactionId", query, &params.ProviderTransactionId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter providerTransactionId: %w", err).Error())
	}

	// ------------- Optional query parameter "providerBetRef" -------------

	err = runtime.BindQueryParameter("form", true, false, "providerBetRef", query, &params.ProviderBetRef)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter providerBetRef: %w", err).Error())
	}

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-Player-Token" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Player-Token")]; found {
		var XPlayerToken SessionToken

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Player-Token", runtime.ParamLocationHeader, value, &XPlayerToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Player-Token: %w", err).Error())
		}

		params.XPlayerToken = XPlayerToken

	} else {
		err = fmt.Errorf("Header parameter X-Player-Token is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// ------------- Required header parameter "X-Correlation-ID" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Correlation-ID")]; found {
		var XCorrelationID CorrelationId

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Correlation-ID", runtime.ParamLocationHeader, value, &XCorrelationID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Correlation-ID: %w", err).Error())
		}

		params.XCorrelationID = XCorrelationID

	} else {
		err = fmt.Errorf("Header parameter X-Correlation-ID is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.GetTransactions(c, playerId, params)
}

// AddTransaction operation middleware
func (siw *ServerInterfaceWrapper) AddTransaction(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "playerId" -------------
	var playerId PlayerId

	err = runtime.BindStyledParameter("simple", false, "playerId", c.Params("playerId"), &playerId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter playerId: %w", err).Error())
	}

	c.Context().SetUserValue(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params AddTransactionParams

	query, err := url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Required query parameter "provider" -------------

	if paramValue := c.Query("provider"); paramValue != "" {

	} else {
		err = fmt.Errorf("Query argument provider is required, but not found")
		c.Status(fiber.StatusBadRequest).JSON(err)
		return err
	}

	err = runtime.BindQueryParameter("form", true, true, "provider", query, &params.Provider)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter provider: %w", err).Error())
	}

	headers := c.GetReqHeaders()

	// ------------- Required header parameter "X-Player-Token" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Player-Token")]; found {
		var XPlayerToken SessionToken

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Player-Token", runtime.ParamLocationHeader, value, &XPlayerToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Player-Token: %w", err).Error())
		}

		params.XPlayerToken = XPlayerToken

	} else {
		err = fmt.Errorf("Header parameter X-Player-Token is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// ------------- Required header parameter "X-Correlation-ID" -------------
	if value, found := headers[http.CanonicalHeaderKey("X-Correlation-ID")]; found {
		var XCorrelationID CorrelationId

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Correlation-ID", runtime.ParamLocationHeader, value, &XCorrelationID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter X-Correlation-ID: %w", err).Error())
		}

		params.XCorrelationID = XCorrelationID

	} else {
		err = fmt.Errorf("Header parameter X-Correlation-ID is required, but not found: %w", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return siw.Handler.AddTransaction(c, playerId, params)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	for _, m := range options.Middlewares {
		router.Use(m)
	}

	router.Get(options.BaseURL+"/players/session", wrapper.GetSession)

	router.Put(options.BaseURL+"/players/session", wrapper.RefreshSession)

	router.Get(options.BaseURL+"/players/:playerId/balance", wrapper.GetBalance)

	router.Get(options.BaseURL+"/players/:playerId/gamerounds/:providerGameRoundId", wrapper.GetGameRound)

	router.Get(options.BaseURL+"/players/:playerId/transactions", wrapper.GetTransactions)

	router.Post(options.BaseURL+"/players/:playerId/transactions", wrapper.AddTransaction)

}

type UnauthorizedJSONResponse struct {
	// Error Error details describing why PAM rejected the request
	Error struct {
		// Code Pam Error code "PAM_ERR_SESSION_NOT_FOUND" or "PAM_ERR_UNDEFINED"
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`

	// Status Status is Error
	Status string `json:"status"`
}

type GetSessionRequestObject struct {
	Params GetSessionParams
}

type GetSessionResponseObject interface {
	VisitGetSessionResponse(ctx *fiber.Ctx) error
}

type GetSession200JSONResponse SessionResponse

func (response GetSession200JSONResponse) VisitGetSessionResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetSession401JSONResponse = UnauthorizedJSONResponse

func (response GetSession401JSONResponse) VisitGetSessionResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type RefreshSessionRequestObject struct {
	Params RefreshSessionParams
}

type RefreshSessionResponseObject interface {
	VisitRefreshSessionResponse(ctx *fiber.Ctx) error
}

type RefreshSession200JSONResponse SessionResponse

func (response RefreshSession200JSONResponse) VisitRefreshSessionResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type RefreshSession401JSONResponse = UnauthorizedJSONResponse

func (response RefreshSession401JSONResponse) VisitRefreshSessionResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type GetBalanceRequestObject struct {
	PlayerId PlayerId `json:"playerId"`
	Params   GetBalanceParams
}

type GetBalanceResponseObject interface {
	VisitGetBalanceResponse(ctx *fiber.Ctx) error
}

type GetBalance200JSONResponse BalanceResponse

func (response GetBalance200JSONResponse) VisitGetBalanceResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetBalance401JSONResponse = UnauthorizedJSONResponse

func (response GetBalance401JSONResponse) VisitGetBalanceResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type GetGameRoundRequestObject struct {
	PlayerId            PlayerId        `json:"playerId"`
	ProviderGameRoundId ProviderRoundId `json:"providerGameRoundId"`
	Params              GetGameRoundParams
}

type GetGameRoundResponseObject interface {
	VisitGetGameRoundResponse(ctx *fiber.Ctx) error
}

type GetGameRound200JSONResponse GameRoundResponse

func (response GetGameRound200JSONResponse) VisitGetGameRoundResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetGameRound401JSONResponse = UnauthorizedJSONResponse

func (response GetGameRound401JSONResponse) VisitGetGameRoundResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type GetTransactionsRequestObject struct {
	PlayerId PlayerId `json:"playerId"`
	Params   GetTransactionsParams
}

type GetTransactionsResponseObject interface {
	VisitGetTransactionsResponse(ctx *fiber.Ctx) error
}

type GetTransactions200JSONResponse GetTransactionsResponse

func (response GetTransactions200JSONResponse) VisitGetTransactionsResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetTransactions400JSONResponse struct {
	// Error Error details describing why PAM rejected the request
	Error struct {
		// Code Pam Error code "PAM_ERR_MISSING_PROVIDER" or "PAM_ERR_UNDEFINED"
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`

	// Status Status is error
	Status string `json:"status"`
}

func (response GetTransactions400JSONResponse) VisitGetTransactionsResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(400)

	return ctx.JSON(&response)
}

type GetTransactions401JSONResponse = UnauthorizedJSONResponse

func (response GetTransactions401JSONResponse) VisitGetTransactionsResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

type AddTransactionRequestObject struct {
	PlayerId PlayerId `json:"playerId"`
	Params   AddTransactionParams
	Body     *AddTransactionJSONRequestBody
}

type AddTransactionResponseObject interface {
	VisitAddTransactionResponse(ctx *fiber.Ctx) error
}

type AddTransaction200JSONResponse AddTransactionResponse

func (response AddTransaction200JSONResponse) VisitAddTransactionResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type AddTransaction400Response struct {
}

func (response AddTransaction400Response) VisitAddTransactionResponse(ctx *fiber.Ctx) error {
	ctx.Status(400)
	return nil
}

type AddTransaction401JSONResponse = UnauthorizedJSONResponse

func (response AddTransaction401JSONResponse) VisitAddTransactionResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(401)

	return ctx.JSON(&response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get session details
	// (GET /players/session)
	GetSession(ctx context.Context, request GetSessionRequestObject) (GetSessionResponseObject, error)
	// Refresh player session
	// (PUT /players/session)
	RefreshSession(ctx context.Context, request RefreshSessionRequestObject) (RefreshSessionResponseObject, error)
	// Get player account balance
	// (GET /players/{playerId}/balance)
	GetBalance(ctx context.Context, request GetBalanceRequestObject) (GetBalanceResponseObject, error)
	// Get game rounds
	// (GET /players/{playerId}/gamerounds/{providerGameRoundId})
	GetGameRound(ctx context.Context, request GetGameRoundRequestObject) (GetGameRoundResponseObject, error)
	// Get transactions
	// (GET /players/{playerId}/transactions)
	GetTransactions(ctx context.Context, request GetTransactionsRequestObject) (GetTransactionsResponseObject, error)
	// Create transaction for one player
	// (POST /players/{playerId}/transactions)
	AddTransaction(ctx context.Context, request AddTransactionRequestObject) (AddTransactionResponseObject, error)
}

type StrictHandlerFunc func(ctx *fiber.Ctx, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetSession operation middleware
func (sh *strictHandler) GetSession(ctx *fiber.Ctx, params GetSessionParams) error {
	var request GetSessionRequestObject

	request.Params = params

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetSession(ctx.UserContext(), request.(GetSessionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetSession")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetSessionResponseObject); ok {
		if err := validResponse.VisitGetSessionResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// RefreshSession operation middleware
func (sh *strictHandler) RefreshSession(ctx *fiber.Ctx, params RefreshSessionParams) error {
	var request RefreshSessionRequestObject

	request.Params = params

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.RefreshSession(ctx.UserContext(), request.(RefreshSessionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "RefreshSession")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(RefreshSessionResponseObject); ok {
		if err := validResponse.VisitRefreshSessionResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetBalance operation middleware
func (sh *strictHandler) GetBalance(ctx *fiber.Ctx, playerId PlayerId, params GetBalanceParams) error {
	var request GetBalanceRequestObject

	request.PlayerId = playerId
	request.Params = params

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetBalance(ctx.UserContext(), request.(GetBalanceRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetBalance")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetBalanceResponseObject); ok {
		if err := validResponse.VisitGetBalanceResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetGameRound operation middleware
func (sh *strictHandler) GetGameRound(ctx *fiber.Ctx, playerId PlayerId, providerGameRoundId ProviderRoundId, params GetGameRoundParams) error {
	var request GetGameRoundRequestObject

	request.PlayerId = playerId
	request.ProviderGameRoundId = providerGameRoundId
	request.Params = params

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetGameRound(ctx.UserContext(), request.(GetGameRoundRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetGameRound")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetGameRoundResponseObject); ok {
		if err := validResponse.VisitGetGameRoundResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetTransactions operation middleware
func (sh *strictHandler) GetTransactions(ctx *fiber.Ctx, playerId PlayerId, params GetTransactionsParams) error {
	var request GetTransactionsRequestObject

	request.PlayerId = playerId
	request.Params = params

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetTransactions(ctx.UserContext(), request.(GetTransactionsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetTransactions")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetTransactionsResponseObject); ok {
		if err := validResponse.VisitGetTransactionsResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// AddTransaction operation middleware
func (sh *strictHandler) AddTransaction(ctx *fiber.Ctx, playerId PlayerId, params AddTransactionParams) error {
	var request AddTransactionRequestObject

	request.PlayerId = playerId
	request.Params = params

	var body AddTransactionJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.AddTransaction(ctx.UserContext(), request.(AddTransactionRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AddTransaction")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(AddTransactionResponseObject); ok {
		if err := validResponse.VisitAddTransactionResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}
