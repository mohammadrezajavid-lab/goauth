package http

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/service/goauth"
	"github.com/mohammadrezajavid-lab/goauth/pkg/phonenumber"
	nethttp "net/http"
	"strconv"
)

type Handler struct {
	authSvc       goauth.Service
	authValidator goauth.Validator
}

func NewHandler(authSvc goauth.Service, authValidator goauth.Validator) Handler {
	return Handler{
		authSvc:       authSvc,
		authValidator: authValidator,
	}
}

// --- Auth Endpoints ---

// GenerateOTPCode godoc
// @Summary Generate OTP
// @Description Generates a one-time password for a given phone number.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body goauth.GenerateOTPRequest true "Phone Number"
// @Success 200 {object} map[string]string "message: OTP code generated successfully"
// @Failure 400 {object} map[string]interface{} "Validation failed or bad request"
// @Failure 422 {object} map[string]interface{} "Invalid phone number format"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/generateotp [post]
func (h *Handler) GenerateOTPCode(ctx echo.Context) error {
	request, ok := ctx.Get("request").(*goauth.GenerateOTPRequest)
	if !ok {
		request = new(goauth.GenerateOTPRequest)
		if err := ctx.Bind(request); err != nil {
			return echo.NewHTTPError(nethttp.StatusBadRequest, "invalid request body")
		}
	}

	phoneNumber, err := phonenumber.NewNormalizer().NormalizePhoneNumber(request.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(nethttp.StatusUnprocessableEntity, "invalid phone number format")
	}
	request.PhoneNumber = phoneNumber

	if fieldErrors, validateErr := h.authValidator.ValidateGenerateOTPRequest(request); validateErr != nil {
		return ctx.JSON(nethttp.StatusBadRequest, echo.Map{
			"message": "Validation failed",
			"errors":  fieldErrors,
		})
	}

	_, gErr := h.authSvc.GenerateOTP(ctx.Request().Context(), request)
	if gErr != nil {
		return handleServiceError(ctx, gErr)
	}

	return ctx.JSON(nethttp.StatusOK, echo.Map{
		"message": "OTP code has been generated and printed to the console.",
	})
}

// VerifyAndLoginOrRegister godoc
// @Summary Verify OTP and Login/Register
// @Description Verifies an OTP and returns a JWT token for the user. Registers the user if they don't exist.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body goauth.VerifyOTPRequest true "Phone Number and OTP"
// @Success 200 {object} goauth.VerifyOTPResponse
// @Failure 400 {object} map[string]interface{} "Validation failed or bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 422 {object} map[string]interface{} "Invalid phone number format"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/verify [post]
func (h *Handler) VerifyAndLoginOrRegister(ctx echo.Context) error {
	var request = &goauth.VerifyOTPRequest{}
	if err := ctx.Bind(request); err != nil {
		return echo.NewHTTPError(nethttp.StatusBadRequest, err.Error())
	}

	normalizedPhone, err := phonenumber.NewNormalizer().NormalizePhoneNumber(request.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(nethttp.StatusUnprocessableEntity, "invalid phone number format")
	}
	request.PhoneNumber = normalizedPhone

	if fieldErrors, validateErr := h.authValidator.ValidateVerifyOTPRequest(request); validateErr != nil {
		return ctx.JSON(nethttp.StatusBadRequest, echo.Map{
			"message": "Validation failed",
			"errors":  fieldErrors,
		})
	}

	verifyRes, vErr := h.authSvc.VerifyAndLogin(ctx.Request().Context(), request)
	if vErr != nil {
		return handleServiceError(ctx, vErr)
	}

	return ctx.JSON(nethttp.StatusOK, verifyRes)
}

// --- User Management Endpoints ---

// GetUser godoc
// @Summary Get User by ID
// @Description Retrieves details for a single user by their ID.
// @Tags users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} goauth.User
// @Failure 400 {object} map[string]interface{} "Invalid user ID format"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [get]
func (h *Handler) GetUser(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(nethttp.StatusBadRequest, echo.Map{"error": "invalid user ID format"})
	}

	req := &goauth.GetUserRequest{ID: id}
	user, err := h.authSvc.GetUser(ctx.Request().Context(), req)
	if err != nil {
		return handleServiceError(ctx, err)
	}

	return ctx.JSON(nethttp.StatusOK, user)
}

// ListUsers godoc
// @Summary List Users
// @Description Retrieves a paginated list of users with search functionality.
// @Tags users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param search query string false "Search term for phone number"
// @Success 200 {object} goauth.ListUsersResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [get]
func (h *Handler) ListUsers(ctx echo.Context) error {
	// Parse query parameters with defaults
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	pageSize, _ := strconv.Atoi(ctx.QueryParam("pageSize"))
	search := ctx.QueryParam("search")

	req := &goauth.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	response, err := h.authSvc.ListUsers(ctx.Request().Context(), req)
	if err != nil {
		return handleServiceError(ctx, err)
	}

	return ctx.JSON(nethttp.StatusOK, response)
}
