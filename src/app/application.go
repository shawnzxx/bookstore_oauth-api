package app

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_oauth-api/src/clients/cassandra"
	"github.com/shawnzxx/bookstore_oauth-api/src/http"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/db"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/rest"
	"github.com/shawnzxx/bookstore_oauth-api/src/services/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {
	// When app bootup we want to test db connection first, if can not connect we want to stop process
	_ = cassandra.GetSession()

	atService := access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository())
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	router.Run(":8080")
}
