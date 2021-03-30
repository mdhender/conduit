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

// Package memory implements a non-persistent in-memory data store for Conduit data.
package memory

import (
	"strings"
	"sync"
	"time"
)

func New() *Store {
	return &Store{}
}

func (s *Store) CreateUser(username, email, password string) (User, map[string][]string) {
	s.Lock()
	defer s.Unlock()
	errs := make(map[string][]string)
	if username = strings.TrimSpace(username); username == "" {
		errs["username"] = append(errs["username"], "can't be blank")
	}
	if email = strings.TrimSpace(email); email == "" {
		errs["email"] = append(errs["email"], "can't be blank")
	}
	if password = strings.TrimSpace(password); password == "" {
		errs["password"] = append(errs["password"], "can't be blank")
	}
	for i := range s.users {
		if s.users[i].Username == username {
			errs["username"] = append(errs["username"], "has already been taken")
		}
		if s.users[i].Email == email {
			errs["email"] = append(errs["email"], "has already been taken")
		}
	}
	if len(errs) != 0 {
		return User{}, errs
	}

	s.seq++

	var u User
	u.Id = s.seq
	u.Username = username
	u.Email = email
	u.Password = password
	u.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.99999999Z")
	u.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05.99999999Z")
	s.users = append(s.users, u)

	return u, nil
}

func (s *Store) GetUser(id int) (User, bool) {
	s.Lock()
	defer s.Unlock()
	for i := range s.users {
		if s.users[i].Id == id {
			return s.users[i], true
		}
	}
	return User{}, false
}

func (s *Store) Login(email, password string) (User, bool) {
	s.Lock()
	defer s.Unlock()
	for i := range s.users {
		if s.users[i].Email == email { // yeah, we know, timing attack
			if s.users[i].Password != password {
				break
			}
			return s.users[i], true
		}
	}
	return User{}, false
}

func (s *Store) UpdateUser(id int, email, bio, image *string) (User, map[string][]string) {
	s.Lock()
	defer s.Unlock()
	errs := make(map[string][]string)
	for i := range s.users {
		if s.users[i].Id == id {
			u := s.users[i]
			changes := false
			if bio != nil {
				val := strings.TrimSpace(*bio)
				u.Bio = &val
				changes = true
			}
			if email != nil {
				val := strings.TrimSpace(*email)
				if *email != val {
					errs["email"] = append(errs["email"], "can't have leading or trailing spaces")
				} else if val == "" {
					errs["email"] = append(errs["email"], "must not be empty if provided")
				} else {
					u.Email = val
					changes = true
				}
			}
			if image != nil {
				val := strings.TrimSpace(*image)
				u.Image = &val
				changes = true
			}
			if len(errs) != 0 {
				return User{}, errs
			}
			if changes {
				s.users[i] = u
			}
			return u, nil
		}
	}
	errs["email"] = append(errs["email"], "no such email")
	return User{}, errs
}

type Store struct {
	sync.RWMutex
	seq   int
	users []User
}

type User struct {
	Id        int
	Username  string
	Email     string
	Password  string
	CreatedAt string // "2021-03-27T16:58:01.233Z"
	UpdatedAt string // "2021-03-27T16:58:01.245Z"
	Bio       *string
	Image     *string
}
