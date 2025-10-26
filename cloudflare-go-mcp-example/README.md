# Cloudflare Workers MCP Server Example (Go)

A Calculator MCP server running on Cloudflare Workers using Go and the official Model Context Protocol Go SDK.

## About MCP

The Model Context Protocol (MCP) is an open protocol that standardizes how applications provide context to LLMs. MCP servers expose resources, tools, and prompts that LLM clients can use to enhance their capabilities.

## About This Project

This project demonstrates how to build and deploy an MCP server on Cloudflare Workers using Go. It uses the [`workers`](https://github.com/syumai/workers) package to run an HTTP server compiled to WebAssembly and the official [`go-sdk`](https://github.com/modelcontextprotocol/go-sdk) for MCP protocol implementation.

### Key Features

- **MCP server implementation in Go** using the official SDK
- **Runs on Cloudflare Workers** (edge computing)
- **WebAssembly-based deployment**
- **Four calculator tools** for basic arithmetic operations:
  - `add` - Addition of two numbers
  - `subtract` - Subtraction of two numbers
  - `multiply` - Multiplication of two numbers
  - `divide` - Division with zero-division error handling

## Notice

Go (not TinyGo) with many dependencies may exceed the size limit of the Worker (3MB for free plan, 10MB for paid plan). In that case, consider optimizing dependencies or using the TinyGo compiler.

## Requirements

- Node.js
- Go 1.24.0 or later
- A Cloudflare account (for deployment)

## Getting Started

### Installation

```console
# Clone or initialize the project
cd cloudflare-go-mcp-example

# Install Go dependencies
go mod tidy

# Install Node dependencies
pnpm install
```

### Development

Start the development server:

```console
pnpm start
# or
pnpm run dev
```

The MCP server will be available at `http://localhost:8787`.

### Testing the Server

#### Quick Test with Test Client

Run the test client to verify all calculator tools:

```console
go run scripts/test_client.go
```

This will test all four calculator operations and verify error handling.

#### Manual Testing with MCP Protocol

The server uses the Streamable HTTP transport which requires maintaining a session. For quick testing, use the provided test client above.

For information about the available endpoints:

```console
curl http://localhost:8787/hello
```

## Available Commands

```console
pnpm start       # Start dev server with Wrangler
pnpm run dev     # Same as start
pnpm run build   # Build Go Wasm binary
pnpm run deploy  # Deploy to Cloudflare Workers
go run .         # Run dev server without Wrangler (Cloudflare features unavailable)
```

## Connecting as an MCP Client

To connect to this MCP server from an MCP client (such as Claude Desktop), configure your client with the server URL:

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

After deployment to Cloudflare Workers, replace the URL with your deployed worker URL:

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

## Project Structure

- `main.go` - Main MCP server implementation with calculator tools
- `scripts/test_client.go` - Test client for verifying MCP server functionality
- `wrangler.jsonc` - Cloudflare Workers configuration
- `build/` - Compiled WebAssembly output

## Available Tools

The calculator MCP server provides four tools:

### add
Add two numbers together.

**Input:**
- `a` (number): First number
- `b` (number): Second number

**Output:**
- `result` (number): The sum of a and b

### subtract
Subtract the second number from the first number.

**Input:**
- `a` (number): First number
- `b` (number): Second number

**Output:**
- `result` (number): The difference (a - b)

### multiply
Multiply two numbers together.

**Input:**
- `a` (number): First number
- `b` (number): Second number

**Output:**
- `result` (number): The product of a and b

### divide
Divide the first number by the second number.

**Input:**
- `a` (number): Numerator
- `b` (number): Denominator

**Output:**
- `result` (number): The quotient (a / b)

**Error handling:** Returns an error if b is zero.

## Deployment

Deploy to Cloudflare Workers:

```console
pnpm run deploy
```

Your MCP server will be available at `https://your-worker-name.workers.dev`.

## Resources

- [Model Context Protocol Documentation](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Cloudflare Workers Documentation](https://developers.cloudflare.com/workers/)
- [syumai/workers Package](https://github.com/syumai/workers)

## License

MIT
