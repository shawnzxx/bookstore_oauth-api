package access_token

import (
	"strings"

	"github.com/shawnzxx/bookstore_oauth-api/src/domain/access_token"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/db"
	"github.com/shawnzxx/bookstore_oauth-api/src/repository/rest"
	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *rest_errors.RestErr
}

type service struct {
	restUserRepo rest.RestUserRepository
	dbRepo       db.DbRepository
}

func NewService(restUserRepo rest.RestUserRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUserRepo: restUserRepo,
		dbRepo:       dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, *rest_errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *rest_errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	//TODO: Support both grant types: client_credentials and password

	// Below is password grant_type flow
	// Authenticate the user against the Users API
	user, err := s.restUserRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	//Generate a new access token:
	at := access_token.GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) *rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
