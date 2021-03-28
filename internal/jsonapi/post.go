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

package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// todo: add source and license for this

type MalformedRequest struct {
	status int
	msg    string
}

func (mr *MalformedRequest) Status() int {
	return mr.status
}

func (mr *MalformedRequest) Error() string {
	return mr.msg
}

var ErrBadRequest = errors.New("bad request")
var ErrRequestEntityTooLarge = errors.New("request entity too larg")
var ErrUnsupportedMediaType = errors.New("unsupported media type")

func GetFormData(w http.ResponseWriter, r *http.Request, rejectUnknown bool, dst interface{}) error {
	switch r.Header.Get("Content-Type") {
	case "application/json": // ok
	case "application/json; charset=utf-8": // ok
	default:
		return fmt.Errorf("Content-Type %q is not supported: %w", r.Header.Get("Content-Type"), ErrUnsupportedMediaType)
	}

	// enforce a max body of 1mb. should make this a parameter.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)

	// maybe enforce checking for unknown fields when parsing the body.
	if rejectUnknown {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("Request body contains badly-formed JSON (at position %d): %w", syntaxError.Offset, ErrBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("Request body contains badly-formed JSON: %w", ErrBadRequest)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("Request body contains an invalid value for the %q field (at position %d): %w", unmarshalTypeError.Field, unmarshalTypeError.Offset, ErrBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fmt.Errorf("Request body contains unknown field %s: %w", strings.TrimPrefix(err.Error(), "json: unknown field "), ErrBadRequest)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("Request body must not be empty: %w", ErrBadRequest)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("Request body must not be larger than 1MB: %w", ErrRequestEntityTooLarge)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("Request body must only contain a single JSON object: %w", ErrBadRequest)
	}

	return nil
}
