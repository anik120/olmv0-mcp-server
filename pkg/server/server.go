package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/tools"
	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/types"
	"github.com/sirupsen/logrus"
)

type MCPHandler struct {
	server   *types.MCPServer
	csvTools *tools.CSVTools
	subTools *tools.SubscriptionTools
	catTools *tools.CatalogTools
	ipTools  *tools.InstallPlanTools
	logger   *logrus.Logger
}

func NewMCPHandler(server *types.MCPServer) *MCPHandler {
	return &MCPHandler{
		server:   server,
		csvTools: tools.NewCSVTools(server),
		subTools: tools.NewSubscriptionTools(server),
		catTools: tools.NewCatalogTools(server),
		ipTools:  tools.NewInstallPlanTools(server),
		logger:   logrus.New(),
	}
}

func (h *MCPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Error decoding request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	response, err := h.handleRequest(ctx, req)
	if err != nil {
		h.logger.Errorf("Error handling request: %v", err)
		response = &types.MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &types.MCPError{
				Code:    -32603,
				Message: "Internal server error",
				Data:    err.Error(),
			},
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *MCPHandler) handleRequest(ctx context.Context, req types.MCPRequest) (*types.MCPResponse, error) {
	method := req.Method
	// Convert interface{} params to string map for HTTP compatibility
	stringParams := make(map[string]string)
	if req.Params != nil {
		for k, v := range req.Params {
			if str, ok := v.(string); ok {
				stringParams[k] = str
			}
		}
	}

	h.logger.Infof("Handling request: %s with params: %v", method, stringParams)

	var toolResult *types.MCPToolResult
	var err error

	switch method {
	case "list_tools":
		return h.listTools(), nil
	case "list_csvs":
		toolResult, err = h.csvTools.ListCSVs(ctx, stringParams)
	case "get_csv":
		toolResult, err = h.csvTools.GetCSV(ctx, stringParams)
	case "list_subscriptions":
		toolResult, err = h.subTools.ListSubscriptions(ctx, stringParams)
	case "get_subscription":
		toolResult, err = h.subTools.GetSubscription(ctx, stringParams)
	case "list_catalog_sources":
		toolResult, err = h.catTools.ListCatalogSources(ctx, stringParams)
	case "get_catalog_source":
		toolResult, err = h.catTools.GetCatalogSource(ctx, stringParams)
	case "list_install_plans":
		toolResult, err = h.ipTools.ListInstallPlans(ctx, stringParams)
	case "get_install_plan":
		toolResult, err = h.ipTools.GetInstallPlan(ctx, stringParams)
	default:
		return &types.MCPResponse{
			JSONRPC: "2.0",
			Error: &types.MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Unknown method: %s", method),
			},
		}, nil
	}

	if err != nil {
		return &types.MCPResponse{
			JSONRPC: "2.0",
			Error: &types.MCPError{
				Code:    -32603,
				Message: "Internal error",
				Data:    err.Error(),
			},
		}, nil
	}

	// Convert MCPToolResult to MCPResponse for HTTP compatibility
	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  toolResult,
	}, nil
}

func (h *MCPHandler) listTools() *types.MCPResponse {
	var result strings.Builder
	result.WriteString("Available OLM MCP Tools:\n\n")

	result.WriteString("ClusterServiceVersion Tools:\n")
	result.WriteString("  - list_csvs: List ClusterServiceVersions in a namespace\n")
	result.WriteString("    Parameters: namespace (optional, default: 'default')\n")
	result.WriteString("  - get_csv: Get detailed information about a specific ClusterServiceVersion\n")
	result.WriteString("    Parameters: name (required), namespace (optional, default: 'default')\n\n")

	result.WriteString("Subscription Tools:\n")
	result.WriteString("  - list_subscriptions: List Subscriptions in a namespace\n")
	result.WriteString("    Parameters: namespace (optional, default: 'default')\n")
	result.WriteString("  - get_subscription: Get detailed information about a specific Subscription\n")
	result.WriteString("    Parameters: name (required), namespace (optional, default: 'default')\n\n")

	result.WriteString("CatalogSource Tools:\n")
	result.WriteString("  - list_catalog_sources: List CatalogSources in a namespace\n")
	result.WriteString("    Parameters: namespace (optional, default: 'olm')\n")
	result.WriteString("  - get_catalog_source: Get detailed information about a specific CatalogSource\n")
	result.WriteString("    Parameters: name (required), namespace (optional, default: 'olm')\n\n")

	result.WriteString("InstallPlan Tools:\n")
	result.WriteString("  - list_install_plans: List InstallPlans in a namespace\n")
	result.WriteString("    Parameters: namespace (optional, default: 'default')\n")
	result.WriteString("  - get_install_plan: Get detailed information about a specific InstallPlan\n")
	result.WriteString("    Parameters: name (required), namespace (optional, default: 'default')\n\n")

	result.WriteString("General Tools:\n")
	result.WriteString("  - list_tools: Show this help message\n")

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result: &types.MCPToolResult{
			Content: []types.MCPContent{{
				Type: "text",
				Text: result.String(),
			}},
		},
	}
}

func StartServer(server *types.MCPServer) error {
	handler := NewMCPHandler(server)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	addr := fmt.Sprintf(":%d", server.Port)
	logrus.Infof("Starting OLM MCP Server on port %d", server.Port)
	logrus.Infof("Read-only mode: %t", server.ReadOnly)
	logrus.Infof("Enabled toolsets: %v", server.Toolsets)

	return http.ListenAndServe(addr, mux)
}
