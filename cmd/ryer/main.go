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

// Package main implements a Conduit server in the style of Mat Ryer's server.
// (see https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html)
// (see https://svlapin.github.io/engineering/2019/09/14/go-patterns.html)
package main

import (
	"github.com/mdhender/conduit/internal/config"
	"github.com/mdhender/conduit/internal/jwt"
	"github.com/mdhender/conduit/internal/store/memory"
	"github.com/mdhender/conduit/internal/way"
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC) // force logs to be UTC
	log.Println("[main] entered")

	cfg := config.Default()
	if err := cfg.Load(); err != nil {
		log.Printf("[main] %+v\n", err)
		os.Exit(2)
	}

	if err := run(cfg); err != nil {
		log.Printf("[main] %+v\n", err)
		os.Exit(2)
	}
}

func run(cfg *config.Config) error {
	s := &Server{
		db:           memory.New(),
		dtfmt:        cfg.App.TimestampFormat,
		router:       way.NewRouter(),
		tokenFactory: jwt.NewFactory(cfg.Server.Salt + cfg.Server.Key),
	}
	s.Addr = net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	s.IdleTimeout = cfg.Server.Timeout.Idle
	s.ReadTimeout = cfg.Server.Timeout.Read
	s.WriteTimeout = cfg.Server.Timeout.Write
	s.MaxHeaderBytes = 1 << 20
	s.Handler = s.router

	s.routes()

	log.Printf("[server] listening on %s\n", s.Addr)
	return s.ListenAndServe()
}
