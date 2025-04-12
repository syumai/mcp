# idgen-mcp-server

idgen is a stdout-based MCP server that generates various types of IDs.

## Installation

```
go install github.com/syumai/mcp/idgen-mcp-server@latest
```

## Usage

```jsonc
{
    "idgen": {
        "command": "idgen-mcp-server",
    }
    // or
    "idgen": {
        "command": "go",
        "args": ["run", "github.com/syumai/mcp/idgen-mcp-server@latest"],
    }
}
```

## Available Tools

The following tools are available:

-   `generate_uuid`: Generates a new UUID (v4).
-   `generate_xid`: Generates a new XID.
-   `generate_ulid`: Generates a new ULID.
-   `generate_shortuuid`: Generates a new shortuuid.
