package app

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_oauth-api/src/clients/cassandra"
	"github.com/shawnzxx/bookstore_oauth-api/src/domain/access_token"
	"github.com/shawnzxx/bookstore_oauth-api/src/http"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	// When app bootup we want to test db connection first, if can not connect we want to stop process
	_ = cassandra.GetSession()

	dbRepository := db.NewRepository()
	atService := access_token.NewService(dbRepository)
	atHandler := http.NewHandler(atService)

	router.GET("/ouath/access_token/:access_token_id", atHandler.GetById)
	router.POST("/ouath/access_token", atHandler.Create)
	router.Run(":8080")
}
