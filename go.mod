module github.com/nexgou/server/src/module/websocket

go 1.26.3

require (
	github.com/gobwas/ws v1.4.0
	github.com/nexgou/server v0.0.0
)

require (
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	golang.org/x/sys v0.45.0 // indirect
)

replace github.com/nexgou/server => ../../server
