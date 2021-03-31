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

func Authentication(newServer TestServer, t *testing.T) {
	// Specification: Authentication API

	// Given a new server
	// And the user with e-mail "jake@jake.jake" and password "jakejake" has been added
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a LoginUserRequest with the values
	//   { "user": { "email": "jake@jake.jake", "password": "jakejake" } }
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid UserResponse with a valid User
	// And the User e-mail address should be "jake@jake.jake"
	srv := newServer(secret)
	srv.ServeHTTP(httptest.NewRecorder(), request("POST", "/api/users", conduit.NewUserRequest{User: conduit.NewUser{Username: "Jacob", Email: "jake@jake.jake", Password: "jakejake"}}, contentType))
	loginUser := conduit.LoginUser{Email: "jake@jake.jake", Password: "jakejake"}
	req := request("POST", "/api/users/login", conduit.LoginUserRequest{User: loginUser}, contentType)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("authentication: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("authentication: %q %q response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if userResponse.User.Email != loginUser.Email {
				t.Errorf("authentication: %q %q email expected %q: got %q\n", req.Method, req.URL.Path, loginUser.Email, userResponse.User.Username)
			}
		}
	}

	// Given the prior server
	// And the request is POST /api/users/login
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a LoginUserRequest with the values
	//   { "user": { "email": "jake@jake.jake", "password": "fakefake" } }
	// When we execute the request
	// Then the response should have a status of 401 (not authorized)
	loginUser = conduit.LoginUser{Email: "jake@jake.jake", Password: "fakefake"}
	req = request("POST", "/api/users/login", conduit.LoginUserRequest{User: loginUser}, contentType)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusUnauthorized; w.Code != expected {
		t.Errorf("authentication: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	}
}
