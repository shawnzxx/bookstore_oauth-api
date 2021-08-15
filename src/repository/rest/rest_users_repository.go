package rest

import (
	"encoding/json"
	"time"

	"github.com/federicoleon/golang-restclient/rest"
	"github.com/shawnzxx/bookstore_oauth-api/src/domain/users"
	"github.com/shawnzxx/bookstore_oauth-api/src/utils/errors"
)

var (
	userRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 100 * time.Microsecond,
	}
)

type RestUserRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type restUserRepository struct {
}

func NewRestUsersRepository() RestUserRepository {
	return &restUserRepository{}
}

func (r *restUserRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := userRestClient.Post("/users/login", request)

	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerError("restclient request timeout when trying to login user")
	}

	// means error happened
	if response.StatusCode > 299 {
		var restErr errors.RestErr
		// since we use same rest error struct for both auth and user service
		// if can not unmarshal response means someone changed the struct
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		// error struct no change return real response error
		return nil, &restErr
	}

	var user users.User
	err := json.Unmarshal(response.Bytes(), &user)
	if err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users login response")
	}
	return &user, nil
}
