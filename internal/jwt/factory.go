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
	"encoding/json"
	"time"
)

// NewFactory returns an initialized factory.
// The signer is used to sign the generated tokens.
func NewFactory(secret string) Factory {
	return Factory{s: HS256Signer([]byte(secret))}
}

// Consider setting up factory to use only one algorithm.
// If we do, can it still use a key id from the header?
// You should create a new factory when you rotate keys!
type Factory struct {
	s         Signer
	tokenType string
}

// Validate will return an error if the JWT is not properly signed.
func (f *Factory) Validate(j *JWT) error {
	expectedSignature, err := f.s.Sign([]byte(j.h.b64 + "." + j.p.b64))
	if err != nil {
		return err
	}

	if j.isSigned = j.s != encode(expectedSignature); !j.isSigned {
		return ErrUnauthorized
	}

	return nil // valid signature
}

func (f *Factory) NewToken(id int, username, email string, roles []string, ttl time.Duration) (string, error) {
	var j JWT

	j.h.TokenType = "JWT"
	j.h.Algorithm = f.s.Algorithm()
	j.p.IssuedAt = time.Now().Unix()
	j.p.ExpirationTime = time.Now().Add(ttl).Unix()
	j.p.Private.TokenType = j.h.TokenType
	j.p.Private.Algorithm = j.h.Algorithm
	j.p.Private.Id = id
	j.p.Private.Username = username
	j.p.Private.Email = email
	j.p.Private.Roles = roles

	if h, err := json.MarshalIndent(j.h, "  ", "  "); err != nil {
		return "", err
	} else {
		j.h.b64 = encode(h)
	}

	if p, err := json.MarshalIndent(j.p, "  ", "  "); err != nil {
		return "", err
	} else {
		j.p.b64 = encode(p)
	}

	rawSignature, err := f.s.Sign([]byte(j.h.b64 + "." + j.p.b64))
	if err != nil {
		return "", err
	}
	j.s = encode(rawSignature)

	return j.h.b64 + "." + j.p.b64 + "." + j.s, nil
}
