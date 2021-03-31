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
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrNotAuthorized = errors.New("not authorized")
var ErrNotFound = errors.New("not found")

func New() *Store {
	db := &Store{}
	db.users.email = make(map[string]*User)
	db.users.id = make(map[int]*User)
	db.users.name = make(map[string]*User)
	return db
}

func (db *Store) CreateUser(username, email, password string) (*User, map[string][]string) {
	db.Lock()
	defer db.Unlock()
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
	if _, ok := db.users.name[username]; ok {
		errs["username"] = append(errs["username"], "has already been taken")
	}
	if _, ok := db.users.email[email]; ok {
		errs["email"] = append(errs["email"], "has already been taken")
	}
	if len(errs) != 0 {
		return nil, errs
	}

	db.seq++
	u := &User{
		Id:        db.seq,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.99999999Z"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.99999999Z"),
		Following: make(map[int]bool),
	}
	db.users.id[u.Id] = u
	db.users.name[u.Username] = u
	db.users.email[u.Email] = u

	return u.Copy(), nil
}

func (db *Store) FollowUserByUsername(id int, username string) (*Profile, error) {
	db.Lock()
	defer db.Unlock()

	user := db.users.id[id]
	if id == 0 || user == nil {
		return nil, ErrNotAuthorized
	} else if user.Username == username { // wants to follow self
		return nil, ErrNotAuthorized
	}

	target := db.users.name[username]
	if target == nil {
		return nil, ErrNotFound
	}
	user.Following[target.Id] = true

	return target.Profile(user), nil
}

func (db *Store) GetProfileByUsername(id int, username string) (*Profile, error) {
	db.Lock()
	defer db.Unlock()

	user := db.users.id[id]
	target := db.users.name[username]
	if target == nil {
		return nil, ErrNotFound
	}

	profile := target.Profile(user)
	if user != nil {
		profile.Following = user.Following[target.Id]
	}

	return profile, nil
}

func (db *Store) GetUser(id int) (*User, error) {
	db.Lock()
	defer db.Unlock()

	user := db.users.id[id]
	if user == nil {
		return nil, ErrNotFound
	}
	return user.Copy(), nil
}

func (db *Store) Login(email, password string) (*User, error) {
	db.Lock()
	defer db.Unlock()
	user := db.users.email[email]
	if user == nil || user.Password != password { // yeah, we know, timing attack
		return nil, ErrNotAuthorized
	}
	return user.Copy(), nil
}

func (db *Store) UpdateUser(id int, email, bio, image *string) (*User, map[string][]string) {
	db.Lock()
	defer db.Unlock()
	errs := make(map[string][]string)

	user := db.users.id[id]
	if user == nil {
		errs["email"] = append(errs["email"], "no such email")
		return nil, errs
	}

	cp, changes := user.Copy(), false
	if bio != nil {
		cp.bio = strings.TrimSpace(*bio)
		cp.Bio = &cp.bio
		changes = true
	}
	if email != nil {
		val := strings.TrimSpace(*email)
		if *email != val {
			errs["email"] = append(errs["email"], "can't have leading or trailing spaces")
		} else if val == "" {
			errs["email"] = append(errs["email"], "must not be empty if provided")
		} else {
			cp.Email = val
			changes = true
		}
	}
	if image != nil {
		cp.image = strings.TrimSpace(*image)
		cp.Image = &cp.image
		changes = true
	}
	if len(errs) != 0 {
		return nil, errs
	}

	if !changes {
		return user.Copy(), nil
	}

	// clear out pointers to the old record
	delete(db.users.id, user.Id)
	delete(db.users.name, user.Username)
	delete(db.users.email, user.Email)

	// and then insert the new record
	db.users.id[id] = cp
	db.users.name[user.Username] = cp
	db.users.email[user.Email] = cp

	return cp.Copy(), nil
}

type Store struct {
	sync.RWMutex
	seq   int
	users struct {
		id    map[int]*User
		name  map[string]*User
		email map[string]*User
	}
}

type Profile struct {
	Id         int
	Username   string
	Bio        *string
	Image      *string
	Following  bool
	bio, image string
}

type User struct {
	Id         int
	Username   string
	Email      string
	Password   string
	CreatedAt  string // "2021-03-27T16:58:01.233Z"
	UpdatedAt  string // "2021-03-27T16:58:01.245Z"
	Bio        *string
	Image      *string
	Following  map[int]bool // map of Id of users being followed
	bio, image string
}

func (u *User) Copy() *User {
	if u == nil {
		return &User{}
	}
	cp := &User{
		Id:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Following: make(map[int]bool),
		bio:       u.bio,
		image:     u.image,
	}
	if u.Bio != nil {
		cp.Bio = &cp.bio
	}
	if u.Image != nil {
		cp.Image = &cp.image
	}
	for id := range u.Following {
		cp.Following[id] = true
	}
	return cp
}

func (u *User) Profile(p *User) *Profile {
	if u == nil {
		return &Profile{}
	}
	profile := &Profile{
		Id:        u.Id,
		Username:  u.Username,
		Following: p != nil && p.Following[u.Id],
		bio:       u.bio,
		image:     u.image,
	}
	if u.Bio != nil {
		profile.Bio = &profile.bio
	}
	if u.Image != nil {
		profile.Image = &profile.image
	}
	return profile
}
