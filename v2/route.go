package lightmux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// Middleware defines a function to process middleware logic for HTTP handlers.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// RouteGroup represents a group of routes with a common prefix and shared middlewares.
type RouteGroup struct {
	prefix      string
	middlewares []Middleware
	mux         *LightMux
}

// Route represents an HTTP route with its path, supported methods, and middlewares.
type Route struct {
	Path        string
	Methods     map[string]http.HandlerFunc
	Middlewares []Middleware
}

func isValidMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect, http.MethodTrace:
		return true
	default:
		return false
	}
}

// NewGroup creates a new RouteGroup with the given prefix and optional middlewares.
func (l *LightMux) NewGroup(prefix string, middlewares ...Middleware) *RouteGroup {
	return &RouteGroup{
		prefix:      prefix,
		middlewares: middlewares,
		mux:         l,
	}
}

// NewRoute creates a new Route within the RouteGroup with the given path and optional middlewares.
func (g *RouteGroup) NewRoute(path string, middlewares ...Middleware) *Route {
	fullPath := g.prefix + path
	allMiddleware := append(g.middlewares, middlewares...)
	return g.mux.NewRoute(fullPath, allMiddleware...)
}

// NewRoute creates a new Route with the given path and optional middlewares.
func (l *LightMux) NewRoute(path string, middlewares ...Middleware) *Route {
	// Check for duplicate path
	if _, exists := l.routeMap[path]; exists {
		panic(fmt.Sprintf("route with path %v already exists", path))
	}

	r := &Route{
		Path:        path,
		Methods:     make(map[string]http.HandlerFunc),
		Middlewares: middlewares,
	}

	l.routeStack = append(l.routeStack, r)
	l.routeMap[path] = struct{}{}

	return r
}

// Use adds middlewares into route middlewares.
func (r *Route) Use(middlewares ...Middleware)  {
	r.Middlewares = append(r.Middlewares, middlewares...)
}

// wrapMiddlewares applies the route's middlewares to the given handler.
func (r *Route) wrapMiddlewares(handler http.HandlerFunc) http.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		handler = r.Middlewares[i](handler)
	}
	return handler
}

// Handle registers a handler for a specific HTTP method on the route.
// Middlewares are not wrapped here; they are applied when serving the request.
func (r *Route) Handle(method string, handler http.HandlerFunc) {
	if !isValidMethod(method) {
		panic("invalid HTTP method: " + method)
	}

	// check if method already exists
	if _, exists := r.Methods[method]; exists {
		panic("duplicate method for path: " + method + " " + r.Path)
	}

	r.Methods[method] = handler
}

// ApplyRoutes registers all routes that have been created with NewRoute.
//
// Run() calls this before starting HTTP server, and before applying any global middlewares.
// This ensures all route handlers are registered to the underlying mux.
func (l *LightMux) ApplyRoutes() {
	for _, route := range l.routeStack {
		route := route
		l.mux.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			if handler, ok := route.Methods[r.Method]; ok {
				// Wrap middlewares here before serving
				route.wrapMiddlewares(handler).ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{
					"error": fmt.Sprintf("No handler for method %s for path %s", r.Method, r.URL.Path),
				})
				return
			}
		})
	}
}

// PrintRoutes prints all registered routes and their supported methods.
func (l *LightMux) PrintRoutes() {
	for _, r := range l.routeStack {
		fmt.Printf("Route: %s\n", r.Path)
		for method, handler := range r.Methods {
			fmt.Printf("\t- %s (handler: %s)\n", method, getFuncName(handler))
		}
		fmt.Printf("\tMiddlewares: %d\n", len(r.Middlewares))
		for i, mw := range r.Middlewares {
			fmt.Printf("\t\t%d: %T (%s)\n", i+1, mw, getFuncName(mw))
		}
	}
}

// getFuncName returns the name of the function for the given handler or middleware.
func getFuncName(h interface{}) string {
	return runtime.FuncForPC(
		reflect.ValueOf(h).Pointer(),
	).Name()
}
