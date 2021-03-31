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

package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/peterbourgon/ff/v3"
)

type Config struct {
	Debug bool
	App   struct {
		Root            string
		TimestampFormat string
	}
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
		TLS struct {
			Serve    bool
			CertFile string
			KeyFile  string
		}
		Salt    string
		Key     string
		WebRoot string
	}
	Cookies struct {
		HttpOnly bool
		Secure   bool
	}
	Data struct {
		Path string
	}
}

// Default returns a default configuration.
// These are the values without loading the environment, configuration file, or command line.
func Default() *Config {
	var cfg Config
	cfg.App.Root = "D:/GoLand/conduit/"
	cfg.App.TimestampFormat = "2006-01-02T15:04:05.99999999Z"
	cfg.Data.Path = cfg.App.Root + "test/data/"
	cfg.Server.Scheme = "http"
	cfg.Server.Host = "localhost"
	cfg.Server.Port = "3000"
	cfg.Server.Timeout.Idle = 10 * time.Second
	cfg.Server.Timeout.Read = 5 * time.Second
	cfg.Server.Timeout.Write = 10 * time.Second
	cfg.Server.Key = "curry.aka.yrruc"
	cfg.Server.Salt = "pepper"
	cfg.Server.WebRoot = cfg.App.Root + "web/"
	return &cfg
}

// Load updates the values in a Config in this order:
//   1. It will load a configuration file if one is given on the
//      command line via the `-config` flag. If provided, the file
//      must contain a valid JSON object.
//   2. Environment variables, using the prefix `CONDUIT_RYER_SERVER`
//   3. Command line flags
func (cfg *Config) Load() error {
	fs := flag.NewFlagSet("Server", flag.ExitOnError)
	fileName := fs.String("config", cfg.FileName, "config file (optional)")
	debug := fs.Bool("debug", cfg.Debug, "log debug information (optional)")
	appRoot := fs.String("root", cfg.App.Root, "path to treat as root for relative file references")
	dataPath := fs.String("data-path", cfg.Data.Path, "path containing data files")
	serverCookiesHttpOnly := fs.Bool("cookies-http-only", cfg.Cookies.HttpOnly, "set HttpOnly flag on cookies")
	serverCookiesSecure := fs.Bool("cookies-secure", cfg.Cookies.Secure, "set Secure flag on cookies")
	serverScheme := fs.String("scheme", cfg.Server.Scheme, "http scheme, either 'http' or 'https'")
	serverHost := fs.String("host", cfg.Server.Host, "host name (or IP) to listen on")
	serverPort := fs.String("port", cfg.Server.Port, "port to listen on")
	serverKey := fs.String("key", cfg.Server.Key, "set key for signing tokens")
	serverSalt := fs.String("salt", cfg.Server.Salt, "set salt for hashing")
	serverTimeoutIdle := fs.Duration("idle-timeout", cfg.Server.Timeout.Idle, "http idle timeout")
	serverTimeoutRead := fs.Duration("read-timeout", cfg.Server.Timeout.Read, "http read timeout")
	serverTimeoutWrite := fs.Duration("write-timeout", cfg.Server.Timeout.Write, "http write timeout")
	serverTLSServe := fs.Bool("https", cfg.Server.TLS.Serve, "serve https")
	serverTLSCertFile := fs.String("https-cert-file", cfg.Server.Host, "https certificate file")
	serverTLSKeyFile := fs.String("https-key-file", cfg.Server.Host, "https certificate key file")
	serverWebRoot := fs.String("web-root", cfg.Server.WebRoot, "path to serve web assets from")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("CONDUIT_RYER_SERVER"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return err
	}

	cfg.Debug = *debug
	cfg.App.Root = path.Clean(*appRoot)
	cfg.FileName = *fileName
	cfg.Cookies.HttpOnly = *serverCookiesHttpOnly
	cfg.Cookies.Secure = *serverCookiesSecure
	cfg.Data.Path = path.Clean(*dataPath)
	cfg.Server.Scheme = *serverScheme
	cfg.Server.Host = *serverHost
	cfg.Server.Port = *serverPort
	cfg.Server.Key = *serverKey
	cfg.Server.Salt = *serverSalt
	cfg.Server.Timeout.Idle = *serverTimeoutIdle
	cfg.Server.Timeout.Read = *serverTimeoutRead
	cfg.Server.Timeout.Write = *serverTimeoutWrite
	cfg.Server.TLS.Serve = *serverTLSServe
	cfg.Server.TLS.CertFile = *serverTLSCertFile
	cfg.Server.TLS.KeyFile = *serverTLSKeyFile
	cfg.Server.WebRoot = path.Clean(*serverWebRoot)

	if cfg.Server.TLS.Serve == true {
		if cfg.Server.TLS.CertFile == "" {
			return fmt.Errorf("must supply certificates file when serving HTTPS")
		}
		if cfg.Server.TLS.KeyFile == "" {
			return fmt.Errorf("must supply certificate key file when serving HTTPS")
		}
	}

	return nil
}
