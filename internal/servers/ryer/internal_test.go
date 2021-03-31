/*
 * conduit - current practices for Go web servers
 *
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"github.com/mdhender/conduit/internal/tests"
	"github.com/mdhender/conduit/internal/way"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	testServer := func(secret string) tests.Server {
		srv := &Server{
			db:           memory.New(),
			dtfmt:        "2006-01-02T15:04:05.99999999Z",
			router:       way.NewRouter(),
			tokenFactory: jwt.NewFactory(secret),
		}
		srv.Handler = srv.router
		srv.routes()
		return srv
	}
	tests.User(testServer, t)
}

func TestRegistration(t *testing.T) {
	// Specification: Registration API

	// When given a new Server
	// And the request is POST /api/users
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a NewUserRequest with the values
	//   { "user": { "username": "Jacob", "email": "jake@jake.jake", "password": "jakejake" } }
	srv := newServer()
	newUser := conduit.NewUser{Username: "Jacob", Email: "jake@jake.jake", Password: "jakejake"}
	req := request("POST", "/api/users", conduit.NewUserRequest{User: newUser}, contentType)
	// Then executing the request should succeed with status of 200 (ok)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusOK, http.StatusText(http.StatusOK), w.Code, http.StatusText(w.Code))
	} else {
		// And return a valid User
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("api: %q %q response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if userResponse.User.Email != newUser.Email {
				t.Errorf("api: %q %q email expected %q: got %q\n", req.Method, req.URL.Path, newUser.Email, userResponse.User.Username)
			}
			if userResponse.User.Username != newUser.Username {
				t.Errorf("api: %q %q username expected %q: got %q\n", req.Method, req.URL.Path, newUser.Username, userResponse.User.Username)
			}
		}
	}

	// When given the prior Server
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a NewUserRequest with the values
	//   { "user": { "username": "Jacob", "email": "jake@jake.jake", "password": "jakejake" } }
	req = request("POST", "/api/users", conduit.NewUserRequest{User: newUser}, contentType)
	// Then executing the request should fail with status of 422 (unprocessable entity)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), w.Code, http.StatusText(w.Code))
	}

	// When given the prior Server
	// And the request content type header is "text/plain"
	// And the request body is a NewUserRequest with the values
	//   { "user": { "username": "Jacob", "email": "jake@jake.jake", "password": "jakejake" } }
	srv = newServer()
	req = request("POST", "/api/users", conduit.NewUserRequest{User: newUser}, keyValue{key: "Content-Type", value: "text/plain"})
	// Then executing the request should fail with status of 422 (unprocessable entity)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType), w.Code, http.StatusText(w.Code))
	}
}

func TestAuthentication(t *testing.T) {
	// Specification: Authentication API

	// When given a new server
	srv := newServer()
	// And the user with e-mail "jake@jake.jake" and password "jakejake" has been added
	srv.ServeHTTP(httptest.NewRecorder(), request("POST", "/api/users", conduit.NewUserRequest{User: conduit.NewUser{Username: "Jacob", Email: "jake@jake.jake", Password: "jakejake"}}, contentType))
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a LoginUserRequest with the values
	//   { "user": { "email": "jake@jake.jake", "password": "jakejake" } }
	loginUser := conduit.LoginUser{Email: "jake@jake.jake", Password: "jakejake"}
	req := request("POST", "/api/users/login", conduit.LoginUserRequest{User: loginUser}, contentType)
	// Then executing the request should succeed with status of 200 (ok)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusOK, http.StatusText(http.StatusOK), w.Code, http.StatusText(w.Code))
	} else {
		// And return a valid User with the e-mail address of "jake@jake.jake"
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("api: %q %q response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if userResponse.User.Email != loginUser.Email {
				t.Errorf("api: %q %q email expected %q: got %q\n", req.Method, req.URL.Path, loginUser.Email, userResponse.User.Username)
			}
		}
	}

	// When given the prior server
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a LoginUserRequest with the values
	//   { "user": { "email": "jake@jake.jake", "password": "fakefake" } }
	loginUser = conduit.LoginUser{Email: "jake@jake.jake", Password: "fakefake"}
	req = request("POST", "/api/users/login", conduit.LoginUserRequest{User: loginUser}, contentType)
	// Then executing the request should fail with status of 401 (unauthorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), w.Code, http.StatusText(w.Code))
	}
}

func TestUser(t *testing.T) {
	srv := newServer()
	validBearerToken := keyValue{key: "Authorization", value: "Bearer " + srv.tokenFactory.NewToken(15*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated")}
	expiredBearerToken := keyValue{key: "Authorization", value: "Bearer " + srv.tokenFactory.NewToken(0*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated")}

	// Specification: User API

	// When given a new Server
	srv = newServer()
	// And the user with username "Jacob," e-mail "jake@jake.jake," and password "jakejake" has been added
	srv.ServeHTTP(httptest.NewRecorder(), request("POST", "/api/users", conduit.NewUserRequest{User: conduit.NewUser{Username: "Jacob", Email: "jake@jake.jake", Password: "jakejake"}}, contentType))
	// And the request is GET /api/user
	// And the request content type header is "application/json; charset=utf-8"
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	req := request("GET", "/api/user", nil, contentType, validBearerToken)
	// Then executing the request should succeed with status of 200 (ok)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusOK, http.StatusText(http.StatusOK), w.Code, http.StatusText(w.Code))
	} else {
		// And return a valid user with the email of "jake@jake.jake"
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("api: %q %q response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if userResponse.User.Email != "jake@jake.jake" {
				t.Errorf("api: %q %q email expected %q: got %q\n", req.Method, req.URL.Path, "jake@jake.jake", userResponse.User.Email)
			}
		}
	}

	// When given the prior Server
	// And the request is PUT /api/user
	// And the request content type header is "application/json; charset=utf-8"
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// And the request body is an UpdateUserRequest with the values
	//   { "user":{ "email": "jake@jake.jake", "bio": "I like to skateboard", "image": "https://i.stack.imgur.com/xHWG8.jpg" } }
	email, bio, image := "jake@jake.jake", "I like to skateboard", "https://i.stack.imgur.com/xHWG8.jpg"
	updateUser := conduit.UpdateUser{Email: &email, Bio: &bio, Image: &image}
	req = request("PUT", "/api/user", conduit.UpdateUserRequest{User: updateUser}, contentType, validBearerToken)
	// Then executing the request should succeed with status of 200 (ok)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusOK, http.StatusText(http.StatusOK), w.Code, http.StatusText(w.Code))
	} else {
		// And return a valid user with the username of "Jacob" with the updated fields
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("api: %q %q response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if userResponse.User.Username != "Jacob" {
				t.Errorf("api: %q %q username expected %q: got %q\n", req.Method, req.URL.Path, "Jacob", userResponse.User.Username)
			}
			if userResponse.User.Email != email {
				t.Errorf("api: %q %q email expected %q: got %q\n", req.Method, req.URL.Path, email, userResponse.User.Email)
			}
			if userResponse.User.Bio == nil {
				t.Errorf("api: %q %q bio expected %q: got nil\n", req.Method, req.URL.Path, bio)
			} else if *userResponse.User.Bio != bio {
				t.Errorf("api: %q %q bio expected %q: got %q\n", req.Method, req.URL.Path, bio, *userResponse.User.Bio)
			}
			if userResponse.User.Image == nil {
				t.Errorf("api: %q %q image expected %q: got nil\n", req.Method, req.URL.Path, image)
			} else if *userResponse.User.Image != image {
				t.Errorf("api: %q %q image expected %q: got %q\n", req.Method, req.URL.Path, image, *userResponse.User.Image)
			}
		}
	}

	// When given the prior Server
	// And the request is GET /api/user
	// And the request content type header is "application/json; charset=utf-8"
	// And the request does not include a valid bearer token
	// Then executing the request should fail with status of 401 (not authorized)
	t.Errorf("!implemented")

	// When given the prior Server
	// And the request is PUT /api/user
	// And the request content type header is "application/json; charset=utf-8"
	// And the request does not include a valid bearer token
	// And the request body is an UpdateUserRequest with the values
	//   { "user":{ "email": "jake@jake.jake", "bio": "Change is good" } }
	// Then executing the request should fail with status of 401 (not authorized)
	t.Errorf("!implemented")

	// When given a new Server
	// And the request is GET /api/user with no bearer token
	srv = newServer()
	req = request("GET", "/api/user", nil)
	// Then executing the request should fail with status of 401 (not authorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}

	// When given a new Server
	// And the request is GET /api/user with a valid bearer token
	srv = newServer()
	req = httptest.NewRequest("GET", "/api/user", nil)
	req = request("GET", "/api/user", nil, contentType, validBearerToken)
	// Then executing the request should success with status of 200 (ok)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusOK, w.Code)
	}

	// When given a new Server
	// And the request is GET /api/user with an expired bearer token
	srv = newServer()
	req = httptest.NewRequest("GET", "/api/user", nil)
	req = request("GET", "/api/user", nil, contentType, expiredBearerToken)
	// Then executing the request should fail with status of 401 (not authorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}
}

func TestProfile(t *testing.T) {
	// Specification: Profile API

	// When given a new Server
	// And the user with username "Jacob," e-mail "jake@jake.jake," and password "jakejake" has been added
	// And the user with username "Anne," e-mail "anne@anne.anne," and password "anneanne" has been added
	// And the request is GET /api/profiles/Anne
	// And the request content type header is "application/json; charset=utf-8"
	// Then executing the request should succeed with status of 200 (ok)
	// And return a valid Profile with the username of "Anne"
	t.Errorf("!implemented")

	// When given the prior Server
	// And the request is POST /api/profiles/:username/follow
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// Then executing the request should succeed with status of 200 (ok)
	// And return a valid Profile with the username of "Anne"
	t.Errorf("!implemented")

	// When given the prior Server
	// And the request is DELETE /api/profiles/Anne/follow
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// Then executing the request should succeed with status of 200 (ok)
	// And return a valid Profile with the username of "Anne"
	t.Errorf("!implemented")
}

func TestUsers(t *testing.T) {
	validNewUser := conduit.NewUser{Email: "jake@jake.jake", Username: "jake", Password: "jakejake"}
	validUserRequest := conduit.NewUserRequest{User: validNewUser}

	// Specification: Users API

	// When given a new Server
	// And the request is POST /api/users with a valid user request with no Content-Type header
	srv := newServer()
	req := request("POST", "/api/users", validUserRequest)
	// Then executing the request should succeed with status of 200 (OK)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnsupportedMediaType, w.Code)
	}

	// When given a new Server
	// And the request is POST /api/users with a valid user request and Content-Type
	srv = newServer()
	req = request("POST", "/api/users", validUserRequest, contentType)
	// Then executing the request should succeed with status of 200 (OK)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusOK, w.Code)
	} else {
		// And return a valid User
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("api: %q %q response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if userResponse.User.Email != validNewUser.Email {
				t.Errorf("api: %q %q email expected %q: got %q\n", req.Method, req.URL.Path, validNewUser.Email, userResponse.User.Username)
			}
			if userResponse.User.Username != validNewUser.Username {
				t.Errorf("api: %q %q username expected %q: got %q\n", req.Method, req.URL.Path, validNewUser.Username, userResponse.User.Username)
			}
		}
	}

	// When given a new Server
	// And the test user Jacob is loaded
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a LoginUserRequest with the values
	//   { "user": { "email": "jake@jake.jake", "password": "jakejake" } }
	srv = newServer()
	srv.ServeHTTP(httptest.NewRecorder(), request("POST", "/api/users", validUserRequest, contentType))
	loginUser := conduit.LoginUser{
		Email:    "jake@jake.jake",
		Password: "jakejake",
	}
	loginUserRequest := conduit.LoginUserRequest{User: loginUser}
	req = request("POST", "/api/users/login", loginUserRequest, contentType)
	// Then executing the request should succeed with status of 200 (ok)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusOK, http.StatusText(http.StatusOK), w.Code, http.StatusText(w.Code))
	} else {
		// And return a the test user Jacob
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("api: %q %q response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if userResponse.User.Email != loginUser.Email {
				t.Errorf("api: %q %q email expected %q: got %q\n", req.Method, req.URL.Path, loginUser.Email, userResponse.User.Username)
			}
		}
	}

	// When given the prior Server
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a LoginUserRequest with the values
	//   { "user": { "email": "jake@jake.jake", "password": "fakejake" } }
	loginUser.Password = "fakejake"
	loginUserRequest = conduit.LoginUserRequest{User: loginUser}
	req = request("POST", "/api/users/login", loginUserRequest, contentType)
	// Then executing the request should fail with status of 401 (not authorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), w.Code, http.StatusText(w.Code))
	}
}

type keyValue struct {
	key, value string
}

var contentType = keyValue{"Content-Type", "application/json; charset=utf-8"}

func request(method, target string, v interface{}, keys ...keyValue) *http.Request {
	var body io.Reader
	if v != nil {
		buf, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Sprintf("assert(err != %+v)", err))
		}
		body = bytes.NewReader(buf)
	}
	req := httptest.NewRequest(method, target, body)
	for _, kv := range keys {
		req.Header.Set(kv.key, kv.value)
	}
	return req
}

func fetch(body io.Reader, data interface{}) error {
	if dec := json.NewDecoder(body); dec == nil {
		return fmt.Errorf("failed to create decoder")
	} else if err := dec.Decode(&data); err != nil {
		return err
	} else if err = dec.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("body must only contain a single JSON object")
	}

	return nil
}

func jsonReader(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	buf, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("assert(err != %+v)", err))
	}
	return bytes.NewReader(buf)
}

func newServer() *Server {
	srv := &Server{
		db:           memory.New(),
		dtfmt:        "2006-01-02T15:04:05.99999999Z",
		router:       way.NewRouter(),
		tokenFactory: jwt.NewFactory("salt+pepper"),
	}
	srv.MaxHeaderBytes = 1 << 20
	srv.Handler = srv.router
	srv.routes()
	return srv
}
