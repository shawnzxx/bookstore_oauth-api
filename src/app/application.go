package app

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_oauth-api/src/clients/cassandra"
	"github.com/shawnzxx/bookstore_oauth-api/src/http"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/db"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/rest"
	"github.com/shawnzxx/bookstore_oauth-api/src/services/access_token"
	"github.com/shawnzxx/bookstore_utils-go/logger"
	"net"
	"os"
)

var (
	_logger         = logger.GetLogger()
	_router         = gin.Default()
	env, ipv4, port string
)

func StartApplication() {
	printOutServiceInfo()
	// When app boot up we want to test db connection first, if can not connect we want to stop process
	_ = cassandra.GetSession()

	atService := access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository())
	atHandler := http.NewHandler(atService)

	_router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	_router.POST("/oauth/access_token", atHandler.Create)

	_router.Run(":" + port)
}

func printOutServiceInfo() {
	//print out service Hostname
	containerHostname, _ := os.Hostname()
	_logger.Printf("users api's hostname: %s", containerHostname)
	//print out service IPs
	ips, err := net.LookupHost(containerHostname)
	_logger.Printf("LookupHost: %v, error: %v\n", ips, err)
	//print out service ip, env, port
	env = os.Getenv("ENV")
	ipv4 = ips[0]
	port = os.Getenv("PORT")
	_logger.Printf("users api running in environment: %s, ip: %s, port: %s", env, ipv4, port)
}
