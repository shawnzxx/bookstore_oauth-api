package app

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_oauth-api/src/domain/access_token"
	"github.com/shawnzxx/bookstore_oauth-api/src/http"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	dbRepository := db.NewRepository()
	atService := access_token.NewService(dbRepository)
	atHandler := http.NewHandler(atService)

	router.GET("/ouath/access_token/:access_token_id", atHandler.GetById)
	router.Run(":8080")
}
