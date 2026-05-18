package handler

import (
	"go-order-inventory/internal/apperror"
	"go-order-inventory/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error, fallbackCode int, fallbackMsg string) {
	if appErr, ok := apperror.FromError(err); ok {
		response.Fail(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
		return
	}
	response.Fail(c, http.StatusInternalServerError, fallbackCode, fallbackMsg)
}
