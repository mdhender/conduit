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

import "net/http"

// routes initializes all routes exposed by the server
func (s *server) routes() {
	// routes are taken from https://github.com/gothinkster/realworld/blob/9686244365bf5681e27e2e9ea59a4d905d8080db/api/swagger.json
	for _, route := range []struct {
		pattern string
		method  string
		handler http.HandlerFunc
	}{
		{"/api/articles", "GET", s.handleNotImplemented()},
		{"/api/articles", "POST", s.handleNotImplemented()},
		{"/api/articles/feed", "GET", s.handleGetArticlesFeed()},
		{"/api/articles/:slug", "DELETE", s.handleNotImplemented()},
		{"/api/articles/:slug", "GET", s.handleGetArticles()},
		{"/api/articles/:slug", "PUT", s.handleNotImplemented()},
		{"/api/articles/:slug/comments", "GET", s.handleNotImplemented()},
		{"/api/articles/:slug/comments", "POST", s.handleNotImplemented()},
		{"/api/articles/:slug/comments/:id", "DELETE", s.handleNotImplemented()},
		{"/api/articles/:slug/favorite", "DELETE", s.handleNotImplemented()},
		{"/api/articles/:slug/favorite", "POST", s.handleNotImplemented()},
		{"/api/profiles/:username", "GET", s.handleNotImplemented()},
		{"/api/profiles/:username/follow", "DELETE", s.handleNotImplemented()},
		{"/api/profiles/:username/follow", "POST", s.handleNotImplemented()},
		{"/api/tags", "GET", s.handleNotImplemented()},
		{"/api/user", "GET", s.handleNotImplemented()},
		{"/api/user", "POST", s.handleNotImplemented()},
		{"/api/user", "PUT", s.handleNotImplemented()},
		{"/api/users/login", "POST", s.handleNotImplemented()},
	} {
		s.router.HandleFunc(route.method, route.pattern, route.handler)
	}
	s.router.NotFound = s.handleNotFound()
}
