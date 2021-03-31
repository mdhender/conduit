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

// Package postgres implements a make-believe connection to an external database
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mdhender/conduit/internal/store/model"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	pg *sql.DB
}

func New(user, password, host string) (*Store, error) {
	db := &Store{}
	var err error
	db.pg, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=verify-full", user, password, host))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Store) GetUser(id int) (*model.User, error) {
	row := db.pg.QueryRow(`SELECT USERNAME, EMAIL FROM USERS WHERE ID = $1`, id)
	var username, email string
	err := row.Scan(&username, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &model.User{
		Id:       id,
		Username: username,
		Email:    email,
	}, nil
}
