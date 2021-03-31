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
	"time"
)

// Specification: User API
func User(newServer TestServer, t *testing.T) {
	srv := newServer(secret)
	validBearerToken := keyValue{key: "Authorization", value: "Bearer " + srv.NewJWT(15*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated")}
	expiredBearerToken := keyValue{key: "Authorization", value: "Bearer " + srv.NewJWT(0*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated")}

	// When given a new server
	// And the request is GET /api/user with no bearer token
	srv = newServer(secret)
	req := request("GET", "/api/user", nil)
	// Then executing the request should fail with status of 401 (not authorized)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}

	// When given a new server
	// And the request is GET /api/user with a valid bearer token
	srv = newServer(secret)
	req = httptest.NewRequest("GET", "/api/user", nil)
	req = request("GET", "/api/user", nil, contentType, validBearerToken)
	// Then executing the request should success with status of 200 (OK)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusOK, w.Code)
	}

	// When given a new server
	// And the request is GET /api/user with an expired bearer token
	srv = newServer(secret)
	req = httptest.NewRequest("GET", "/api/user", nil)
	req = request("GET", "/api/user", nil, contentType, expiredBearerToken)
	// Then executing the request should fail with status of 401 (not authorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}
}

func UserX(newServer TestServer, t *testing.T) {
	srv := newServer(secret)
	validBearerToken := keyValue{key: "Authorization", value: "Bearer " + srv.NewJWT(15*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated")}
	expiredBearerToken := keyValue{key: "Authorization", value: "Bearer " + srv.NewJWT(0*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated")}

	// Specification: User API

	// When given a new Server
	srv = newServer(secret)
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
	srv = newServer(secret)
	req = request("GET", "/api/user", nil)
	// Then executing the request should fail with status of 401 (not authorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}

	// When given a new Server
	// And the request is GET /api/user with a valid bearer token
	srv = newServer(secret)
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
	srv = newServer(secret)
	req = httptest.NewRequest("GET", "/api/user", nil)
	req = request("GET", "/api/user", nil, contentType, expiredBearerToken)
	// Then executing the request should fail with status of 401 (not authorized)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}
}
