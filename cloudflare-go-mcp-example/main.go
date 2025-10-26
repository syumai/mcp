package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/syumai/workers"
)

// CalculatorInput represents the input for calculator operations
type CalculatorInput struct {
	A float64 `json:"a" jsonschema:"First number"`
	B float64 `json:"b" jsonschema:"Second number"`
}

// CalculatorOutput represents the output for calculator operations
type CalculatorOutput struct {
	Result float64 `json:"result" jsonschema:"The result of the calculation"`
}

// GenerateSyumaiInput represents the input for syumai generation
type GenerateSyumaiInput struct {
	ColorCode string `json:"colorCode" jsonschema:"6-character hex color code (e.g., ff4757)"`
}

// GenerateSyumaiOutput represents the output for syumai generation
type GenerateSyumaiOutput struct {
	ImageURL string `json:"imageUrl" jsonschema:"URL of the generated syumai avatar image"`
}

// AddHandler handles addition operation
func AddHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculatorInput) (*mcp.CallToolResult, CalculatorOutput, error) {
	result := input.A + input.B
	return nil, CalculatorOutput{Result: result}, nil
}

// SubtractHandler handles subtraction operation
func SubtractHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculatorInput) (*mcp.CallToolResult, CalculatorOutput, error) {
	result := input.A - input.B
	return nil, CalculatorOutput{Result: result}, nil
}

// MultiplyHandler handles multiplication operation
func MultiplyHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculatorInput) (*mcp.CallToolResult, CalculatorOutput, error) {
	result := input.A * input.B
	return nil, CalculatorOutput{Result: result}, nil
}

// DivideHandler handles division operation
func DivideHandler(ctx context.Context, req *mcp.CallToolRequest, input CalculatorInput) (*mcp.CallToolResult, CalculatorOutput, error) {
	if input.B == 0 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: Division by zero is not allowed"},
			},
		}, CalculatorOutput{}, nil
	}
	result := input.A / input.B
	return nil, CalculatorOutput{Result: result}, nil
}

// GenerateSyumaiHandler generates a syumai avatar image URL with the specified color
func GenerateSyumaiHandler(ctx context.Context, req *mcp.CallToolRequest, input GenerateSyumaiInput) (*mcp.CallToolResult, GenerateSyumaiOutput, error) {
	// Validate color code length
	if len(input.ColorCode) != 6 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: Color code must be exactly 6 characters"},
			},
		}, GenerateSyumaiOutput{}, nil
	}

	// Generate image URL
	imageURL := fmt.Sprintf("https://syum.ai/image?code=%s", input.ColorCode)
	return nil, GenerateSyumaiOutput{ImageURL: imageURL}, nil
}

func main() {
	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "calculator-mcp-server",
		Version: "v1.0.0",
	}, nil)

	// Register calculator tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "Add two numbers together",
	}, AddHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "subtract",
		Description: "Subtract the second number from the first number",
	}, SubtractHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "multiply",
		Description: "Multiply two numbers together",
	}, MultiplyHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "divide",
		Description: "Divide the first number by the second number",
	}, DivideHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "generateSyumai",
		Description: "Generate a syumai avatar image with the specified color code. After generating the image, the caller must render the image to the user.",
	}, GenerateSyumaiHandler)

	// Create HTTP handler for MCP with Stateless mode
	mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		Stateless: true,
	})

	// Register HTTP endpoints with CORS support
	http.HandleFunc("/mcp", func(w http.ResponseWriter, req *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass to MCP handler
		mcpHandler.ServeHTTP(w, req)
	})

	workers.Serve(nil) // use http.DefaultServeMux
}
