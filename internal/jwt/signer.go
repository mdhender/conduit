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

package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
)

// Signer interface
type Signer interface {
	Algorithm() string
	Sign(msg []byte) ([]byte, error)
}

// HS256 implements a Signer using HMAC256.
type HS256 struct {
	secret []byte
}

func HS256Signer(secret []byte) *HS256 {
	h := HS256{secret: make([]byte, len(secret))}
	copy(h.secret, secret)
	return &h
}

// Algorithm implements the Signer interface
func (h *HS256) Algorithm() string {
	return "HS256"
}

// Sign implements the Signer interface
func (h *HS256) Sign(msg []byte) ([]byte, error) {
	hm := hmac.New(sha256.New, h.secret)
	if _, err := hm.Write(msg); err != nil {
		return nil, err
	}
	return hm.Sum(nil), nil
}
