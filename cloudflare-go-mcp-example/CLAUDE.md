# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a fully functional MCP (Model Context Protocol) calculator server running on Cloudflare Workers using Go compiled to WebAssembly. The project demonstrates:
- MCP protocol implementation using the official `github.com/modelcontextprotocol/go-sdk` package
- Cloudflare Workers deployment via `github.com/syumai/workers` package
- Four calculator tools: `add`, `subtract`, `multiply`, `divide` (with division-by-zero error handling)

## Architecture

### WebAssembly Deployment Model

- Go code is compiled to WebAssembly (WASM) target `GOOS=js GOARCH=wasm`
- The compiled `app.wasm` binary is placed in `build/` directory
- Wrangler loads the WASM via `build/worker.mjs` (main entry point in wrangler.jsonc)
- The `workers.Serve(nil)` call bridges Go's `http.DefaultServeMux` to Cloudflare's Workers runtime

### MCP Server Architecture

The MCP implementation follows this pattern:
1. Create an MCP server with `mcp.NewServer()` specifying implementation name and version
2. Register tool handlers with `mcp.AddTool()` - each handler receives context, request, and typed input
3. Create HTTP handler using `mcp.NewStreamableHTTPHandler()` with `Stateless: true` option
4. Mount the handler at `/mcp` endpoint with CORS headers
5. Tool handlers return `(*mcp.CallToolResult, OutputType, error)` - use `CallToolResult` for errors, return nil for success

### Tool Handler Pattern

Tool handlers in main.go:265-108 follow this signature:
```go
func Handler(ctx context.Context, req *mcp.CallToolRequest, input InputType) (*mcp.CallToolResult, OutputType, error)
```

- Return `(nil, output, nil)` for successful execution
- Return `(&mcp.CallToolResult{IsError: true, Content: ...}, OutputType{}, nil)` for expected errors (like division by zero)
- Return `(nil, OutputType{}, err)` for unexpected errors

Input/output types use JSON schema tags for automatic schema generation.

## Development Commands

### Build
```bash
pnpm run build
```
Runs two steps:
1. `go run github.com/syumai/workers/cmd/workers-assets-gen -mode=go` - generates asset handling code
2. `GOOS=js GOARCH=wasm go build -o ./build/app.wasm .` - compiles Go to WebAssembly

### Development Server
```bash
pnpm start
# or
pnpm run dev
```
Starts Wrangler dev server on `http://localhost:8787`. The MCP endpoint is at `http://localhost:8787/mcp`.

### Testing
```bash
go run scripts/test_client.go
```
Runs in-memory test client that verifies all four calculator tools. Uses `mcp.NewInMemoryTransports()` to create client-server pair without HTTP.

### Deploy
```bash
pnpm run deploy
```
Deploys to Cloudflare Workers. The MCP endpoint will be at `https://your-worker-name.workers.dev/mcp`.

### Go Dependencies
```bash
go mod tidy
```

## Key Technical Constraints

### Size Limitations
Go binaries with many dependencies may exceed Cloudflare Workers size limits:
- Free plan: 3MB limit
- Paid plan: 10MB limit

If you exceed these limits:
- Optimize dependencies (remove unused packages)
- Use TinyGo compiler instead of standard Go
- Split functionality across multiple Workers

### Go Version
Requires Go 1.24.0 or later (go.mod specifies 1.25.0)

### CORS Configuration
The `/mcp` endpoint includes CORS headers allowing all origins (`Access-Control-Allow-Origin: *`). Restrict this in production if needed.

## MCP Client Configuration

Configure MCP clients (like Claude Desktop) to connect to this server:

**Local development:**
```json
{
  "mcpServers": {
    "cloudflare-go-example": {
      "url": "http://localhost:8787/mcp",
      "transport": "http"
    }
  }
}
```

**Production (after deployment):**
```json
{
  "mcpServers": {
    "cloudflare-go-example": {
      "url": "https://your-worker-name.workers.dev/mcp",
      "transport": "http"
    }
  }
}
```

## File Structure

- `main.go` - MCP server implementation with calculator tools
- `scripts/test_client.go` - In-memory test client (run with `go run`, has `+build ignore` tag)
- `wrangler.jsonc` - Cloudflare Workers configuration
- `build/` - Compiled WebAssembly output (generated)
