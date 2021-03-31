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

package ryer

import (
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"net/http"
)

// currentUser extracts data for the user making the request.
// It always returns a user struct, even if the request does
// not have a valid bearer token.
// TODO: should return a Conduit User.
func (s *Server) currentUser(r *http.Request) (user struct {
	IsAdmin         bool
	IsAuthenticated bool
	User            memory.User
}) {
	j, err := jwt.GetBearerToken(r)
	if err != nil {
		//log.Printf("currentUser: bearerToken %v\n", j)
		//log.Printf("currentUser: getBearerToken %+v\n", err)
		return user
	} else if err = s.TokenFactory.Validate(j); err != nil {
		//log.Printf("currentUser: validateToken %+v\n", err)
		return user
	} else if !j.IsValid() {
		return user
	}
	user.User, _ = s.DB.GetUser(j.Data().Id)
	for _, role := range j.Data().Roles {
		switch role {
		case "admin":
			user.IsAdmin = true
		case "authenticated":
			user.IsAuthenticated = true
		}
	}
	return user
}
