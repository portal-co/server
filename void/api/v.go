package api

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"net/netip"

	stream "github.com/nknorg/encrypted-stream"
	"nhooyr.io/websocket"
)

type Vl struct {
	*http.Client
	Ctx     context.Context
	Path    string
	Encrypt string
}

// Accept implements net.Listener.
func (v Vl) Accept() (net.Conn, error) {
	w, _, err := websocket.Dial(v.Ctx, v.Path, &websocket.DialOptions{HTTPClient: v.Client})
	if err != nil {
		return nil, err
	}
	n := websocket.NetConn(v.Ctx, w, websocket.MessageBinary)
	b, err := base64.StdEncoding.DecodeString(v.Encrypt)
	if err != nil {
		return nil, err
	}
	return stream.NewEncryptedStream(n, &stream.Config{
		Cipher:    stream.NewXSalsa20Poly1305Cipher((*[32]byte)(b)),
		Initiator: true,
	})
}

// Addr implements net.Listener.
func (Vl) Addr() net.Addr {
	return net.TCPAddrFromAddrPort(netip.MustParseAddrPort("127.0.0.1:80"))
}

// Close implements net.Listener.
func (Vl) Close() error {
	return nil
}

var _ net.Listener = Vl{}
