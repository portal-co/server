package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	stream "github.com/nknorg/encrypted-stream"
	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

type R struct {
	net.Conn
	Defer func()
}

func main() {
	fmt.Println("void:starting")
	b, err := base64.StdEncoding.DecodeString(os.Getenv("ENCRYPT"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		rs := make(chan R)
		try := true
		http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" && try {
				try = false
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			if strings.Contains(r.UserAgent(), "ytespider") {
				return
			}
			if r.URL.Path == "/void.ws" {
				ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
				if err != nil {
					return
				}
				n := websocket.NetConn(r.Context(), ws, websocket.MessageBinary)
				defer n.Close()
				n, err = stream.NewEncryptedStream(n, &stream.Config{
					Cipher: stream.NewXSalsa20Poly1305Cipher((*[32]byte)(b)),
				})
				if err != nil {
					return
				}
				n.SetDeadline(time.Now().Add(3 * time.Second))
				c := make(chan struct{})
				rs <- R{n, func() {
					c <- struct{}{}
				}}
				<-c
				return
			}
			x, y, err := w.(http.Hijacker).Hijack()
			if err != nil {
				http.Error(w, "void: "+err.Error(), http.StatusInternalServerError)
				return
			}
			for {
				s := <-rs
				s.SetDeadline(time.Time{})
				err = r.Write(s)
				if err != nil {
					fmt.Println("void:", err)
					continue
				}
				var g errgroup.Group
				g.Go(func() error {
					defer s.Close()
					_, err := io.Copy(s, io.MultiReader(y, x))
					return err
				})
				g.Go(func() error {
					defer x.Close()
					_, err := io.Copy(x, s)
					return err
				})
				err = g.Wait()
				s.Defer()
				if err != nil {
					fmt.Println(err)

				} else {
					break
				}
			}
		}))
	}
}
