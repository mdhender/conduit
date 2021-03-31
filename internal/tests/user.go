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

	// Specification: User API

	// Given a new server
	// And the user with username "Jacob," e-mail "jake@jake.jake," and password "jakejake" has been added
	// And the request is GET /api/user
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid UserResponse with a valid User
	// And the User email should be "jake@jake.jake"
	srv = newServer(secret)
	srv.ServeHTTP(httptest.NewRecorder(), request("POST", "/api/users", conduit.NewUserRequest{User: conduit.NewUser{Username: "Jacob", Email: "jake@jake.jake", Password: "jakejake"}}, contentType))
	req := request("GET", "/api/user", nil, validBearerToken)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("user: %s %s response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := "jake@jake.jake"; userResponse.User.Email != expected {
				t.Errorf("user: %s %s email expected %q: got %q\n", req.Method, req.URL.Path, expected, userResponse.User.Email)
			}
		}
	}

	// Given the prior server
	// And the request is GET /api/user
	// And the request does not include a bearer token
	// When we execute the request
	// Then the response should have a status of 401 (not authorized)
	req = request("GET", "/api/user", nil)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusUnauthorized; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	}

	// Given the prior server
	// And the request is GET /api/user
	// And the request includes a bearer token for an unauthenticated user
	// When we execute the request
	// Then the response should have a status of 401 (not authorized)
	req = request("GET", "/api/user", nil, keyValue{key: "Authorization", value: "Bearer " + srv.NewJWT(15*time.Second, 1, "Guest", "guest@guest.guest")})
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusUnauthorized; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	}

	// Given the prior server
	// And the request is GET /api/user
	// And the request includes a bearer token for a non-existent user
	// Then the response should have a status of 200 (ok)
	// And contain a valid UserResponse with a valid User
	// And the User Id should be 0
	// And the User email should be empty
	req = request("GET", "/api/user", nil, keyValue{key: "Authorization", value: "Bearer " + srv.NewJWT(15*time.Second, 0, "Guest", "guest@guest.guest", "authenticated")})
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("user: %s %s response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := 0; userResponse.User.Id != expected {
				t.Errorf("user: %s %s id expected %q: got %q\n", req.Method, req.URL.Path, expected, userResponse.User.Id)
			}
			if expected := ""; userResponse.User.Email != expected {
				t.Errorf("user: %s %s email expected %q: got %q\n", req.Method, req.URL.Path, expected, userResponse.User.Email)
			}
		}
	}

	// Given the prior server
	// And the request is GET /api/user
	// And the request includes an expired bearer token
	// When we execute the request
	// Then the response should have a status of 401 (not authorized)
	req = request("GET", "/api/user", nil, expiredBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusUnauthorized; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	}

	// Given the prior server
	// And the request is PUT /api/user
	// And the request content type header is "application/json; charset=utf-8"
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// And the request body is an UpdateUserRequest with the values
	//   { "user":{ "email": "jake@jake.jake", "bio": "I like to skateboard", "image": "https://i.stack.imgur.com/xHWG8.jpg" } }
	// Then the response should have a status of 200 (ok)
	// And contain a valid UserResponse with a valid User
	// And the User username should be "Jacob"
	// And the User email should be "jake@jake.jake"
	// And the User bio should be "I like to skateboard"
	// And the User image should be "https://i.stack.imgur.com/xHWG8.jpg"
	username, email, bio, image := "Jacob", "jake@jake.jake", "I like to skateboard", "https://i.stack.imgur.com/xHWG8.jpg"
	updateUser := conduit.UpdateUser{Email: &email, Bio: &bio, Image: &image}
	req = request("PUT", "/api/user", conduit.UpdateUserRequest{User: updateUser}, contentType, validBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var userResponse conduit.UserResponse
		if err := fetch(w.Result().Body, &userResponse); err != nil {
			t.Errorf("user: %s %s response did not contain valid UserResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := username; userResponse.User.Username != expected {
				t.Errorf("user: %s %s username expected %q: got %q\n", req.Method, req.URL.Path, expected, userResponse.User.Username)
			}
			if expected := email; userResponse.User.Email != expected {
				t.Errorf("user: %s %s email expected %q: got %q\n", req.Method, req.URL.Path, expected, userResponse.User.Email)
			}
			if expected := bio; userResponse.User.Bio == nil {
				t.Errorf("user: %s %s bio expected %q: got nil\n", req.Method, req.URL.Path, expected)
			} else if *userResponse.User.Bio != expected {
				t.Errorf("user: %s %s bio expected %q: got %q\n", req.Method, req.URL.Path, expected, *userResponse.User.Bio)
			}
			if expected := image; userResponse.User.Image == nil {
				t.Errorf("user: %s %s image expected %q: got nil\n", req.Method, req.URL.Path, expected)
			} else if *userResponse.User.Image != expected {
				t.Errorf("user: %s %s image expected %q: got %q\n", req.Method, req.URL.Path, expected, *userResponse.User.Image)
			}
		}
	}

	// Given the prior server
	// And the request is PUT /api/user
	// And the request content type header is "application/json; charset=utf-8"
	// And the request does not include a valid bearer token
	// And the request body is an UpdateUserRequest with the values
	//   { "user":{ "email": "jake@jake.jake", "bio": "Change is good" } }
	// When we execute the request
	// Then the response should have a status of 401 (not authorized)
	email, bio = "jake@jake.jake", "Change is good"
	updateUser = conduit.UpdateUser{Email: &email, Bio: &bio}
	req = request("PUT", "/api/user", conduit.UpdateUserRequest{User: updateUser}, contentType, expiredBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusUnauthorized; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	}
}
