package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/models"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/runner/common"
	"github.com/gin-gonic/gin"
)

var ErrInternalServerError = errors.New("Something unexpected happened on the server")

func simpleError(err error) *models.Error {
	return &models.Error{Error: &models.ErrorBody{Message: err.Error()}}
}

var errStatusCode = map[error]int{
	models.ErrAppsNotFound:        http.StatusNotFound,
	models.ErrAppsAlreadyExists:   http.StatusConflict,
	models.ErrRoutesNotFound:      http.StatusNotFound,
	models.ErrRoutesAlreadyExists: http.StatusConflict,
}

func handleErrorResponse(c *gin.Context, err error) {
	ctx := c.MustGet("ctx").(context.Context)
	log := common.Logger(ctx)
	log.Error(err)

	if code, ok := errStatusCode[err]; ok {
		c.JSON(code, simpleError(err))
	} else {
		c.JSON(http.StatusInternalServerError, simpleError(err))
	}
}
