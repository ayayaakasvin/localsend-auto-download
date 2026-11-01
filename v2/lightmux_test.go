package lightmux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareExecution(t *testing.T) {

	var called []string

	logMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			called = append(called, "log")
			next(w, r)
		}
	}

	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			called = append(called, "auth")
			next(w, r)
		}
	}

	lmux := NewLightMux(&http.Server{})
	route := lmux.NewRoute("/mw", logMiddleware, authMiddleware)
	route.Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		called = append(called, "handler")
		w.WriteHeader(http.StatusOK)
	})

	lmux.ApplyRoutes()
	lmux.ApplyGlobalMiddlewares()

	req := httptest.NewRequest(http.MethodGet, "/mw", nil)
	w := httptest.NewRecorder()
	lmux.Mux().ServeHTTP(w, req)

	t.Log(called)

	mustResult := []string{"log", "auth", "handler"}
	for i := range mustResult {
		if mustResult[i] != called[i] {
			t.Fatalf("mw call order failed: %s != %s", mustResult[i], called[i])
		}
	}
}

func TestHandlerResponse(t *testing.T) {

	var called []string

	foo := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		called = append(called, "foo")
	}

	bar := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		called = append(called, "bar")
	}

	lmux := NewLightMux(&http.Server{})
	route := lmux.NewRoute("/call")
	route.Handle(http.MethodGet, foo)
	route.Handle(http.MethodPost, bar)

	lmux.ApplyRoutes()
	lmux.ApplyGlobalMiddlewares()

	req := httptest.NewRequest(http.MethodGet, "/call", nil)
	w := httptest.NewRecorder()
	lmux.Mux().ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodPost, "/call", nil)
	w = httptest.NewRecorder()
	lmux.Mux().ServeHTTP(w, req)

	t.Log(called)

	mustResult := []string{"foo", "bar"}
	for i := range mustResult {
		if mustResult[i] != called[i] {
			t.Fatalf("mw call order failed: %s != %s", mustResult[i], called[i])
		}
	}
}

func TestRouteDuplicatePanic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but got none")
		} else {
			t.Logf("panic value: %v", r)
		}
	}()

	lmux := NewLightMux(&http.Server{})
	lmux.NewRoute("/call1")
	lmux.NewRoute("/call2")
	lmux.NewRoute("/call1")
}

func TestRouteHandlerDuplicatePanic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but got none")
		} else {
			t.Logf("panic value: %v", r)
		}
	}()

	lmux := NewLightMux(&http.Server{})
	call1 := lmux.NewRoute("/foo")
	call1.Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {})
	call1.Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {})
}

func TestRouteDifferentHandler(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Log("Panic was not recieved")
		} else {
			t.Errorf("panic value: %v", r)
		}
	}()

	lmux := NewLightMux(&http.Server{})
	call1 := lmux.NewRoute("/api")
	call1.Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {})
	call1.Handle(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {})
}

func TestServerState(t *testing.T) {

	logMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("log")
			next.ServeHTTP(w, r)
		})
	}

	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	}

	lmux := NewLightMux(&http.Server{})
	lmux.Use(logMiddleware)
	route := lmux.NewRoute("/call", authMiddleware)
	route.Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	lmux.ApplyRoutes()
	lmux.ApplyGlobalMiddlewares()

	lmux.PrintRoutes()
	lmux.PrintMiddlewareInfo()
}

func Test404HandlerResponse(t *testing.T) {

	var called string

	foo := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		called = "foo"
	}

	bar := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		called = "bar"
	}

	lmux := NewLightMux(&http.Server{})
	route := lmux.NewRoute("/call")
	route.Handle(http.MethodGet, foo)
	lmux.Mux().HandleFunc("/", bar)

	lmux.ApplyRoutes()
	lmux.ApplyGlobalMiddlewares()

	req := httptest.NewRequest(http.MethodGet, "/random", nil)
	w := httptest.NewRecorder()
	lmux.Mux().ServeHTTP(w, req)

	t.Log(called)

	if called != "bar" {
		t.Fatalf("unexpected called value: %s, wanted: bar", called)
	}
}
