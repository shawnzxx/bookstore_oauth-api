package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
	"time"
)

const (
	DBHost = "DB_HOST"
)

var (
	session *gocql.Session
	host    = os.Getenv(DBHost)
)

func init() {
	// Connect to Cassandra cluster
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = "oauth"

	var err error
	retryCount := 30
	for {
		session, err = cluster.CreateSession()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish connection to host %s", host)
			}
			log.Printf(fmt.Sprintf("Could not connect to database. Wait 5 seconds. %d retries left...", retryCount))
			retryCount--
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

// re-use same private session variable and return outside
func GetSession() *gocql.Session {
	return session
}
