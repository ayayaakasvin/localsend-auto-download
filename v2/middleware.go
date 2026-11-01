package lightmux

import (
	"fmt"
	"net/http"
)


// Global middleware functions are applied to all incoming HTTP requests handled by the server.
// Func registers Middleware and can be used for logging, authentication, etc.
// Changes will be applied to server after runnung LightMux.Run func.
func (l *LightMux) Use(middlewares ...Middleware) {
	if len(middlewares) != 0 {
		l.globalMiddlewareStack = append(l.globalMiddlewareStack, middlewares...)
	}
}

func chainMiddlewares(handler http.HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ApplyGlobalMiddlewares applies all registered global middlewares to the HTTP handler.
// This method is called after all routes have been registered and
// before starting the HTTP server (inside Run() method).
func (l *LightMux) ApplyGlobalMiddlewares() {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.mux.ServeHTTP(w, r)
	})

	finalHandler := base
	if len(l.globalMiddlewareStack) > 0 {
		finalHandler = chainMiddlewares(base, l.globalMiddlewareStack)
	}
	l.server.Handler = finalHandler
}

// Prints count of registered middlewares
func (l *LightMux) PrintMiddlewareInfo() {
	fmt.Printf("Global middleware count: %d\n", len(l.globalMiddlewareStack))
}
