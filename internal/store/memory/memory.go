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
	return &Store{
		users: make(map[string]User),
	}
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
	for _, u := range s.users {
		if u.Username == username {
			errs["username"] = append(errs["username"], "has already been taken")
		}
		if u.Email == email {
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
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()
	s.users[u.Username] = u

	return u, nil
}

func (s *Store) GetUser(id int) (User, bool) {
	s.Lock()
	defer s.Unlock()
	for _, u := range s.users {
		if u.Id == id {
			return u, true
		}
	}
	return User{Id: id}, false
}

type Store struct {
	sync.RWMutex
	seq   int
	users map[string]User
}

type User struct {
	Id        int
	Username  string
	Email     string
	Password  string
	Bio       string
	Image     string
	CreatedAt time.Time // "2021-03-27T16:58:01.233Z"
	UpdatedAt time.Time // "2021-03-27T16:58:01.245Z"
}
