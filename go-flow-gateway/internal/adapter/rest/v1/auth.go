package v1

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bgg/go-flow-gateway/config"
	"github.com/bgg/go-flow-gateway/internal/adapter/rest/token"
	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/infra/repo"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	domainUrl string
	u         usecase.UserProfile
	l         logger.Logger
	o         usecase.OAuthDetail
	c         usecase.UserCredential
}

func NewAuthRoutes(cfg *config.Config, handler *gin.RouterGroup, u usecase.UserProfile, l logger.Logger, o usecase.OAuthDetail, c usecase.UserCredential) {
	r := authRoutes{domainUrl: cfg.App.DomainUrl, u: u, l: l, o: o, c: c}
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
	DisplayName string `json:"displayName"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// register godoc
//
// @Summary Register
// @Description Register
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "body"
// @Success 200 {string} string	"ok"
// @Router /auth/register [post]
func (r *authRoutes) register(c *gin.Context) {
	var req RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userProfile := entity.UserProfile{
		DisplayName: req.DisplayName,
	}

	up, err := r.u.Create(c.Request.Context(), userProfile)
	if err != nil {
		r.l.Error(err, "http - v1 - register")
		sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
		return
	}

	err = r.c.Create(c.Request.Context(), up.UserID, req.Username, req.Password)
	if err != nil {
		r.l.Error(err, "http - v1 - register")
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// login godoc
//
// @Summary Login
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "body"
// @Success 200 {string} string	"ok"
// @Router /auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	uc, err := r.c.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "login failed")
		return
	}

	session := sessions.Default(c)
	session.Set("userID", uc.UserID)
	err = session.Save()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
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
// @Summary Line Login
// @Description Redirect to Line Login
// @Tags Auth
// @Success 302 {string} string	"ok"
// @Router /auth/line-login [get]
func (r *authRoutes) lineLogin(c *gin.Context) {
	state := generateState()
	lineAuthUrl := fmt.Sprintf("https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=profile%%20openid&nonce=%s",
		url.QueryEscape(os.Getenv("LINE_CHANNEL_ID")),
		url.QueryEscape(fmt.Sprintf("%s/api/v1/auth/line-callback", r.domainUrl)),
		state,
		state, // nonce can be the same as state for simplicity in this example
	)

	c.Redirect(http.StatusTemporaryRedirect, lineAuthUrl)
}

// lineCallback godoc
//
// @Summary Line Callback
// @Description Handler the redirect from Line Login
// @Tags Auth
// @Param code query string true "code"
// @Success 200 {string} string	"ok"
// @Router /auth/line-callback [get]
func (r *authRoutes) lineCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		sendErrorResponse(c, http.StatusBadRequest, "code is empty")
		return
	}

	tokens, err := r.exchangeCodeForTokens(code)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	lineUserProfile, err := token.VerifyLineIDToken(tokens.IDToken, os.Getenv("LINE_CHANNEL_ID"))
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// check the oauth detail is exists
	od, err := r.o.GetByOAuthID(c.Request.Context(), lineUserProfile.Sub)
	if err != nil {
		// oauth detail is not found,then create a new oauth detail
		if _, ok := repo.AsNoRowsAffectedError(err); ok {

			userProfile := entity.UserProfile{
				DisplayName: lineUserProfile.Name,
				PictureURL:  lineUserProfile.Picture,
			}

			up, err := r.u.Create(c.Request.Context(), userProfile)
			if err != nil {
				r.l.Error(err, "http - v1 - lineCallback")
				sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
				return
			}

			oauthDetail := entity.OAuthDetail{
				OAuthID:      lineUserProfile.Sub,
				UserID:       up.UserID,
				Provider:     "line",
				AccessToken:  tokens.AccessToken,
				RefreshToken: tokens.RefreshToken,
			}

			err = r.o.Create(c.Request.Context(), oauthDetail)
			if err != nil {
				r.l.Error(err, "http - v1 - lineCallback")
				sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
				return
			}

			// TODO: think about this way is good or not
			od.UserID = up.UserID
		} else {
			r.l.Error(err, "http - v1 - lineCallback")
			sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
			return
		}
	}

	// user profile already exists, update the refresh token from provider
	if od.OAuthID == lineUserProfile.Sub {
		err = r.o.UpdateRefreshToken(c.Request.Context(), lineUserProfile.Sub, tokens.RefreshToken)
		if err != nil {
			r.l.Error(err, "http - v1 - lineCallback")
			sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
			return
		}
	}

	session := sessions.Default(c)
	session.Set("userID", od.UserID)
	err = session.Save()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, r.domainUrl)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func (r *authRoutes) exchangeCodeForTokens(code string) (*TokenResponse, error) {

	tokenEndpoint := "https://api.line.me/oauth2/v2.1/token"

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {fmt.Sprintf("%s/api/v1/auth/line-callback", r.domainUrl)},
		"client_id":     {os.Getenv("LINE_CHANNEL_ID")},
		"client_secret": {os.Getenv("LINE_CHANNEL_SECRET")},
	}

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tr TokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}

	return &tr, nil
}

// logout godoc
//
// @Summary Logout
// @Description Logout
// @Tags Auth
// @Success 200 {string} string	"ok"
// @Router /auth/logout [get]
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
