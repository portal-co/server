package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
)

func SpaceFieldsJoin(str string) string {
	return strings.Join(strings.Fields(str), "")
}

func main() {
	u := "https://raw.githubusercontent.com/portal-co/server/main/pinned"
	c, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	defer c.Body.Close()
	a, err := io.ReadAll(c.Body)
	if err != nil {
		panic(err)
	}
	d, err := http.Get("https://ipfs.io/" + SpaceFieldsJoin(string(a)))
	if err != nil {
		panic(err)
	}
	defer d.Body.Close()
	f, err := os.CreateTemp("/tmp", "prtl-srvr-")
	if err != nil {
		panic(err)
	}
	p := f.Name()
	_, err = io.Copy(f, d.Body)
	if err != nil {
		f.Close()
		panic(err)
	}
	f.Close()
	err = os.Chmod(p, 0777)
	if err != nil {
		panic(err)
	}
	for {
		syscall.Exec(p, []string{p}, os.Environ())
	}
}
