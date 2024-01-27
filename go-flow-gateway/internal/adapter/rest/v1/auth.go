package v1

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/bgg/go-flow-gateway/config"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	domainUrl      string
	userProfile    usecase.UserProfile
	logger         logger.Logger
	oauthDetail    usecase.OAuthDetail
	userCredential usecase.UserCredential
	lineChannelID  string
}

func NewAuthRoutes(cfg *config.Config, handler *gin.RouterGroup, u usecase.UserProfile, l logger.Logger, o usecase.OAuthDetail, c usecase.UserCredential, lineChannelID string) {
	r := authRoutes{domainUrl: cfg.App.DomainUrl, userProfile: u, logger: l, oauthDetail: o, userCredential: c, lineChannelID: lineChannelID}
	auth := handler.Group("/auth")
	{
		auth.POST("/register", r.register)
		auth.POST("/login", r.login)
		auth.GET("/line-login", r.lineLogin)       // Initiate Line Login
		auth.GET("/line-callback", r.lineCallback) // Handler the redirect from Line Login
		auth.GET("/logout", r.logout)
	}
}

type RegisterRequest struct {
	DisplayName string `json:"displayName" example:"billyang"`
	Username    string `json:"username" binding:"required" example:"useraname"`
	Password    string `json:"password" binding:"required" example:"password"`
}

type RegisterResponse struct {
	UserID string `json:"userID" example:"userID"`
}

// register godoc
//
//	@Summary		Register
//	@Description	Register
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			RegisterRequest	body		RegisterRequest		true	"register information"
//	@Success		200				{object}	RegisterResponse	"Succesfully registered"
//	@Router			/auth/register [post]
func (r *authRoutes) register(c *gin.Context) {
	var req RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		r.logger.Warn("http - v1 - register: invalid request body", err)
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := r.userCredential.Register(c.Request.Context(), req.DisplayName, req.Username, req.Password)
	if err != nil {
		r.logger.Error("http - v1 - register: register failed", err)
		sendErrorResponse(c, http.StatusInternalServerError, "register failed")
		return
	}

	r.logger.Info("http - v1 - register: register successfully", "userID", userID)
	c.JSON(http.StatusOK, RegisterResponse{UserID: fmt.Sprint(userID)})
}

type LoginRequest struct {
	Username string `json:"username" example:"username" binding:"required"`
	Password string `json:"password" example:"password" binding:"required"`
}

// login godoc
//
//	@Summary		Login
//	@Description	Login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			LoginRequest	body		LoginRequest	true	"login information"
//	@Success		204				{object}	nil				"No content"
//	@Router			/auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		r.logger.Error("http - v1 - login: invalid request body", err)
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	uc, err := r.userCredential.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		r.logger.Error("http - v1 - login: login failed", err)
		sendErrorResponse(c, http.StatusUnauthorized, "login failed")
		return
	}

	err = r.setUserSession(c, uc.UserID)
	if err != nil {
		r.logger.Error("http - v1 - login: failed to set user session", err)
		sendErrorResponse(c, http.StatusInternalServerError, "login failed")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// generateState godoc
//
// prepare a random string for the state parameter to prevent CSRF attacks
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// lineLogin godoc
//
//	@Summary		Line Login
//	@Description	Redirect to Line Login
//	@Tags			Auth
//	@Produce		html
//	@Success		302	{string}	string	Location	"Redirect URL"
//	@Router			/auth/line-login [get]
func (r *authRoutes) lineLogin(c *gin.Context) {
	state := generateState()
	lineAuthUrl := fmt.Sprintf("https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=profile%%20openid&nonce=%s",
		url.QueryEscape(os.Getenv("LINE_CHANNEL_ID")),
		url.QueryEscape(fmt.Sprintf("%s/api/v1/auth/line-callback", r.domainUrl)),
		state,
		state, // nonce can be the same as state for simplicity in this example
	)

	r.logger.Info("http - v1 - lineLogin: redirect to line login")
	c.Redirect(http.StatusTemporaryRedirect, lineAuthUrl)
}

// lineCallback godoc
//
//	@Summary		Line Callback
//	@Description	Handler the redirect from Line Login and set user session
//	@Tags			Auth
//	@Param			code	query		string	true		"Authorization code returned from Line Login"
//	@Success		302		{string}	string	Location	"Redirect URL"
//	@Router			/auth/line-callback [get]
func (r *authRoutes) lineCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		r.logger.Warn("http - v1 - lineCallback: code is empty")
		sendErrorResponse(c, http.StatusBadRequest, "code is empty")
		return
	}

	oAuthDetail, err := r.oauthDetail.HandleOAuthCallback(c.Request.Context(), code, r.domainUrl, "line", r.lineChannelID)
	if err != nil {
		r.logger.Error("http - v1 - lineCallback: failed to handle oauth callback", err)
		sendErrorResponse(c, http.StatusInternalServerError, "failed to handle oauth callback")
		return
	}

	err = r.setUserSession(c, oAuthDetail.UserID)
	if err != nil {
		r.logger.Error("http - v1 - lineCallback: failed to set user session", err)
		sendErrorResponse(c, http.StatusInternalServerError, "login failed")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, r.domainUrl)
}

// logout godoc
//
//	@Summary		Logout
//	@Description	Logout
//	@Tags			Auth
//	@Success		204	{object}	nil	"No content"
//	@Router			/auth/logout [get]
func (r *authRoutes) logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userID")
	err := session.Save()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (r *authRoutes) setUserSession(c *gin.Context, userID int) error {
	session := sessions.Default(c)
	session.Set("userID", userID)
	err := session.Save()
	if err != nil {
		return err
	}

	return nil
}
