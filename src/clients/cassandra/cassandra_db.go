package cassandra

import (
	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func init() {
	// Connect to Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "oauth"

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}

// re-use same private session variable and return outside
func GetSession() *gocql.Session {
	return session
}
