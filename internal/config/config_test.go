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

// Package config_test implements tests against the exposed Config API
//
// Apparently https://github.com/golang/go/issues/31859#issuecomment-489889428
// kind of broke tests that need command line flags. Limiting test to just the
// default values.
package config_test

import (
	"github.com/mdhender/conduit/internal/config"
	"testing"
)

func TestDefault(t *testing.T) {
	// Specification: Config API

	// Given an empty environment
	// When we create a new configuration
	cfg := config.Default()
	// Then we should have a valid Config
	// And it should have the expected default values
	if cfg == nil {
		t.Fatalf("default: expected config to be non-nil: got nil\n")
	}
	// And default values should be ...
	if expected := false; cfg.Debug != expected {
		t.Errorf("default: expected debug to be %v: got %v\n", expected, cfg.Debug)
	}
	if expected := "http"; cfg.Server.Scheme != expected {
		t.Errorf("default: expected server.scheme to be %q: got %q\n", expected, cfg.Server.Scheme)
	}
	if expected := "localhost"; cfg.Server.Host != expected {
		t.Errorf("default: expected server.host to be %q: got %q\n", expected, cfg.Server.Host)
	}
	if expected := "3000"; cfg.Server.Port != expected {
		t.Errorf("default: expected server.port to be %q: got %q\n", expected, cfg.Server.Port)
	}
}
