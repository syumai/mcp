// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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

func main() {
	testClient()
}

func testClient() {
	ctx := context.Background()

	// Create test server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "calculator-mcp-server",
		Version: "v1.0.0",
	}, nil)

	// Register tools
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

	// Create in-memory transports
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	// Connect server
	go func() {
		if err := server.Run(ctx, serverTransport); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Create and connect client
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	fmt.Println("✓ Server initialized successfully")

	// List tools
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("✓ Found %d tools:\n", len(tools.Tools))
	for _, tool := range tools.Tools {
		fmt.Printf("  - %s: %s\n", tool.Name, tool.Description)
	}

	// Test add
	fmt.Println("\nTesting tools:")
	addResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "add",
		Arguments: map[string]any{"a": 10.0, "b": 5.0},
	})
	if err != nil {
		log.Fatalf("Failed to call add: %v", err)
	}
	fmt.Printf("✓ add(10, 5) = %v\n", addResult.Content)

	// Test subtract
	subResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "subtract",
		Arguments: map[string]any{"a": 10.0, "b": 3.0},
	})
	if err != nil {
		log.Fatalf("Failed to call subtract: %v", err)
	}
	fmt.Printf("✓ subtract(10, 3) = %v\n", subResult.Content)

	// Test multiply
	mulResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "multiply",
		Arguments: map[string]any{"a": 4.0, "b": 7.0},
	})
	if err != nil {
		log.Fatalf("Failed to call multiply: %v", err)
	}
	fmt.Printf("✓ multiply(4, 7) = %v\n", mulResult.Content)

	// Test divide
	divResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "divide",
		Arguments: map[string]any{"a": 20.0, "b": 4.0},
	})
	if err != nil {
		log.Fatalf("Failed to call divide: %v", err)
	}
	fmt.Printf("✓ divide(20, 4) = %v\n", divResult.Content)

	// Test divide by zero
	divZeroResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "divide",
		Arguments: map[string]any{"a": 10.0, "b": 0.0},
	})
	if err == nil && divZeroResult.IsError {
		fmt.Printf("✓ divide(10, 0) correctly returns error: %v\n", divZeroResult.Content)
	} else if err != nil {
		log.Fatalf("Failed to call divide with zero: %v", err)
	}

	fmt.Println("\n✓ All tests passed!")
}
