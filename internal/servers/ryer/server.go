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

// Package ryer implements a Conduit server in the style of Mat Ryer's server.
// (see https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html)
// (see https://svlapin.github.io/engineering/2019/09/14/go-patterns.html)
//
// Note: I would normally place this in the `main` package. I'm breaking it out
// to make it easier to run the combined test suite and to keep the `main` packages
// for all the servers as common as possible.
package ryer

import (
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"github.com/mdhender/conduit/internal/way"
	"net/http"
	"time"
)

type Server struct {
	http.Server
	DB                  *memory.Store
	DtFmt               string // format string for timestamps in responses
	Router              *way.Router
	TokenFactory        jwt.Factory
	debug               bool
	rejectUnknownFields bool
}

// NewJWT implements the tests.Server interface
func (s *Server) NewJWT(ttl time.Duration, id int, username, email string, roles ...string) string {
	return s.TokenFactory.NewToken(ttl, id, username, email, roles...)
}

// ServeHTTP implements the http handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
