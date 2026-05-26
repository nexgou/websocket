> **[English](#websocket-module)** | **[Español](#módulo-websocket)**

---

# WebSocket Module

A high-performance WebSocket module for [Nexgou](https://github.com/nexgou/server) applications, built on top of [gobwas/ws](https://github.com/gobwas/ws) — a zero-allocation WebSocket library.

## When to use

- Real-time bidirectional communication (chat, live dashboards, notifications)
- Low-latency streaming where HTTP/SSE is not enough
- You need full-duplex channels per client connection
- High-concurrency scenarios that require minimal allocations per frame

## Installation

```bash
go get github.com/nexgou/server/src/module/websocket
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `WS_READ_BUFFER_SIZE` | `4096` | Read buffer size in bytes |
| `WS_WRITE_BUFFER_SIZE` | `4096` | Write buffer size in bytes |

No external service is required.

## Module wiring

```go
import "github.com/nexgou/server/src/module/websocket"

var AppModule = nexgou.Module(nexgou.ModuleOptions{
    Imports: []nexgou.IModule{
        nexgou.ConfigModule,
        nexgou.LogModule,
        websocket.Module,
    },
})
```

## API Reference

### `WebSocketService`

| Method | Signature | Description |
|---|---|---|
| `Upgrade` | `(ctx *common.Context) (*Conn, error)` | Upgrades an HTTP request to a WebSocket connection |
| `Handler` | `(fn func(*Conn)) common.HandlerFunc` | Returns a route handler that upgrades and calls fn synchronously |
| `HandlerAsync` | `(fn func(*Conn)) common.HandlerFunc` | Same as Handler but runs fn in a new goroutine |

### `Conn`

| Method | Signature | Description |
|---|---|---|
| `Read` | `() ([]byte, ws.OpCode, error)` | Reads the next message from the client |
| `Write` | `(data []byte) error` | Sends a UTF-8 text message |
| `WriteString` | `(s string) error` | Sends a UTF-8 text message from a string |
| `WriteBytes` | `(data []byte) error` | Sends a binary message |
| `Ping` | `(data []byte) error` | Sends a Ping control frame |
| `Close` | `() error` | Sends Close frame and closes the connection |
| `CloseWithError` | `(code ws.StatusCode, reason string) error` | Closes with a specific status code |
| `RemoteAddr` | `() net.Addr` | Returns the client's remote address |
| `LocalAddr` | `() net.Addr` | Returns the local address |
| `NetConn` | `() net.Conn` | Returns the raw net.Conn (escape hatch) |

## Usage examples

### Echo server

```go
func NewEchoController(ws *websocket.WebSocketService) *EchoController {
    return &EchoController{ws: ws}
}

func (c *EchoController) Routes() []nexgou.Route {
    return []nexgou.Route{
        nexgou.Get("/ws/echo", c.ws.Handler(c.echo)),
    }
}

func (c *EchoController) echo(conn *websocket.Conn) {
    defer conn.Close()
    for {
        data, _, err := conn.Read()
        if err != nil {
            return
        }
        if err := conn.Write(data); err != nil {
            return
        }
    }
}
```

### Chat broadcast (async upgrade)

```go
type ChatController struct {
    ws      *websocket.WebSocketService
    clients sync.Map
}

func (c *ChatController) Routes() []nexgou.Route {
    return []nexgou.Route{
        nexgou.Get("/ws/chat", c.ws.HandlerAsync(c.handleClient)),
    }
}

func (c *ChatController) handleClient(conn *websocket.Conn) {
    defer func() {
        c.clients.Delete(conn)
        conn.Close()
    }()
    c.clients.Store(conn, struct{}{})

    for {
        data, _, err := conn.Read()
        if err != nil {
            return
        }
        // broadcast to all
        c.clients.Range(func(k, _ any) bool {
            k.(*websocket.Conn).Write(data)
            return true
        })
    }
}
```

### Manual upgrade with custom logic

```go
func (c *MyController) wsHandler(ctx *nexgou.Context) error {
    conn, err := c.ws.Upgrade(ctx)
    if err != nil {
        return err
    }
    defer conn.Close()

    for {
        data, opCode, err := conn.Read()
        if err != nil {
            return nil
        }
        switch opCode {
        case ws.OpText:
            conn.WriteString("received: " + string(data))
        case ws.OpBinary:
            conn.WriteBytes(data)
        }
    }
}
```

---

---

# Módulo WebSocket

Módulo WebSocket de alto rendimiento para aplicaciones [Nexgou](https://github.com/nexgou/server), construido sobre [gobwas/ws](https://github.com/gobwas/ws) — una librería WebSocket de cero asignaciones.

## Cuándo usarlo

- Comunicación bidireccional en tiempo real (chat, dashboards en vivo, notificaciones)
- Streaming de baja latencia donde HTTP/SSE no es suficiente
- Necesitas canales full-duplex por conexión de cliente
- Escenarios de alta concurrencia que requieren mínimas asignaciones por frame

## Instalación

```bash
go get github.com/nexgou/server/src/module/websocket
```

## Configuración

| Variable | Por defecto | Descripción |
|---|---|---|
| `WS_READ_BUFFER_SIZE` | `4096` | Tamaño del buffer de lectura en bytes |
| `WS_WRITE_BUFFER_SIZE` | `4096` | Tamaño del buffer de escritura en bytes |

No se requiere ningún servicio externo.

## Registro del módulo

```go
import "github.com/nexgou/server/src/module/websocket"

var AppModule = nexgou.Module(nexgou.ModuleOptions{
    Imports: []nexgou.IModule{
        nexgou.ConfigModule,
        nexgou.LogModule,
        websocket.Module,
    },
})
```

## Referencia de la API

### `WebSocketService`

| Método | Firma | Descripción |
|---|---|---|
| `Upgrade` | `(ctx *common.Context) (*Conn, error)` | Hace upgrade de una petición HTTP a una conexión WebSocket |
| `Handler` | `(fn func(*Conn)) common.HandlerFunc` | Devuelve un handler que hace upgrade y llama fn síncronamente |
| `HandlerAsync` | `(fn func(*Conn)) common.HandlerFunc` | Como Handler pero ejecuta fn en una nueva goroutine |

### `Conn`

| Método | Firma | Descripción |
|---|---|---|
| `Read` | `() ([]byte, ws.OpCode, error)` | Lee el siguiente mensaje del cliente |
| `Write` | `(data []byte) error` | Envía un mensaje de texto UTF-8 |
| `WriteString` | `(s string) error` | Envía un mensaje de texto UTF-8 desde un string |
| `WriteBytes` | `(data []byte) error` | Envía un mensaje binario |
| `Ping` | `(data []byte) error` | Envía un frame de control Ping |
| `Close` | `() error` | Envía frame Close y cierra la conexión |
| `CloseWithError` | `(code ws.StatusCode, reason string) error` | Cierra con un código de estado específico |
| `RemoteAddr` | `() net.Addr` | Devuelve la dirección remota del cliente |
| `LocalAddr` | `() net.Addr` | Devuelve la dirección local |
| `NetConn` | `() net.Conn` | Devuelve el net.Conn subyacente (escape hatch) |

## Ejemplos de uso

### Servidor echo

```go
func NewEchoController(ws *websocket.WebSocketService) *EchoController {
    return &EchoController{ws: ws}
}

func (c *EchoController) Routes() []nexgou.Route {
    return []nexgou.Route{
        nexgou.Get("/ws/echo", c.ws.Handler(c.echo)),
    }
}

func (c *EchoController) echo(conn *websocket.Conn) {
    defer conn.Close()
    for {
        data, _, err := conn.Read()
        if err != nil {
            return
        }
        conn.Write(data)
    }
}
```

### Chat con broadcast (upgrade asíncrono)

```go
type ChatController struct {
    ws      *websocket.WebSocketService
    clients sync.Map
}

func (c *ChatController) Routes() []nexgou.Route {
    return []nexgou.Route{
        nexgou.Get("/ws/chat", c.ws.HandlerAsync(c.handleClient)),
    }
}

func (c *ChatController) handleClient(conn *websocket.Conn) {
    defer func() {
        c.clients.Delete(conn)
        conn.Close()
    }()
    c.clients.Store(conn, struct{}{})

    for {
        data, _, err := conn.Read()
        if err != nil {
            return
        }
        c.clients.Range(func(k, _ any) bool {
            k.(*websocket.Conn).Write(data)
            return true
        })
    }
}
```

### Upgrade manual con lógica personalizada

```go
func (c *MiControlador) wsHandler(ctx *nexgou.Context) error {
    conn, err := c.ws.Upgrade(ctx)
    if err != nil {
        return err
    }
    defer conn.Close()

    for {
        data, opCode, err := conn.Read()
        if err != nil {
            return nil
        }
        switch opCode {
        case ws.OpText:
            conn.WriteString("recibido: " + string(data))
        case ws.OpBinary:
            conn.WriteBytes(data)
        }
    }
}
```
