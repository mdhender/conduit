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

package tests

import (
	"github.com/mdhender/conduit/internal/conduit"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Users(newServer TestServer, t *testing.T) {
	validNewUser := conduit.NewUser{Email: "jake@jake.jake", Username: "jake", Password: "jakejake"}
	validUserRequest := conduit.NewUserRequest{User: validNewUser}

	// Specification: Users API

	// When given a new Server
	// And the request is POST /api/users with a valid user request with no Content-Type header
	srv := newServer(secret)
	req := request("POST", "/api/users", validUserRequest)
	// Then executing the request should succeed with status of 200 (OK)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnsupportedMediaType, w.Code)
	}

	// When given a new Server
	// And the request is POST /api/users with a valid user request and Content-Type
	srv = newServer(secret)
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
	srv = newServer(secret)
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
