package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
	idgenserver "github.com/syumai/mcp/idgen-mcp-server/server"
)

func main() {
	srv := idgenserver.NewServer()
	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
}
