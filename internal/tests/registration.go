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

func Registration(newServer TestServer, t *testing.T) {
	// Specification: Registration API

	// When given a new Server
	// And the request is POST /api/users
	// And the request content type header is "application/json; charset=utf-8"
	// And the request body is a NewUserRequest with the values
	//   { "user": { "username": "Jacob", "email": "jake@jake.jake", "password": "jakejake" } }
	srv := newServer(secret)
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
	srv = newServer(secret)
	req = request("POST", "/api/users", conduit.NewUserRequest{User: newUser}, keyValue{key: "Content-Type", value: "text/plain"})
	// Then executing the request should fail with status of 422 (unprocessable entity)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("api: %q %q expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType), w.Code, http.StatusText(w.Code))
	}
}
