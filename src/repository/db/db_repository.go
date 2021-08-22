package db

import (
	"errors"

	"github.com/gocql/gocql"
	"github.com/shawnzxx/bookstore_oauth-api/src/clients/cassandra"
	"github.com/shawnzxx/bookstore_oauth-api/src/domain/access_token"
	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(access_token.AccessToken) *rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *rest_errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *rest_errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("no access token found for giving id")
		}
		return nil, rest_errors.NewInternalServerError("error when trying to get access token", errors.New("database error"))
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserId,
		at.ClientId,
		at.Expires,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError("error when trying to create access token", errors.New("database error"))
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError("error when trying to update access token expire time", errors.New("database error"))
	}
	return nil
}
