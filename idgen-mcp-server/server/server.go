package server

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/oklog/ulid/v2"
	"github.com/rs/xid"
)

// NewServer creates a new MCP server instance for idgen.
func NewServer() *server.MCPServer {
	srv := server.NewMCPServer(
		"idgen",
		"0.0.1",
		server.WithLogging(),
		server.WithRecovery(),
	)

	tools := []struct {
		name        string
		description string
		handler     func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	}{
		{
			name:        "generate_uuid",
			description: "Generates a new UUID (v4)",
			handler:     handleGenerateUUID,
		},
		{
			name:        "generate_xid",
			description: "Generates a new XID",
			handler:     handleGenerateXID,
		},
		{
			name:        "generate_ulid",
			description: "Generates a new ULID",
			handler:     handleGenerateULID,
		},
		{
			name:        "generate_shortuuid",
			description: "Generates a new shortuuid",
			handler:     handleGenerateShortUUID,
		},
	}
	for _, tool := range tools {
		t := mcp.NewTool(tool.name,
			mcp.WithDescription(tool.description),
		)
		srv.AddTool(t, tool.handler)
	}

	return srv
}

// handleGenerateUUID handles the generate_uuid tool request.
func handleGenerateUUID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	newUUID := uuid.NewString()
	return mcp.NewToolResultText(newUUID), nil
}

// handleGenerateXID handles the generate_xid tool request.
func handleGenerateXID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	newXID := xid.New().String()
	return mcp.NewToolResultText(newXID), nil
}

// handleGenerateULID handles the generate_ulid tool request.
func handleGenerateULID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	newULID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return mcp.NewToolResultText(newULID.String()), nil
}

// handleGenerateShortUUID handles the generate_shortuuid tool request.
func handleGenerateShortUUID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	newShortUUID := shortuuid.New()
	return mcp.NewToolResultText(newShortUUID), nil
}
