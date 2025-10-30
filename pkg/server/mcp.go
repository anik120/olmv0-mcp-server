package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/tools"
	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/types"
	"github.com/sirupsen/logrus"
)

type MCPStdioServer struct {
	server   *types.MCPServer
	csvTools *tools.CSVTools
	subTools *tools.SubscriptionTools
	catTools *tools.CatalogTools
	ipTools  *tools.InstallPlanTools
	logger   *logrus.Logger
}

func NewMCPStdioServer(server *types.MCPServer) *MCPStdioServer {
	return &MCPStdioServer{
		server:   server,
		csvTools: tools.NewCSVTools(server),
		subTools: tools.NewSubscriptionTools(server),
		catTools: tools.NewCatalogTools(server),
		ipTools:  tools.NewInstallPlanTools(server),
		logger:   logrus.New(),
	}
}

func (s *MCPStdioServer) Start() error {
	s.logger.Info("Starting MCP stdio server")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		var req types.MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		response := s.handleRequest(context.Background(), req)

		responseBytes, err := json.Marshal(response)
		if err != nil {
			s.sendError(req.ID, -32603, "Internal error", err.Error())
			continue
		}

		fmt.Println(string(responseBytes))
	}

	return scanner.Err()
}

func (s *MCPStdioServer) handleRequest(ctx context.Context, req types.MCPRequest) *types.MCPResponse {
	s.logger.Infof("Handling MCP request: %s", req.Method)

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolCall(ctx, req)
	default:
		return &types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &types.MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func (s *MCPStdioServer) handleInitialize(req types.MCPRequest) *types.MCPResponse {
	return &types.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "olmv0-mcp-server",
				"version": "1.0.0",
			},
		},
	}
}

func (s *MCPStdioServer) handleToolsList(req types.MCPRequest) *types.MCPResponse {
	tools := []types.MCPTool{
		{
			Name:        "list_csvs",
			Description: "List ClusterServiceVersions in a namespace",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
			},
		},
		{
			Name:        "get_csv",
			Description: "Get detailed information about a specific ClusterServiceVersion",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the ClusterServiceVersion",
					},
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			Name:        "list_subscriptions",
			Description: "List Subscriptions in a namespace",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
			},
		},
		{
			Name:        "get_subscription",
			Description: "Get detailed information about a specific Subscription",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the Subscription",
					},
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			Name:        "list_catalog_sources",
			Description: "List CatalogSources in a namespace",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: olm)",
					},
				},
			},
		},
		{
			Name:        "get_catalog_source",
			Description: "Get detailed information about a specific CatalogSource",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the CatalogSource",
					},
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: olm)",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			Name:        "list_install_plans",
			Description: "List InstallPlans in a namespace",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
			},
		},
		{
			Name:        "get_install_plan",
			Description: "Get detailed information about a specific InstallPlan",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the InstallPlan",
					},
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Kubernetes namespace (default: default)",
					},
				},
				"required": []string{"name"},
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func (s *MCPStdioServer) handleToolCall(ctx context.Context, req types.MCPRequest) *types.MCPResponse {
	params, ok := req.Params["arguments"].(map[string]interface{})
	if !ok {
		return &types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	toolName, ok := req.Params["name"].(string)
	if !ok {
		return &types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "Tool name required",
			},
		}
	}

	// Convert params to string map for existing tool functions
	stringParams := make(map[string]string)
	for k, v := range params {
		if str, ok := v.(string); ok {
			stringParams[k] = str
		}
	}

	var result *types.MCPToolResult
	var err error

	switch toolName {
	case "list_csvs":
		result, err = s.csvTools.ListCSVs(ctx, stringParams)
	case "get_csv":
		result, err = s.csvTools.GetCSV(ctx, stringParams)
	case "list_subscriptions":
		result, err = s.subTools.ListSubscriptions(ctx, stringParams)
	case "get_subscription":
		result, err = s.subTools.GetSubscription(ctx, stringParams)
	case "list_catalog_sources":
		result, err = s.catTools.ListCatalogSources(ctx, stringParams)
	case "get_catalog_source":
		result, err = s.catTools.GetCatalogSource(ctx, stringParams)
	case "list_install_plans":
		result, err = s.ipTools.ListInstallPlans(ctx, stringParams)
	case "get_install_plan":
		result, err = s.ipTools.GetInstallPlan(ctx, stringParams)
	default:
		return &types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &types.MCPError{
				Code:    -32601,
				Message: "Unknown tool",
			},
		}
	}

	if err != nil {
		return &types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &types.MCPError{
				Code:    -32603,
				Message: "Tool execution failed",
				Data:    err.Error(),
			},
		}
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *MCPStdioServer) sendError(id interface{}, code int, message, data string) {
	response := &types.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &types.MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	responseBytes, _ := json.Marshal(response)
	fmt.Println(string(responseBytes))
}