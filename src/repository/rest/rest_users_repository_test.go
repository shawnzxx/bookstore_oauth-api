package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/federicoleon/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

// TestMain is this package test entry point, each package have only one TestMain
// When use "go test" command, TestMain function will call below each test cases
func TestMain(m *testing.M) {
	fmt.Println("about to start test cases...")
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email1@gmail.com","password":"my_password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`, //return invalid error interfaces
	})

	repository := restUserRepository{}
	user, err := repository.LoginUser("email1@gmail.com", "my_password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "restclient request timeout when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email2@gmail.com","password":"my_password"}`,
		RespHTTPCode: http.StatusNotFound,
		//status return as string instead of number, simulate wrong rest error interface
		RespBody: `{"message": "invalid login credentials", "status": "404", "error": "not_found"}`,
	})

	repository := restUserRepository{}
	user, err := repository.LoginUser("email2@gmail.com", "my_password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentails(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email3@gmail.com","password":"my_password"}`,
		RespHTTPCode: http.StatusNotFound,
		//this return 404 not found error and message as "invalid login credentials"
		RespBody: `{"message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})

	repository := restUserRepository{}
	user, err := repository.LoginUser("email3@gmail.com", "my_password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email4@gmail.com","password":"my_password"}`,
		RespHTTPCode: http.StatusOK,
		//we simulate return wrong json struct, id should be int64 instead of string
		RespBody: `{"id":"1","first_name":"Xiaoxiao","last_name":"Zhang","email":"shawnzhang.dev@gmail.com"}`,
	})

	repository := restUserRepository{}
	user, err := repository.LoginUser("email4@gmail.com", "my_password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal users login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email5@gmail.com","password":"my_password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":1,"first_name":"Xiaoxiao","last_name":"Zhang","email":"shawnzhang.dev@gmail.com"}`,
	})

	repository := restUserRepository{}
	user, err := repository.LoginUser("email5@gmail.com", "my_password")

	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "Xiaoxiao", user.FirstName)
	assert.EqualValues(t, "Zhang", user.LastName)
	assert.EqualValues(t, "shawnzhang.dev@gmail.com", user.Email)
}
