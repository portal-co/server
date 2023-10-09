package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
)

var discordOauthConfig = oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://discord.com/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	},
	Scopes:       []string{"identify"},
	ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
}

type ttokens struct{}

func NewDiscordCfg(x string) oauth2.Config {
	a := discordOauthConfig
	a.RedirectURL = x
	return a
}
func SetupAuth(r chi.Router, o, t string) func(http.Handler) string {
	var c oauth2.Config
	if t == "discord" {
		c = NewDiscordCfg(o + "/auth-callback/" + t)
	}
	m := map[string]http.Handler{}
	r.HandleFunc("/auth-callback/"+t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("state")
		n, ok := m[q]
		if !ok {
			return
		}
		delete(m, q)
		tokens, err := c.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), ttokens{}, tokens))
		n.ServeHTTP(w, r)
	})
	return func(h http.Handler) string {
		var x [32]byte
		rand.Read(x[:])
		s := base64.StdEncoding.EncodeToString(x[:])
		m[s] = h
		return c.AuthCodeURL(s)
	}
}
