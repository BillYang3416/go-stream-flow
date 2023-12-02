package v1

import (
	"net/http"

	"github.com/bgg/go-file-gate/internal/entity"
	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/internal/usecase/repo"
	"github.com/bgg/go-file-gate/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userProfileRoutes struct {
	u usecase.UserProfile
	l logger.Logger
}

func newUserProfileRoutes(handler *gin.RouterGroup, u usecase.UserProfile, l logger.Logger) {

	r := &userProfileRoutes{u, l}

	h := handler.Group("/user-profiles")
	{
		h.Use(CheckSessionMiddleware())
		h.POST("/", r.create)
		h.GET("/:userId", r.get)
	}
}

type createUserProfileRequest struct {
	UserID      string `json:"userId" example:"U1234567890abcdef1234567890abcdef" binding:"required"`
	DisplayName string `json:"displayName" example:"John Doe" binding:"required"`
	PictureURL  string `json:"pictureUrl" example:"https://example.com/picture.jpg" binding:"required,url"`
	AccessToken string `json:"accessToken" example:"1234567890abcdef1234567890abcdef" binding:"required"`
}

type userProfileResponse struct {
	UserID      string `json:"userId" example:"U1234567890abcdef1234567890abcdef"`
	DisplayName string `json:"displayName" example:"John Doe"`
	PictureURL  string `json:"pictureUrl" example:"https://example.com/picture.jpg"`
}

// Create Profile godoc
//
//	@Summary		Create user profile
//	@Description	Create user profile
//	@Tags			User Profile
//	@Accept			json
//	@Produce		json
//	@Param			createUserProfileRequest	body		createUserProfileRequest	true	"user profile information"
//	@Success		200						{object}	userProfileResponse
//	@Failure		400						{object}	errorResponse
//	@Router			/user-profiles [post]
func (r *userProfileRoutes) create(c *gin.Context) {
	var request createUserProfileRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - create")
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendValidationErrorResponse(c, validationErrs)
		} else {
			sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		}
		return
	}

	userProfile, err := r.u.Create(
		c.Request.Context(),
		entity.UserProfile{
			UserID:      request.UserID,
			DisplayName: request.DisplayName,
			PictureURL:  request.PictureURL,
			AccessToken: request.AccessToken,
		},
	)

	if err != nil {
		r.l.Error(err, "http - v1 - create")
		if repo.IsUniqueConstraintError(err) {
			sendErrorResponse(c, http.StatusConflict, "a user profile with the same user id already exists")
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
		}
		return
	}

	response := userProfileResponse{
		UserID:      userProfile.UserID,
		DisplayName: userProfile.DisplayName,
		PictureURL:  userProfile.PictureURL,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get user profile
//	@Tags			User Profile
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	string	true	"user id"
//	@Success		200						{object}	userProfileResponse
//	@Failure		400						{object}	errorResponse
//	@Router			/user-profiles/{userId} [get]
func (r *userProfileRoutes) get(c *gin.Context) {
	userId := c.Param("userId")

	userProfile, err := r.u.GetByID(c.Request.Context(), userId)

	if err != nil {
		r.l.Error(err, "http - v1 - find")
		sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
		return
	}

	response := userProfileResponse{
		UserID:      userProfile.UserID,
		DisplayName: userProfile.DisplayName,
		PictureURL:  userProfile.PictureURL,
	}

	c.JSON(http.StatusOK, response)
}
