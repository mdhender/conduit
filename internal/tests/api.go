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

// Package tests implements a test suite for the Conduit server API.
//
// To use this package, your server must provide a NewServer function
// which returns a struct that satisfies the Server interface.
//
// You must also arrange for your server's test file to call the suite
// functions directly. Those are the functions that have the NewServer
// as a parameter.
package tests

import (
	"net/http"
	"testing"
	"time"
)

type Server interface {
	NewJWT(ttl time.Duration, id int, username, email string, roles ...string) string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type TestServer func(secret string) Server

type keyValue struct {
	key, value string
}

var contentType = keyValue{"Content-Type", "application/json; charset=utf-8"}
var secret = "salt+pepper"

func Suite(newServer TestServer, t *testing.T) {
	Registration(newServer, t)
	Authentication(newServer, t)
	User(newServer, t)
	Users(newServer, t)
	Profile(newServer, t)
}
