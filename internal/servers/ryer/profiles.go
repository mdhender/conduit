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
	"github.com/mdhender/conduit/internal/conduit"
	"github.com/mdhender/conduit/internal/way"
	"log"
	"net/http"
)

func (s *Server) handleFollowUserByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cu := s.currentUser(r).User

		username := way.Param(r.Context(), "username")
		profile, err := s.DB.FollowUserByUsername(cu.Id, username)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		data, err := json.Marshal(conduit.ProfileResponse{Profile: conduit.Profile{
			Bio:       profile.Bio,
			Following: profile.Following,
			Image:     profile.Image,
			Username:  profile.Username,
		}})
		if err != nil {
			log.Printf("followUserByUsername: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

func (s *Server) handleGetProfileByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// client doesn't have to be authenticated, but if she is,
		// we will fetch the following flag for her.
		var userId int
		if cu := s.currentUser(r).User; cu != nil {
			userId = cu.Id
		}

		username := way.Param(r.Context(), "username")
		profile, err := s.DB.GetProfileByUsername(userId, username)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		data, err := json.Marshal(conduit.ProfileResponse{Profile: conduit.Profile{
			Bio:       profile.Bio,
			Following: profile.Following,
			Image:     profile.Image,
			Username:  profile.Username,
		}})
		if err != nil {
			log.Printf("getProfileByUsername: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
