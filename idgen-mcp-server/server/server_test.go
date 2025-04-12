package server

import (
	"regexp"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/oklog/ulid/v2"
	"github.com/rs/xid"
)

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

func TestTools(t *testing.T) {
	mcpServer := NewServer()
	testServer := server.NewTestServer(mcpServer)
	t.Cleanup(testServer.Close)

	testCases := map[string]func(t *testing.T, id string){
		"generate_uuid": func(t *testing.T, id string) {
			if id == "" {
				t.Errorf("Expected non-empty UUID, got empty string")
				return
			}
			if !uuidRegex.MatchString(id) {
				t.Errorf("Expected result text to be a valid UUID v4, but got: %s", id)
			}
		},
		"generate_xid": func(t *testing.T, id string) {
			_, err := xid.FromString(id)
			if err != nil {
				t.Errorf("Expected result text to be a valid XID, but got error: %v. Value: %s", err, id)
			}
		},
		"generate_ulid": func(t *testing.T, id string) {
			if len(id) != 26 {
				t.Errorf("Expected ULID string length 26, but got %d. Value: %s", len(id), id)
				return
			}
			_, err := ulid.Parse(id)
			if err != nil {
				t.Errorf("Expected result text to be a valid ULID, but got error: %v. Value: %s", err, id)
			}
		},
		"generate_shortuuid": func(t *testing.T, id string) {
			if id == "" {
				t.Errorf("Expected non-empty result text for ShortUUID, but got empty string")
				return
			}
			if len(id) != 22 {
				t.Errorf("Expected shortuuid string length 22, but got %d. Value: %s", len(id), id)
			}
		},
	}

	cli, err := client.NewSSEMCPClient(testServer.URL + "/sse")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	t.Cleanup(func() {
		if err := cli.Close(); err != nil {
			t.Errorf("Failed to close client: %v", err)
		}
	})

	if err := cli.Start(t.Context()); err != nil {
		t.Fatalf("Failed to start client: %v", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(t.Context(), initRequest)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	for toolName, validateID := range testCases {
		t.Run(toolName, func(t *testing.T) {
			t.Parallel()

			var request mcp.CallToolRequest
			request.Params.Name = toolName

			result, err := cli.CallTool(t.Context(), request)
			if err != nil {
				t.Fatalf("CallTool for %s failed: %v", toolName, err)
			}
			if result == nil {
				t.Fatal("want a result, but got nil")
			}
			if len(result.Content) == 0 {
				t.Fatalf("want at least one content element, got none. %+v", result)
			}
			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatalf("want result.Content[0] to be mcp.TextContent, got %T. %+v", result.Content[0], result)
			}
			validateID(t, textContent.Text)
		})
	}
}
