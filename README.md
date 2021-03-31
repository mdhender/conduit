# Conduit
Conduit implements several styles of Go web servers for [Conduit](https://github.com/gothinkster/realworld), the RealWorld example web application.

# Servers
The grand plan is to implement several servers.

Look in `cmd/xxx` to find the `main` for each server.
The `main` function creates a `Config`, then calls `run` to initialize and execute the server.
All errors are returned to `main`, which prints the error and then exits.

Run uses the package imported from `internal/servers/xxx` to initialize the server.
It should have `server` declared in the main package, but because of the way I wrote
the test suite, I have to have `Server` and that `Server` has to have exported fields.
(I'd like to revisit and fix that in the future.)

# Configuration
All servers use the `internal/config` package.
Normally, that would be declared in the `main` package.
It's separated out soley to allow it to be reused.

# Test Suite
The servers share a common test suite.


