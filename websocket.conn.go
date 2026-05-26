package websocket

import (
	"bufio"
	"fmt"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// Conn represents an established WebSocket connection.
// It wraps the underlying net.Conn obtained after the HTTP upgrade
// and exposes a simple read/write API built on gobwas/ws.
// All methods are safe to call from a single goroutine per direction
// (one reader, one writer); use external synchronisation if multiple
// goroutines write concurrently.
type Conn struct {
	netConn net.Conn
	rw      *bufio.ReadWriter
}

// newConn wraps the hijacked net.Conn and the buffered ReadWriter produced
// by ws.UpgradeHTTP. The ReadWriter may contain bytes already read from the
// TCP stream during the HTTP handshake, so both its Reader and Writer must be
// used for subsequent WebSocket I/O.
func newConn(netConn net.Conn, rw *bufio.ReadWriter) *Conn {
	return &Conn{
		netConn: netConn,
		rw:      rw,
	}
}

// Read blocks until a complete message is received from the client and
// returns the payload together with the WebSocket operation code.
// Control frames (Ping, Pong, Close) are handled transparently by wsutil.
func (c *Conn) Read() ([]byte, ws.OpCode, error) {
	data, op, err := wsutil.ReadClientData(c.rw)
	if err != nil {
		return nil, op, fmt.Errorf("websocket: read: %w", err)
	}
	return data, op, nil
}

// Write sends data as a UTF-8 text message to the client.
func (c *Conn) Write(data []byte) error {
	if err := wsutil.WriteServerMessage(c.netConn, ws.OpText, data); err != nil {
		return fmt.Errorf("websocket: write text: %w", err)
	}
	return nil
}

// WriteString sends s as a UTF-8 text message to the client.
func (c *Conn) WriteString(s string) error {
	return c.Write([]byte(s))
}

// WriteBytes sends data as a binary message to the client.
func (c *Conn) WriteBytes(data []byte) error {
	if err := wsutil.WriteServerMessage(c.netConn, ws.OpBinary, data); err != nil {
		return fmt.Errorf("websocket: write binary: %w", err)
	}
	return nil
}

// Ping sends a Ping control frame to the client.
func (c *Conn) Ping(data []byte) error {
	if err := wsutil.WriteServerMessage(c.netConn, ws.OpPing, data); err != nil {
		return fmt.Errorf("websocket: ping: %w", err)
	}
	return nil
}

// Close sends a Close control frame and closes the underlying connection.
func (c *Conn) Close() error {
	// Best-effort Close frame; ignore write errors before closing.
	_ = wsutil.WriteServerMessage(c.netConn, ws.OpClose,
		ws.NewCloseFrameBody(ws.StatusNormalClosure, ""),
	)
	if err := c.netConn.Close(); err != nil {
		return fmt.Errorf("websocket: close: %w", err)
	}
	return nil
}

// CloseWithError sends a Close frame with the given status code and reason,
// then closes the underlying connection.
func (c *Conn) CloseWithError(code ws.StatusCode, reason string) error {
	_ = wsutil.WriteServerMessage(c.netConn, ws.OpClose,
		ws.NewCloseFrameBody(code, reason),
	)
	if err := c.netConn.Close(); err != nil {
		return fmt.Errorf("websocket: close: %w", err)
	}
	return nil
}

// NetConn returns the underlying net.Conn for advanced use.
// Callers must not read from or write to this connection using the ws package
// at the same time as other Conn methods.
func (c *Conn) NetConn() net.Conn { return c.netConn }

// RemoteAddr returns the remote network address of the client.
func (c *Conn) RemoteAddr() net.Addr { return c.netConn.RemoteAddr() }

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr { return c.netConn.LocalAddr() }
