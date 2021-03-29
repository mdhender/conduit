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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Specification: User API
func User(t *testing.T, newServer func(string) Server) {
	srv := newServer(secret)
	validBearerToken := keyValue{
		key:   "Authorization",
		value: "Bearer " + srv.NewJWT(15*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated"),
	}
	expiredBearerToken := keyValue{
		key:   "Authorization",
		value: "Bearer " + srv.NewJWT(0*time.Second, 1, "Jacob", "jake@jake.jake", "authenticated"),
	}

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
