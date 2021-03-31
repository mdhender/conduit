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

import "testing"

func Profile(newServer TestServer, t *testing.T) {
	// Specification: Profile API

	// When given a new Server
	// And the user with username "Jacob," e-mail "jake@jake.jake," and password "jakejake" has been added
	// And the user with username "Anne," e-mail "anne@anne.anne," and password "anneanne" has been added
	// And the request is GET /api/profiles/Anne
	// And the request content type header is "application/json; charset=utf-8"
	// Then executing the request should succeed with status of 200 (ok)
	// And return a valid Profile with the username of "Anne"
	t.Errorf("!implemented")

	// When given the prior Server
	// And the request is POST /api/profiles/:username/follow
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// Then executing the request should succeed with status of 200 (ok)
	// And return a valid Profile with the username of "Anne"
	t.Errorf("!implemented")

	// When given the prior Server
	// And the request is DELETE /api/profiles/Anne/follow
	// And the request includes a valid bearer token for the user "jake@jake.jake"
	// Then executing the request should succeed with status of 200 (ok)
	// And return a valid Profile with the username of "Anne"
	t.Errorf("!implemented")
}
