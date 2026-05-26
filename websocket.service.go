// Package websocket provides a high-performance WebSocket module for Nexgou applications.
// It wraps github.com/gobwas/ws — a zero-allocation WebSocket library — and integrates
// cleanly with the Nexgou DI container and routing system.
//
// Configuration is read from environment variables:
//
//	WS_READ_BUFFER_SIZE  — read buffer size in bytes  (default: 4096)
//	WS_WRITE_BUFFER_SIZE — write buffer size in bytes (default: 4096)
//
// Usage:
//
//	func NewChatController(ws *websocket.WebSocketService) *ChatController {
//	    return &ChatController{ws: ws}
//	}
//
//	func (c *ChatController) Routes() []nexgou.Route {
//	    return []nexgou.Route{
//	        nexgou.Get("/ws/chat", c.ws.Handler(c.handleConn)),
//	    }
//	}
//
//	func (c *ChatController) handleConn(conn *websocket.Conn) {
//	    defer conn.Close()
//	    for {
//	        data, _, err := conn.Read()
//	        if err != nil {
//	            return
//	        }
//	        conn.Write(data) // echo
//	    }
//	}
package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/nexgou/server/src/common"
	"github.com/nexgou/server/src/config"
	"github.com/nexgou/server/src/logger"
)

// WebSocketService handles HTTP→WebSocket upgrades using gobwas/ws.
// It is safe for concurrent use; each upgrade produces an independent *Conn.
type WebSocketService struct {
	upgrader        ws.HTTPUpgrader
	readBufferSize  int
	writeBufferSize int
	log             *logger.ScopedLogger
}

// NewWebSocketService creates a new WebSocketService.
// Depends on *config.ConfigService and *logger.LoggerService.
func NewWebSocketService(cfg *config.ConfigService, log *logger.LoggerService) *WebSocketService {
	readBuf := cfg.GetInt("WS_READ_BUFFER_SIZE", 4096)
	writeBuf := cfg.GetInt("WS_WRITE_BUFFER_SIZE", 4096)

	svc := &WebSocketService{
		readBufferSize:  readBuf,
		writeBufferSize: writeBuf,
		log:             log.WithContext("WebSocketService"),
	}

	svc.upgrader = ws.HTTPUpgrader{
		Timeout: 5 * time.Second,
		Header:  http.Header{"X-Powered-By": []string{"Nexgou/WebSocket"}},
	}

	svc.log.Info("ready",
		"read_buffer", readBuf,
		"write_buffer", writeBuf,
	)
	return svc
}

// Upgrade performs the HTTP→WebSocket handshake on the current Nexgou request
// and returns a *Conn ready for reading and writing.
// The caller is responsible for closing the connection when done.
//
// Example:
//
//	conn, err := svc.Upgrade(ctx)
//	if err != nil {
//	    return err
//	}
//	defer conn.Close()
//	// ... read/write loop
func (s *WebSocketService) Upgrade(ctx *common.Context) (*Conn, error) {
	netConn, rw, _, err := s.upgrader.Upgrade(ctx.Request, ctx.Writer)
	if err != nil {
		return nil, fmt.Errorf("websocket: upgrade: %w", err)
	}

	conn := newConn(netConn, rw)
	s.log.Debug("upgraded", "remote", netConn.RemoteAddr().String())
	return conn, nil
}

// Handler returns a Nexgou HandlerFunc that upgrades the connection and invokes
// fn in the same goroutine. fn should run a read/write loop and return when the
// connection is done; the handler will return nil once fn returns.
//
// Example:
//
//	nexgou.Get("/ws", wsvc.Handler(func(conn *websocket.Conn) {
//	    defer conn.Close()
//	    for {
//	        data, _, err := conn.Read()
//	        if err != nil { return }
//	        conn.Write(data)
//	    }
//	}))
func (s *WebSocketService) Handler(fn func(*Conn)) common.HandlerFunc {
	return func(ctx *common.Context) error {
		conn, err := s.Upgrade(ctx)
		if err != nil {
			return err
		}
		fn(conn)
		return nil
	}
}

// HandlerAsync is like Handler but runs fn in a new goroutine, returning
// immediately after the upgrade. Use this when you want the HTTP handler
// to return quickly and manage the connection lifecycle yourself.
//
// Example:
//
//	nexgou.Get("/ws", wsvc.HandlerAsync(func(conn *websocket.Conn) {
//	    defer conn.Close()
//	    // long-running loop here
//	}))
func (s *WebSocketService) HandlerAsync(fn func(*Conn)) common.HandlerFunc {
	return func(ctx *common.Context) error {
		conn, err := s.Upgrade(ctx)
		if err != nil {
			return err
		}
		go fn(conn)
		return nil
	}
}
