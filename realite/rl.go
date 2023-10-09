package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	discord, err := discordgo.New(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	discord.Identify.Intents |= discordgo.IntentMessageContent
	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}
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
	auths := map[string]func(http.Handler) string{}
	auths["discord"] = SetupAuth(r, Host, "discord")

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi world"))
	})
	// r.Connect("/adm/vpn", func(w http.ResponseWriter, r *http.Request) {
	// 	proxyConnect(w, r)
	// })
	for {
		err = http.Serve(netutil.LimitListener(Vd, *limit), r)
		if err != nil {
			fmt.Println(err)
		}
	}
}
