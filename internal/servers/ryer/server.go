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

// Package server implements a Conduit server in the style of Mat Ryer's server.
// (see https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html)
// (see https://svlapin.github.io/engineering/2019/09/14/go-patterns.html)
// Note: I would normally place this in the `main` package. I'm breaking it out
// to make it easier to run the combined test suite and to keep the `main` packages
// for all the servers as common as possible.
package server

import (
	"encoding/json"
	"errors"
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/jsonapi"
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"github.com/mdhender/conduit/internal/way"
	"log"
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

// NewServer returns a Server with default values
func NewServer(secret string) *Server {
	srv := &Server{
		DB:           memory.New(),
		DtFmt:        "2006-01-02T15:04:05.99999999Z",
		Router:       way.NewRouter(),
		TokenFactory: jwt.NewFactory(secret),
	}
	srv.MaxHeaderBytes = 1 << 20
	srv.Handler = srv.Router
	srv.Routes()
	return srv
}

func (s *Server) NewJWT(ttl time.Duration, id int, username, email string, roles ...string) string {
	return s.TokenFactory.NewToken(ttl, id, username, email, roles...)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
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

func (s *Server) authenticatedOnly(h http.HandlerFunc) http.HandlerFunc {
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

func (s *Server) handleCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := s.currentUser(r).User
		user := conduit.User{
			Email:    u.Email,
			Token:    s.TokenFactory.NewToken(24*time.Hour, u.Id, u.Username, u.Email, "authenticated"),
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		}
		data, err := json.Marshal(conduit.UserResponse{User: user})
		if err != nil {
			log.Printf("currentUser: %+v\n", err)
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
func (s *Server) handleCreateUser() http.HandlerFunc {
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

		u, errs := s.DB.CreateUser(req.User.Username, req.User.Email, req.User.Password)
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
		user := conduit.User{
			Id:        u.Id,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Username:  u.Username,
			Token:     s.TokenFactory.NewToken(24*time.Hour, u.Id, u.Username, u.Email, "authenticated"),
		}
		data, err := json.Marshal(conduit.UserResponse{User: user})
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

func (s *Server) getArticlesFeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("getArticlesFeed(%s)\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) handleAdminIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("adminIndex(%s)\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (s *Server) handleGetArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("getArticles(%s)\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req conduit.LoginUserRequest
		err := jsonapi.Data(w, r, s.rejectUnknownFields, &req)
		if err != nil {
			log.Printf("login: %+v\n", err)
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
		u, found := s.DB.Login(req.User.Email, req.User.Password)
		if !found {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		user := conduit.User{
			Email:    u.Email,
			Token:    s.TokenFactory.NewToken(24*time.Hour, u.Id, u.Username, u.Email, "authenticated"),
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		}
		data, err := json.Marshal(conduit.UserResponse{User: user})
		if err != nil {
			if s.debug {
				log.Printf("login: %+v\n", err)
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func (s *Server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("%s: not found\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) handleNotImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.debug {
			log.Printf("%s: not implemented\n", r.URL.Path)
		}
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func (s *Server) handleUpdateCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("updateCurrentUser\n")
		cu := s.currentUser(r).User

		var req conduit.UpdateUserRequest
		err := jsonapi.Data(w, r, s.rejectUnknownFields, &req)
		if err != nil {
			log.Printf("updateCurrentUser: %+v\n", err)
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
		log.Println(req)
		u, errs := s.DB.UpdateUser(cu.Id, req.User.Email, req.User.Bio, req.User.Image)
		if errs != nil {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			var result struct {
				Errors map[string][]string `json:"errors"`
			}
			result.Errors = errs
			data, err := json.Marshal(result)
			if err != nil {
				log.Printf("updateCurrentUser: %+v\n", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(data)
			return
		}
		user := conduit.User{
			Id:        u.Id,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Username:  u.Username,
			Token:     s.TokenFactory.NewToken(24*time.Hour, u.Id, u.Username, u.Email, "authenticated"),
		}
		data, err := json.Marshal(conduit.UserResponse{User: user})
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

func (s *Server) currentUser(r *http.Request) (user struct {
	IsAdmin         bool
	IsAuthenticated bool
	User            memory.User
}) {
	j, err := jwt.GetBearerToken(r)
	if err != nil {
		return user
	} else if err = s.TokenFactory.Validate(j); err != nil {
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
