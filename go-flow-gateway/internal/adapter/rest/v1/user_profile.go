package v1

import (
	"net/http"
	"strconv"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/internal/usecase/apperrors"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userProfileRoutes struct {
	userProfile usecase.UserProfile
	logger      logger.Logger
}

func NewUserProfileRoutes(handler *gin.RouterGroup, u usecase.UserProfile, l logger.Logger) {

	r := &userProfileRoutes{u, l}

	h := handler.Group("/user-profiles")
	{
		h.Use(CheckSessionMiddleware())
		h.POST("/", r.create)
		h.GET("/:userId", r.get)
	}
}

type createUserProfileRequest struct {
	DisplayName string `json:"displayName" example:"John Doe" binding:"required"`
	PictureURL  string `json:"pictureUrl" example:"https://example.com/picture.jpg" binding:"required,url"`
}

type userProfileResponse struct {
	DisplayName string `json:"displayName" example:"John Doe"`
	PictureURL  string `json:"pictureUrl" example:"https://example.com/picture.jpg"`
}

// Create Profile godoc
//
//	@Summary		Create UserProfileRoutes
//	@Description	Create UserProfileRoutes
//	@Tags			UserProfileRoutes
//	@Accept			json
//	@Produce		json
//	@Param			createUserProfileRequest	body		createUserProfileRequest	true	"UserProfileRoutes information"
//	@Success		200							{object}	userProfileResponse
//	@Failure		400							{object}	errorResponse
//	@Router			/user-profiles [post]
func (r *userProfileRoutes) create(c *gin.Context) {
	var request createUserProfileRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.logger.Error("UserProfileRoutes - Create : invalid request body", err)
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			sendValidationErrorResponse(c, validationErrs)
		} else {
			sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		}
		return
	}

	userProfile, err := r.userProfile.Create(
		c.Request.Context(),
		entity.UserProfile{
			DisplayName: request.DisplayName,
			PictureURL:  request.PictureURL,
		},
	)

	if err != nil {
		r.logger.Error("UserProfileRoutes - Create : failed to create UserProfileRoutes", err)
		if apperrors.IsUniqueConstraintError(err) {
			sendErrorResponse(c, http.StatusConflict, "a UserProfileRoutes with the same user id already exists")
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
		}
		return
	}

	response := userProfileResponse{
		DisplayName: userProfile.DisplayName,
		PictureURL:  userProfile.PictureURL,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserProfile godoc
//
//	@Summary		Get UserProfileRoutes
//	@Description	Get UserProfileRoutes
//	@Tags			UserProfileRoutes
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"user id"
//	@Success		200		{object}	userProfileResponse
//	@Failure		400		{object}	errorResponse
//	@Router			/user-profiles/{userId} [get]
func (r *userProfileRoutes) get(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		r.logger.Error("UserProfileRoutes - Get : invalid user id", err)
		sendErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	userProfile, err := r.userProfile.GetByID(c.Request.Context(), userId)

	if err != nil {
		r.logger.Error("UserProfileRoutes - Get : failed to get UserProfileRoutes", err)
		sendErrorResponse(c, http.StatusInternalServerError, "internal server problems")
		return
	}

	response := userProfileResponse{
		DisplayName: userProfile.DisplayName,
		PictureURL:  userProfile.PictureURL,
	}

	c.JSON(http.StatusOK, response)
}
