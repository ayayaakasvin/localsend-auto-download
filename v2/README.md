# lightmux

`lightmux` is a lightweight HTTP server multiplexer for Go, providing simple routing and middleware support. It wraps the standard library's `http.ServeMux` and adds convenient methods for route and middleware management.

## Features

- Lightweight and minimal dependencies
- Route registration with method support
- Global and per-route middleware
- Easy integration with `http.Server`
- Customizable via direct access to `http.ServeMux`

## Installation

```sh
go get github.com/ayayaakasvin/lightmux
```

## Usage

```go
import (
    "net/http"
    "github.com/ayayaakasvin/lightmux"
)

func main() {
    server := &http.Server{Addr: ":8080"}
    mux := lightmux.NewLightMux(server)

    // Global middleware example
    mux.Use(func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // Logging, authentication, etc.
            next(w, r)
        }
    })

    // Route with middleware
    route := mux.NewRoute("/hello")
    route.Handle("GET", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })

    // Start server
    if err := mux.Run(); err != nil {
        panic(err)
    }
}
```

## API Reference

### Types

#### `type Middleware func(http.HandlerFunc) http.HandlerFunc`

Defines a function type for HTTP middleware (both global and per-route).

#### `type LightMux struct`

Manages the HTTP server and routing. Holds a reference to an `http.Server` and an `http.ServeMux` for handler registration.

#### `type Route struct`

Represents an HTTP route with its path, supported methods, and middlewares.

```go
type Route struct {
    Path        string
    Methods     map[string]http.HandlerFunc
    Middlewares []Middleware
}
```

#### `type RouteGroup struct`

Represents a group of routes sharing a common path prefix and optional middlewares.

```go
type RouteGroup struct {
    Prefix      string
    Middlewares []Middleware
    // ...other fields...
}
```

### Functions and Methods

#### `func NewLightMux(server *http.Server) *LightMux`

Creates and returns a new `LightMux` instance using the provided `http.Server`.

#### `func (l *LightMux) ApplyGlobalMiddlewares()`

Applies all registered global middlewares to the HTTP handler. Called after all routes have been registered and before starting the HTTP server (inside `Run()`).

#### `func (l *LightMux) ApplyRoutes()`

Registers all routes that have been created with `NewRoute`. Called by `Run()` before starting the HTTP server and before applying any global middlewares.

#### `func (l *LightMux) Mux() *http.ServeMux`

Returns the internal `http.ServeMux` used by `LightMux` for handler registration. Allows direct access for advanced routing or customization (e.g., adding a custom 404 handler).

#### `func (l *LightMux) NewRoute(path string, middlewares ...Middleware) *Route`

Creates a new `Route` with the given path and optional per-route middlewares.

#### `func (l *LightMux) Group(prefix string, middlewares ...Middleware) *RouteGroup`

Creates a new `RouteGroup` with the given path prefix and optional middlewares.

#### `func (g *RouteGroup) NewRoute(path string, middlewares ...Middleware) *Route`

Creates a new `Route` within the group, combining the group's prefix with the given path and applying both group and route middlewares.

#### `func (g *RouteGroup) Use(middlewares ...Middleware)`

Adds middleware(s) to the group, to be applied to all routes within the group.

#### `func (l *LightMux) PrintMiddlewareInfo()`

Prints the count of registered global and per-route middlewares.

#### `func (l *LightMux) PrintRoutes()`

Prints all registered routes and their supported methods.

#### `func (l *LightMux) Run() error`

Applies routes and global middlewares, then starts the HTTP server. Returns any error encountered while running the server. Shuts down gracefully when the server is stopped.

#### `func (l *LightMux) Use(middlewares ...Middleware)`

Registers global middleware functions to be applied to all incoming HTTP requests handled by the server. Useful for logging, authentication, etc. Global middlewares are applied in the order they are registered, before any per-route middlewares.

#### `func (r *Route) Handle(method string, handler http.HandlerFunc)`

Registers a handler for a specific HTTP method on the route.

#### `func (r *Route) Use(middlewares ...Middleware)`

Adds middleware(s) to the route, to be applied only to this route.

---

For more details, see the [GoDoc](https://pkg.go.dev/github.com/ayayaakasvin/lightmux).

