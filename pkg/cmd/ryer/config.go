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

package main

import (
	"flag"
	"os"
	"time"

	"github.com/peterbourgon/ff/v3"
)

type config struct {
	Debug    bool
	FileName string
	Server   struct {
		Scheme  string
		Host    string
		Port    string
		Timeout struct {
			Idle  time.Duration
			Read  time.Duration
			Write time.Duration
		}
		PublicRoot string
		Salt       string
		Key        string
	}
	Cookies struct {
		HttpOnly bool
		Secure   bool
	}
	Data struct {
		Path            string
		TimestampFormat string
	}
}

// getConfig returns a configuration.
// It will load a configuration file if one is given on the command line via the `-config` flag.
// If provided, the file must contain a valid JSON object.
//
// Priority is:
//   command line         overrides default value or configuration file or environment variable
//   environment variable overrides default value or configuration file
//   configuration file   overrides default value
func getConfig() (*config, error) {
	var cfg config
	cfg.Data.Path = "D:/GoLand/conduit/testdata/"
	cfg.Data.TimestampFormat = "2006-01-02T15:04:05.99999999Z"
	cfg.Server.Scheme = "http"
	cfg.Server.Host = "localhost"
	cfg.Server.Port = "3000"
	cfg.Server.Timeout.Idle = 10 * time.Second
	cfg.Server.Timeout.Read = 5 * time.Second
	cfg.Server.Timeout.Write = 10 * time.Second
	cfg.Server.PublicRoot = "D:/GoLand/conduit/public"
	cfg.Server.Salt = "pepper"
	cfg.Server.Key = "curry"

	var (
		fs                    = flag.NewFlagSet("server", flag.ExitOnError)
		fileName              = fs.String("config", cfg.FileName, "config file (optional)")
		debug                 = fs.Bool("debug", cfg.Debug, "log debug information (optional)")
		dataPath              = fs.String("data-path", cfg.Data.Path, "path containing data files")
		serverCookiesHttpOnly = fs.Bool("cookies-http-only", cfg.Cookies.HttpOnly, "set HttpOnly flag on cookies")
		serverCookiesSecure   = fs.Bool("cookies-secure", cfg.Cookies.Secure, "set Secure flag on cookies")
		serverScheme          = fs.String("scheme", cfg.Server.Scheme, "http scheme, either 'http' or 'https'")
		serverHost            = fs.String("host", cfg.Server.Host, "host name (or IP) to listen on")
		serverPort            = fs.String("port", cfg.Server.Port, "port to listen on")
		serverPublicRoot      = fs.String("public-root", cfg.Server.PublicRoot, "path to serve static files from")
		serverSalt            = fs.String("salt", cfg.Server.Salt, "set salt for hashing")
		serverTimeoutIdle     = fs.Duration("idle-timeout", cfg.Server.Timeout.Idle, "http idle timeout")
		serverTimeoutRead     = fs.Duration("read-timeout", cfg.Server.Timeout.Read, "http read timeout")
		serverTimeoutWrite    = fs.Duration("write-timeout", cfg.Server.Timeout.Write, "http write timeout")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("CONDUIT_RYER_SERVER"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return nil, err
	}

	cfg.Debug = *debug
	cfg.FileName = *fileName
	cfg.Cookies.HttpOnly = *serverCookiesHttpOnly
	cfg.Cookies.Secure = *serverCookiesSecure
	cfg.Data.Path = *dataPath
	cfg.Server.Scheme = *serverScheme
	cfg.Server.Host = *serverHost
	cfg.Server.Port = *serverPort
	cfg.Server.PublicRoot = *serverPublicRoot
	cfg.Server.Salt = *serverSalt
	cfg.Server.Timeout.Idle = *serverTimeoutIdle
	cfg.Server.Timeout.Read = *serverTimeoutRead
	cfg.Server.Timeout.Write = *serverTimeoutWrite

	return &cfg, nil
}
