package websocket

import (
	"github.com/nexgou/server/src/common"
	"github.com/nexgou/server/src/config"
	"github.com/nexgou/server/src/core"
	"github.com/nexgou/server/src/logger"
)

// Module is a ready-to-use Nexgou module that registers and exports WebSocketService.
// It uses github.com/gobwas/ws for zero-allocation, high-performance WebSocket handling.
//
// Optional environment variables:
//
//	WS_READ_BUFFER_SIZE  = 4096   (bytes, default)
//	WS_WRITE_BUFFER_SIZE = 4096   (bytes, default)
//
// Usage:
//
//	var AppModule = nexgou.Module(nexgou.ModuleOptions{
//	    Imports: []nexgou.IModule{
//	        nexgou.ConfigModule,
//	        nexgou.LogModule,
//	        websocket.Module,
//	    },
//	})
//
//	func NewChatController(ws *websocket.WebSocketService) *ChatController {
//	    return &ChatController{ws: ws}
//	}
//
//	// In your controller Routes():
//	//   nexgou.Get("/ws", ws.Handler(func(conn *websocket.Conn) { ... }))
var Module common.IModule = core.NewModule(common.ModuleOptions{
	Imports: []common.IModule{
		config.ConfigModule,
		logger.LogModule,
	},
	Providers: []any{NewWebSocketService},
	Exports:   []any{NewWebSocketService},
})
