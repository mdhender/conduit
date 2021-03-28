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
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/jsonapi"
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"log"
	"net/http"
	"time"

	"github.com/mdhender/conduit/internal/way"
)

type server struct {
	http.Server
	debug               bool
	router              *way.Router
	tokenFactory        jwt.Factory
	db                  *memory.Store
	dtfmt               string // format string for timestamps in responses
	rejectUnknownFields bool
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.currentUser(r).IsAdmin {
			if s.debug {
				log.Printf("%s: not admin\n", r.URL.Path)
			}
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

func (s *server) authenticatedOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.currentUser(r).IsAuthenticated {
			if s.debug {
				log.Printf("%s: not authenticated\n", r.URL.Path)
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

func (s *server) handleCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("currentUser\n")
		}
		u := s.currentUser(r).User
		var result struct {
			User struct {
				Id        int     `json:"id"`
				Email     string  `json:"email"`
				CreatedAt string  `json:"createdAt"`
				UpdatedAt string  `json:"updatedAt"`
				Username  string  `json:"username"`
				Bio       *string `json:"bio"`   // API requires this to be nullable
				Image     *string `json:"image"` // API requires this to be nullable
				Token     string  `json:"token,omitempty"`
			} `json:"user"`
		}
		result.User.Id = u.Id
		result.User.Email = u.Email
		result.User.CreatedAt = u.CreatedAt.UTC().Format(s.dtfmt)
		result.User.UpdatedAt = u.UpdatedAt.UTC().Format(s.dtfmt)
		result.User.Username = u.Username
		if u.Bio != "" {
			result.User.Bio = &u.Bio
		}
		if u.Image != "" {
			result.User.Image = &u.Image
		}

		data, err := json.Marshal(result)
		if err != nil {
			if s.debug {
				log.Printf("currentUser: %+v\n", err)
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

// post body should contain a NewUserRequest which wraps a NewUser
// Returns a UserResponse which wraps a User
func (s *server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("createUser\n")
		}

		var req conduit.NewUserRequest
		err := jsonapi.Data(w, r, s.rejectUnknownFields, &req)
		if err != nil {
			if s.debug {
				log.Printf("createUser: %+v\n", err)
			}
			if errors.Is(err, jsonapi.ErrBadRequest) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			} else if errors.Is(err, jsonapi.ErrRequestEntityTooLarge) {
				http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			} else if errors.Is(err, jsonapi.ErrUnsupportedMediaType) {
				http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		u, errs := s.db.CreateUser(req.User.Username, req.User.Email, req.User.Password)
		if errs != nil {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			var result struct {
				Errors map[string][]string `json:"errors"`
			}
			result.Errors = errs
			data, err := json.Marshal(result)
			if err != nil {
				if s.debug {
					log.Printf("createUser: %+v\n", err)
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(data)
			return
		}

		var result struct {
			User struct {
				Id        int     `json:"id"`
				Email     string  `json:"email"`
				CreatedAt string  `json:"createdAt"`
				UpdatedAt string  `json:"updatedAt"`
				Username  string  `json:"username"`
				Bio       *string `json:"bio"`   // API requires this to be nullable
				Image     *string `json:"image"` // API requires this to be nullable
				Token     string  `json:"token"`
			} `json:"user"`
		}
		result.User.Id = u.Id
		result.User.Email = u.Email
		result.User.CreatedAt = u.CreatedAt.UTC().Format(s.dtfmt)
		result.User.UpdatedAt = u.UpdatedAt.UTC().Format(s.dtfmt)
		result.User.Username = u.Username

		result.User.Token, err = s.tokenFactory.NewToken(u.Id, result.User.Username, result.User.Email, []string{"authenticated"}, 24*time.Hour)
		if err != nil {
			if s.debug {
				log.Printf("createUser: %+v\n", err)
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(result)
		if err != nil {
			if s.debug {
				log.Printf("createUser: %+v\n", err)
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func (s *server) getArticlesFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("getArticlesFeed(%s)\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *server) handleAdminIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("adminIndex(%s)\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (s *server) handleGetArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("getArticles(%s)\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("%s: not found\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *server) handleNotImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("%s: not implemented\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (s *server) currentUser(r *http.Request) (user struct {
	IsAdmin         bool
	IsAuthenticated bool
	User            memory.User
}) {
	j, err := jwt.GetBearerToken(r)
	if err != nil {
		return user
	}
	if err = s.tokenFactory.Validate(j); err != nil || !j.IsValid() {
		return user
	}
	user.User, _ = s.db.GetUser(j.Data().Id)
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
