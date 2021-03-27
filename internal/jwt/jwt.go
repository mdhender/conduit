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

// Package jwt implements naive JSON Web Tokens.
// Don't use this for anything other than testing.
package jwt

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type JWT struct {
	h struct {
		Algorithm   string `json:"alg,omitempty"` // message authentication code algorithm
		TokenType   string `json:"typ,omitempty"`
		ContentType string `json:"cty,omitempty"`
		KeyID       string `json:"kid,omitempty"` // optional identifier used to sign. doesn't work.
		b64         string // header marshalled to JSON and then base-64 encoded
	}
	p struct {
		// The principal that issued the JWT.
		Issuer string `json:"iss,omitempty"`
		// The subject of the JWT.
		Subject string `json:"sub,omitempty"`
		// The recipients that the JWT is intended for.
		// Each principal intended to process the JWT must identify itself with a value in the audience claim.
		// If the principal processing the claim does not identify itself with a value in the aud claim when this claim is present,
		// then the JWT must be rejected.
		Audience []string `json:"aud,omitempty"`
		// The expiration time on and after which the JWT must not be accepted for processing.
		// The value must be a NumericDate:[9] either an integer or decimal, representing seconds past 1970-01-01 00:00:00Z.
		ExpirationTime int64 `json:"exp,omitempty"`
		// The time on which the JWT will start to be accepted for processing.
		// The value must be a NumericDate.
		NotBefore int64 `json:"nbf,omitempty"`
		// The time at which the JWT was issued.
		// The value must be a NumericDate.
		IssuedAt int64 `json:"iat,omitempty"`
		// Case sensitive unique identifier of the token even among different issuers.
		JWTID string `json:"jti,omitempty"`
		// Private data for use by the application.
		Private struct {
			Algorithm string   `json:"alg"`
			TokenType string   `json:"typ"`
			Id        int      `json:"id,omitempty"`
			Username  string   `json:"username,omitempty"`
			Email     string   `json:"email,omitempty"`
			Roles     []string `json:"roles,omitempty"`
		} `json:"private"`
		b64 string // payload marshalled to JSON and then base-64 encoded
	}
	s        string // signature base-64 encoded
	isSigned bool   // true only if the signature has been verified
}

type Data struct {
	Id       int
	Username string
	Email    string
	Roles    []string
}

// pull the bearer token from a request header.
func GetBearerToken(r *http.Request) (*JWT, error) {
	headerAuthText := r.Header.Get("Authorization")
	if headerAuthText == "" {
		return nil, ErrUnauthorized
	}
	authTokens := strings.SplitN(headerAuthText, " ", 2)
	if len(authTokens) != 2 {
		return nil, ErrBadRequest
	}
	authType, authToken := authTokens[0], strings.TrimSpace(authTokens[1])
	if authType != "Bearer" {
		return nil, ErrUnauthorized
	}

	sections := strings.Split(authToken, ".")
	if len(sections) != 3 || len(sections[0]) == 0 || len(sections[1]) == 0 || len(sections[2]) == 0 {
		return nil, ErrBadRequest
	}

	var j JWT
	j.h.b64 = sections[0]
	j.p.b64 = sections[1]
	j.s = sections[2]

	// decode and extract the header from the token
	if rawHeader, err := decode(j.h.b64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(rawHeader, &j.h); err != nil {
		return nil, err
	} else if j.h.Algorithm == "" || j.h.Algorithm == "none" {
		return nil, ErrUnauthorized
	}

	// decode and extract the payload from the token
	if rawPayload, err := decode(j.p.b64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(rawPayload, &j.p); err != nil {
		return nil, err
	} else if j.h.TokenType != j.p.Private.TokenType {
		return nil, ErrUnauthorized
	} else if j.h.Algorithm != j.p.Private.Algorithm {
		return nil, ErrUnauthorized
	}

	return &j, nil
}

func (j *JWT) Data() Data {
	return Data{
		Id:       j.p.Private.Id,
		Username: j.p.Private.Username,
		Email:    j.p.Private.Email,
		Roles:    j.p.Private.Roles,
	}
}

func (j *JWT) IsValid() bool {
	if !j.isSigned || j.h.Algorithm != j.p.Private.Algorithm || j.h.TokenType != j.p.Private.TokenType {
		return false
	}
	now := time.Now()
	if j == nil || j.p.IssuedAt == 0 || j.p.ExpirationTime == 0 {
		return false
	} else if !now.After(time.Unix(j.p.IssuedAt, 0)) {
		return false
	} else if j.p.NotBefore != 0 && !now.Before(time.Unix(j.p.NotBefore, 0)) {
		return false
	}
	return time.Unix(j.p.ExpirationTime, 0).After(now)
}
