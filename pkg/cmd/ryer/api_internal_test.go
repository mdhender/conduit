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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"github.com/mdhender/conduit/internal/way"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newServer() *server {
	srv := &server{
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

var contentType string = "application/json; charset=utf-8"

func fetch(body io.Reader, data interface{}) error {
	dec := json.NewDecoder(body)

	// enforce checking for unknown fields when parsing the body.
	dec.DisallowUnknownFields()

	err := dec.Decode(&data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fmt.Errorf("body contains unknown field %s", strings.TrimPrefix(err.Error(), "json: unknown field "))
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body must only contain a single JSON object")
	}

	return nil
}

func TestUser(t *testing.T) {
	// Specification: User API

	// When given a new server
	srv := newServer()
	// And the request is GET /api/user
	req := httptest.NewRequest("GET", "/api/user", nil)
	// Then executing the request should fail with status of 401 (not authorized)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnauthorized, w.Code)
	}
}

func TestUsers(t *testing.T) {
	// Specification: Users API

	// When given a new server
	srv := newServer()
	// And the request is POST /api/users with a valid user with no Content-Type header
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(`{"user":{"username": "Jacob","email": "jake@jake.jake","password": "jakejake"}}`))
	// Then executing the request should succeed with status of 200 (OK)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("api: %q %q expected %d: got %d\n", req.Method, req.URL.Path, http.StatusUnsupportedMediaType, w.Code)
	}

	// When given a new server
	srv = newServer()
	// And the request is POST /api/users with a valid user
	req = httptest.NewRequest("POST", "/api/users", strings.NewReader(`{"user":{"username": "Jacob","email": "jake@jake.jake","password": "jakejake"}}`))
	req.Header.Set("Content-Type", contentType)
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
			if userResponse.User.Email != "jake@jake.jake" {
				t.Errorf("api: %q %q username expected %q: got %q\n", req.Method, req.URL.Path, "Jacob", userResponse.User.Username)
			}
			if userResponse.User.Username != "Jacob" {
				t.Errorf("api: %q %q username expected %q: got %q\n", req.Method, req.URL.Path, "Jacob", userResponse.User.Username)
			}
		}
	}
}
