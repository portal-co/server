package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/portal-co/server/void/api"
	"golang.org/x/net/netutil"
)

var Vd = api.Vl{
	nil, context.Background(), os.Getenv("REALITE_URL"), os.Getenv("ENCRYPT"),
}

var Host = strings.TrimSuffix(os.Getenv("REALITE_URL"), "/void.ws")

var limit = flag.Int("rate_limit", 256, "The rate limit")

func main() {
	r := chi.NewMux()
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.UserAgent(), "ByteLo") {
				if strings.Contains(r.UserAgent(), "Channel/g") {

				} else {
					r.URL.Scheme = "googlechromes"
					http.Redirect(w, r, r.URL.String(), http.StatusFound)
				}
			} else {
				h.ServeHTTP(w, r)
			}
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi world"))
	})

	err := http.Serve(netutil.LimitListener(Vd, *limit), r)
	if err != nil {
		log.Fatal(err)
	}
}
