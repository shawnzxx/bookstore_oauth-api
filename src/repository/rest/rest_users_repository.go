package rest

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/federicoleon/golang-restclient/rest"
	"github.com/shawnzxx/bookstore_oauth-api/src/domain/users"
	"github.com/shawnzxx/bookstore_utils-go/rest_errors"
)

var (
	userRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8081",
		Timeout: 100 * time.Microsecond,
	}
)

type RestUserRepository interface {
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
}

type restUserRepository struct {
}

func NewRestUsersRepository() RestUserRepository {
	return &restUserRepository{}
}

func (r *restUserRepository) LoginUser(email string, password string) (*users.User, *rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := userRestClient.Post("/users/login", request)

	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("rest client request error when trying to login user", errors.New("database error"))
	}

	// means error happened
	if response.StatusCode > 299 {
		var restErr rest_errors.RestErr
		// since we use same rest error struct for both auth and user service
		// if can not unmarshal response means someone changed the struct
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", errors.New("database error"))
		}
		// error struct no change return real response error
		return nil, &restErr
	}

	var user users.User
	err := json.Unmarshal(response.Bytes(), &user)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users login response", errors.New("database error"))
	}
	return &user, nil
}
