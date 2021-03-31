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

func Profile(newServer TestServer, t *testing.T) {
	srv := newServer(secret)
	validBearerToken := keyValue{key: "Authorization", value: "Bearer " + srv.NewJWT(15*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated")}

	// Specification: Profile API

	// Given a new server
	// And the user with username "Jacob," e-mail "jake@jake.jake," and password "jakejake" has been added
	// And the user with username "Anne," e-mail "anne@anne.anne," and password "anneanne" has been added
	// And the request is GET /api/profiles/Anne
	// And the request content type header is "application/json; charset=utf-8"
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid Profile
	// And the username should be "Anne"
	// And the following flag should be false
	srv = newServer(secret)
	newUser := conduit.NewUser{Username: "Jacob", Email: "jake@jake.jake", Password: "jakejake"}
	srv.ServeHTTP(httptest.NewRecorder(), request("POST", "/api/users", conduit.NewUserRequest{User: newUser}, contentType))
	newUser = conduit.NewUser{Username: "Anne", Email: "anne@anne.anne", Password: "anneanne"}
	srv.ServeHTTP(httptest.NewRecorder(), request("POST", "/api/users", conduit.NewUserRequest{User: newUser}, contentType))
	req := request("GET", "/api/profiles/Anne", nil, contentType)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var profileResponse conduit.ProfileResponse
		if err := fetch(w.Result().Body, &profileResponse); err != nil {
			t.Errorf("profile: %s %s response did not contain valid ProfileResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := "Anne"; profileResponse.Profile.Username != expected {
				t.Errorf("profile: %s %s username expected %q: got %q\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Username)
			}
			if expected := false; profileResponse.Profile.Following != expected {
				t.Errorf("profile: %s %s following expected %v: got %v\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Following)
			}
		}
	}

	// Given the prior Server
	// And the request is GET /api/profiles/Anne
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid Profile
	// And the username should be "Anne"
	// And the following flag should be false
	req = request("GET", "/api/profiles/Anne", nil, validBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var profileResponse conduit.ProfileResponse
		if err := fetch(w.Result().Body, &profileResponse); err != nil {
			t.Errorf("profile: %s %s response did not contain valid ProfileResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := "Anne"; profileResponse.Profile.Username != expected {
				t.Errorf("profile: %s %s username expected %q: got %q\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Username)
			}
			if expected := false; profileResponse.Profile.Following != expected {
				t.Errorf("profile: %s %s following expected %v: got %v\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Following)
			}
		}
	}

	// Given the prior Server
	// And the request is POST /api/profiles/Anne/follow
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid Profile
	// And the username should be "Anne"
	// And the following flag should be true
	req = request("POST", "/api/profiles/Anne/follow", nil, validBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var profileResponse conduit.ProfileResponse
		if err := fetch(w.Result().Body, &profileResponse); err != nil {
			t.Errorf("profile: %s %s response did not contain valid ProfileResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := "Anne"; profileResponse.Profile.Username != expected {
				t.Errorf("profile: %s %s username expected %q: got %q\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Username)
			}
			if expected := true; profileResponse.Profile.Following != expected {
				t.Errorf("profile: %s %s following expected %v: got %v\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Following)
			}
		}
	}

	// Given the prior Server
	// And the request is GET /api/profiles/Anne
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid Profile
	// And the username should be "Anne"
	// And the following flag should be true
	req = request("GET", "/api/profiles/Anne", nil, validBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var profileResponse conduit.ProfileResponse
		if err := fetch(w.Result().Body, &profileResponse); err != nil {
			t.Errorf("profile: %s %s response did not contain valid ProfileResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := "Anne"; profileResponse.Profile.Username != expected {
				t.Errorf("profile: %s %s username expected %q: got %q\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Username)
			}
			if expected := true; profileResponse.Profile.Following != expected {
				t.Errorf("profile: %s %s following expected %v: got %v\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Following)
			}
		}
	}

	// Given the prior Server
	// And the request is DELETE /api/profiles/Anne/follow
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid Profile
	// And the username should be "Anne"
	// And the following flag should be false
	req = request("DELETE", "/api/profiles/Anne/follow", nil, validBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var profileResponse conduit.ProfileResponse
		if err := fetch(w.Result().Body, &profileResponse); err != nil {
			t.Errorf("profile: %s %s response did not contain valid ProfileResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := "Anne"; profileResponse.Profile.Username != expected {
				t.Errorf("profile: %s %s username expected %q: got %q\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Username)
			}
			if expected := false; profileResponse.Profile.Following != expected {
				t.Errorf("profile: %s %s following expected %v: got %v\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Following)
			}
		}
	}

	// Given the prior Server
	// And the request is GET /api/profiles/Anne
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// When we execute the request
	// Then the response should have a status of 200 (ok)
	// And contain a valid Profile
	// And the username should be "Anne"
	// And the following flag should be false
	req = request("GET", "/api/profiles/Anne", nil, validBearerToken)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if expected := http.StatusOK; w.Code != expected {
		t.Errorf("user: %s %s expected %d(%s): got %d(%s)\n", req.Method, req.URL.Path, expected, http.StatusText(expected), w.Code, http.StatusText(w.Code))
	} else {
		var profileResponse conduit.ProfileResponse
		if err := fetch(w.Result().Body, &profileResponse); err != nil {
			t.Errorf("profile: %s %s response did not contain valid ProfileResponse: %+v\n", req.Method, req.URL.Path, err)
		} else {
			if expected := "Anne"; profileResponse.Profile.Username != expected {
				t.Errorf("profile: %s %s username expected %q: got %q\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Username)
			}
			if expected := false; profileResponse.Profile.Following != expected {
				t.Errorf("profile: %s %s following expected %v: got %v\n", req.Method, req.URL.Path, expected, profileResponse.Profile.Following)
			}
		}
	}
}
