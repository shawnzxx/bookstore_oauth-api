package http

import (
	"encoding/json"
	"github.com/shawnzxx/bookstore_utils-go/app_logger"
	"net/http"

	"github.com/gin-gonic/gin"
	atDomain "github.com/shawnzxx/bookstore_oauth-api/src/domain/access_token"
	"github.com/shawnzxx/bookstore_oauth-api/src/services/access_token"
	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

var (
	logger = app_logger.GetLogger()
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

// GetById If you need to change to other http framework only handler params will change, domain layer and service layer no dependency
func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, restErr := handler.service.GetById(c.Param("access_token_id"))
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		errByte, _ := json.Marshal(restErr)
		logger.Error("restErr %v", string(errByte))
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var request atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		errByte, _ := json.Marshal(restErr)
		logger.Error("restErr %v", string(errByte))
		return
	}

	accessToken, restErr := handler.service.Create(request)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		errByte, _ := json.Marshal(restErr)
		logger.Error("restErr %v", string(errByte))
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}
