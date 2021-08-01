package db

import (
	"github.com/shawnzxx/bookstore_oauth-api/src/clients/cassandra"
	"github.com/shawnzxx/bookstore_oauth-api/src/domain/access_token"
	"github.com/shawnzxx/bookstore_oauth-api/src/utils/errors"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	session, err := cassandra.GetSession()
	if err != nil {
		// if can not connect to db we stop the process
		panic(err)
	}
	defer session.Close()
	//TODO implement get access token from DB
	return nil, errors.NewInternalServerError("database not implemented yet")
}
