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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"github.com/mdhender/conduit/internal/way"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	srv := newServer()
	validBearerToken := keyValue{
		key:   "Authorization",
		value: "Bearer " + srv.tokenFactory.NewToken(15*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated"),
	}
	expiredBearerToken := keyValue{
		key:   "Authorization",
		value: "Bearer " + srv.tokenFactory.NewToken(0*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated"),
	}

	// Specification: User API

	// When given a new Server
	// And the request is GET /api/user with no bearer token
	srv = newServer()
	req := request("GET", "/api/user", nil)
	// Then executing the request should fail with status of 401 (not authorized)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}

	// When given a new Server
	// And the request is GET /api/user with a valid bearer token
	srv = newServer()
	req = httptest.NewRequest("GET", "/api/user", nil)
	req = request("GET", "/api/user", nil, contentType, validBearerToken)
	// Then executing the request should success with status of 200 (OK)
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
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a LoginUserRequest with the values
	//   { "user": { "email": "jake@jake.jake", "password": "jakejake" } }
	srv = newServer()
	loginUser := conduit.LoginUser{
		Email:    "jake@jake.jake",
		Password: "jakejake",
	}
	loginUserRequest := conduit.LoginUserRequest{User: loginUser}
	req = request("POST", "/api/users/login", loginUserRequest, contentType)
	// Then executing the request should fail with status of 401 (not authorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), w.Code, http.StatusText(w.Code))
	}

	// When given a new Server
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is
	//   { "user":{
	//	     "email": "jake@jake.jake",
	//		 "password": "jakejake" }}
	// Then executing the request should succeed with status of 200 (OK)
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
