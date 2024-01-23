package v1

import (
	"io"
	"net/http"
	"strconv"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/logger"
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
		h.GET("/", r.getPaginatedFiles)
	}
}

type createUserUploadedFileRequest struct {
	EmailRecipient string `form:"emailRecipient" example:"johndoe@email.com"  binding:"required,email"`
}

// create user uploaded file godoc
//
//	@Summary		Create user uploaded file
//	@Description	Create user uploaded file
//	@Tags			User Uploaded File
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			emailRecipient	formData	string	true	"email recipient"
//	@Param			file			formData	file	true	"file"
//	@Success		204
//	@Failure		400	{object}	errorResponse
//	@Router			/user-uploaded-files [post]
func (r *userUploadedFileRoutes) create(c *gin.Context) {
	var request createUserUploadedFileRequest
	if err := c.ShouldBind(&request); err != nil {
		r.l.Error(err, "http - v1 - create")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
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
	userID, exists := session.Get("userID").(int)
	if !exists {
		r.l.Error(err, "http - v1 - create")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
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

type getPaginatedFilesResponse struct {
	Files        []entity.UserUploadedFile `json:"files"`
	TotalRecords int                       `json:"totalRecords"`
}

// get paginated files godoc
//
//	@Summary		Get paginated files
//	@Description	Get paginated files
//	@Tags			User Uploaded File
//	@Accept			json
//	@Produce		json
//	@Param			lastID	query		int	false	"last id"
//	@Param			limit	query		int	false	"limit"
//	@Success		200		{object}	getPaginatedFilesResponse
//	@Failure		400		{object}	errorResponse
//	@Router			/user-uploaded-files [get]
func (r *userUploadedFileRoutes) getPaginatedFiles(c *gin.Context) {
	lastID, err := strconv.Atoi(c.Query("lastID"))
	if err != nil {
		r.l.Error(err, "http - v1 - getPaginatedFiles")
		sendErrorResponse(c, http.StatusBadRequest, "invalid lastID query parameter")
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		r.l.Error(err, "http - v1 - getPaginatedFiles")
		sendErrorResponse(c, http.StatusBadRequest, "invalid limit query parameter")
		return
	}

	session := sessions.Default(c)
	userID, exists := session.Get("userID").(int)
	if !exists {
		r.l.Error(err, "http - v1 - getPaginatedFiles")
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	files, totalRecords, err := r.u.GetPaginatedFiles(c.Request.Context(), lastID, userID, limit)
	if err != nil {
		r.l.Error(err, "http - v1 - getPaginatedFiles")
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to get paginated files")
		return
	}

	c.JSON(http.StatusOK, getPaginatedFilesResponse{Files: files, TotalRecords: totalRecords})
}
