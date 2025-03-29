package exception

import (
	"farm-scurity/domain/web"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandle(gin *gin.Context, err interface{}) {
	if badRequestError(gin, err) {
		return
	} else if notFoundError(gin, err) {
		return
	}

	internalServerError(gin, err)
}

func notFoundError(gin *gin.Context, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		gin.JSON(http.StatusNotFound, web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: exception.error,
		})
		return true
	}
	return false
}

func badRequestError(gin *gin.Context, err interface{}) bool {
	exception, ok := err.(BadRequestError)
	if ok {
		gin.JSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: exception.error,
		})
		return true
	}
	return false
}

func internalServerError(gin *gin.Context, err interface{}) {
	fmt.Println(err)
	gin.JSON(http.StatusInternalServerError, web.ErrorResponse{
		Code:    http.StatusInternalServerError,
		Status:  "Internal Server Error",
		Message: err.(error).Error(),
	})
}
