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

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

type Server interface {
	NewJWT(duration time.Duration, id int, username, email string, roles ...string) string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type keyValue struct {
	key, value string
}

var contentType = keyValue{"Content-Type", "application/json; charset=utf-8"}
var secret = "salt+pepper"

func request(method, target string, v interface{}, keys ...keyValue) *http.Request {
	var body io.Reader
	if v != nil {
		buf, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Sprintf("assert(err != %+v)", err))
		}
		body = bytes.NewReader(buf)
	}
	req := httptest.NewRequest(method, target, body)
	for _, kv := range keys {
		req.Header.Set(kv.key, kv.value)
	}
	return req
}
