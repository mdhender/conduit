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
	"encoding/json"
	"errors"
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/jsonapi"
	"log"
	"net/http"
	"time"
)

func (s *Server) handleCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := conduit.User{}
		if u := s.currentUser(r).User; u != nil {
			user.Email = u.Email
			user.Token = s.TokenFactory.NewToken(24*time.Hour, u.Id, u.Username, u.Email, "authenticated")
			user.Username = u.Username
			user.Bio = u.Bio
			user.Image = u.Image
		}
		data, err := json.Marshal(conduit.UserResponse{User: user})
		if err != nil {
			log.Printf("currentUser: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func (s *Server) handleUpdateCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		u, errs := s.DB.UpdateUser(cu.Id, req.User.Email, req.User.Bio, req.User.Image)
		if errs != nil {
			w.Header().Add("Content-Type", contentType)
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
			Bio:       u.Bio,
			CreatedAt: u.CreatedAt,
			Email:     u.Email,
			Image:     u.Image,
			UpdatedAt: u.UpdatedAt,
			Username:  u.Username,
			Token:     s.TokenFactory.NewToken(24*time.Hour, u.Id, u.Username, u.Email, "authenticated"),
		}
		data, err := json.Marshal(conduit.UserResponse{User: user})
		if err != nil {
			log.Printf("updateCurrentUser: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
