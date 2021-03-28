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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	type request struct {
		method string
		target string
		body   io.Reader
	}

	// Specification: API

	// When given a new server
	srv := defaultServer()
	srv.routes()
	// And the request is GET /api/user
	r := request{"GET", "/api/user", strings.NewReader("")}
	// Then executing the request should fail with status of 401 (not authorized)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, httptest.NewRequest(r.method, r.target, r.body))
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", r.method, r.target, http.StatusUnauthorized, w.Code)
	}

	// When given a new server
	srv = defaultServer()
	srv.routes()
	// And the request is POST /api/users with a valid user with no Content-Type header
	r = request{"POST", "/api/users", strings.NewReader(`{"user":{"username": "Jacob","email": "jake@jake.jake","password": "jakejake"}}`)}
	// Then executing the request should succeed with status of 200 (OK)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, httptest.NewRequest(r.method, r.target, r.body))
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("api: %q %q expected %d: got %d\n", r.method, r.target, http.StatusUnsupportedMediaType, w.Code)
	}

	// When given a new server
	srv = defaultServer()
	srv.routes()
	// And the request is POST /api/users with a valid user
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(`{"user":{"username": "Jacob","email": "jake@jake.jake","password": "jakejake"}}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// Then executing the request should succeed with status of 200 (OK)
	w = httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("api: %q %q expected %d: got %d\n", r.method, r.target, http.StatusOK, w.Code)
	}
	// And return a valid User
	t.Errorf("api: %q %q test not fully implemented\n", r.method, r.target)
}
