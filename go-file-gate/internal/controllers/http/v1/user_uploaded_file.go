package v1

import (
	"io"
	"net/http"

	"github.com/bgg/go-file-gate/internal/entity"
	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userUploadedFileRoutes struct {
	u usecase.UserUploadedFile
	l logger.Logger
}

func NewUserUploadedFileRoutes(handler *gin.RouterGroup, u usecase.UserUploadedFile, l logger.Logger) {

	r := &userUploadedFileRoutes{u, l}

	h := handler.Group("/user-uploaded-files")
	{
		h.Use(CheckSessionMiddleware())
		h.POST("/", r.create)
	}
}

type createUserUploadedFileRequest struct {
	EmailRecipient string `form:"emailRecipient" example:"johndoe@email.com"  binding:"required,email"`
}

// create user uploaded file godoc
//
// @Summary		Create user uploaded file
// @Description	Create user uploaded file
// @Tags			User Uploaded File
// @Accept			multipart/form-data
// @Produce		json
// @Param			emailRecipient	formData	string	true	"email recipient"
// @Param			file			formData	file	true	"file"
// @Success		204
// @Failure		400						{object}	errorResponse
// @Router			/user-uploaded-files [post]
func (r *userUploadedFileRoutes) create(c *gin.Context) {
	var request createUserUploadedFileRequest
	if err := c.ShouldBind(&request); err != nil {
		r.l.Error(err, "http - v1 - create")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	file, err := c.FormFile("file")
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	uploadedFile, err := file.Open()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer uploadedFile.Close()

	fileContent, err := io.ReadAll(uploadedFile)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	session := sessions.Default(c)
	userID, exists := session.Get("userID").(string)
	if !exists {
		r.l.Error(err, "http - v1 - create")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	_, err = r.u.Create(c.Request.Context(),
		entity.UserUploadedFile{
			UserID:         userID,
			EmailRecipient: request.EmailRecipient,
			Name:           file.Filename,
			Size:           file.Size,
			Content:        fileContent,
		})
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to create user uploaded file")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
