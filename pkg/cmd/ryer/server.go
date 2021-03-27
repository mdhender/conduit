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
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/jwt"
	"log"
	"net/http"
	"time"

	"github.com/mdhender/conduit/internal/way"
)

type server struct {
	http.Server
	router       *way.Router
	tokenFactory jwt.Factory
}

func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAdmin {
			log.Printf("%s: not admin\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

func (s *server) authenticatedOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAuthenticated {
			log.Printf("%s: not authenticated\n", r.URL.Path)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

// post body should contain a NewUserRequest which wraps a NewUser
// Returns a UserResponse which wraps a User
func (s *server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("createUser\n")
		var req conduit.NewUserRequest
		req.User.Username = "Jacob"
		req.User.Email = "jake@jake.jake"
		req.User.Password = "jakejake"

		var err error
		var result conduit.UserResponse
		result.User.Username = req.User.Username
		result.User.Email = req.User.Email
		result.User.Token, err = s.tokenFactory.NewToken("", result.User.Username, result.User.Email, []string{"authenticated"}, 24*time.Hour)
		if err != nil {
			log.Printf("createUser: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(result)
		if err != nil {
			log.Printf("createUser: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
	}
}

func (s *server) getArticlesFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getArticlesFeed(%s)\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *server) handleAdminIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("adminIndex(%s)\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (s *server) handleGetArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getArticles(%s)\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: not found\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *server) handleNotImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: not implemented\n", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func currentUser(r *http.Request) (user struct {
	IsAdmin         bool
	IsAuthenticated bool
}) {
	j, err := jwt.GetBearerToken(r)
	if err != nil || !j.IsValid() {
		return user
	}
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
