package app

import (
	"github.com/gin-gonic/gin"
	"github.com/shawnzxx/bookstore_oauth-api/src/clients/cassandra"
	"github.com/shawnzxx/bookstore_oauth-api/src/http"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/db"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/rest"
	"github.com/shawnzxx/bookstore_oauth-api/src/services/access_token"
	"github.com/shawnzxx/bookstore_utils-go/app_logger"
	"net"
	"os"
)

var (
	logger          = app_logger.GetLogger()
	router          = gin.Default()
	env, ipv4, port string
)

func StartApplication() {
	printOutServiceInfo()
	// When app boot up we want to test db connection first, if can not connect we want to stop process
	_ = cassandra.GetSession()

	atService := access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository())
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":" + port)
}

func printOutServiceInfo() {
	//get local or container's host name
	hostName, _ := os.Hostname()
	logger.Info("users api's hostname: %s", hostName)
	//print out service IPs
	ips, err := net.LookupHost(hostName)
	if err != nil {
		logger.Info("Can not find ips list for the host %v", err)
	}
	//print out service ip, env, port
	env = os.Getenv("ENV")
	ipv4 = ips[0]
	port = os.Getenv("PORT")
	logger.Info("bookstore_oauth-api running on %s environment, ip is %s, port is %s", env, ipv4, port)
}
