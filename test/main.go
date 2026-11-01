package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ayayaakasvin/lightmux"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	lst := []string{"a", "b", "c"}
	lm := lightmux.NewLightMux(&http.Server{
		Addr: "localhost:8080",
	})

	apiGroup := lm.NewGroup("/api", printHelloMiddleware)
	apiGroup.NewRoute("/list").Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strings.Join(lst, ", ")))
	})
	apiGroup.NewRoute("/asd", func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("List POST triggered %s", r.RemoteAddr)
			hf.ServeHTTP(w, r)
		}
	}).Handle(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		lst = append(lst, "d")
		w.Write([]byte(strings.Join(lst, ", ")))
	})

	pingFunction := func (w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte(fmt.Sprintf("PONG %s", r.RemoteAddr)))
	}

	apiGroup.NewRoute("/ping").Handle("GET", pingFunction)

	lm.NewRoute("/writeA").Handle(http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{'a'})
	})

	lm.PrintMiddlewareInfo()
	lm.PrintRoutes()
	if err := lm.Run(ctx); err != nil {
		log.Fatalf("lightmux error: %v", err)
	}
}

func printHelloMiddleware (next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("hello epta %s", r.URL.Path)
		next.ServeHTTP(w, r)
	}
}