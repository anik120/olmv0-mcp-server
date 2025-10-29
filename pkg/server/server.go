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
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Internal server error: %v", err),
			}},
			IsError: true,
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *MCPHandler) handleRequest(ctx context.Context, req types.MCPRequest) (*types.MCPResponse, error) {
	method := req.Method
	params := req.Params

	h.logger.Infof("Handling request: %s with params: %v", method, params)

	switch method {
	case "list_tools":
		return h.listTools(), nil
	case "list_csvs":
		return h.csvTools.ListCSVs(ctx, params)
	case "get_csv":
		return h.csvTools.GetCSV(ctx, params)
	case "list_subscriptions":
		return h.subTools.ListSubscriptions(ctx, params)
	case "get_subscription":
		return h.subTools.GetSubscription(ctx, params)
	case "list_catalog_sources":
		return h.catTools.ListCatalogSources(ctx, params)
	case "get_catalog_source":
		return h.catTools.GetCatalogSource(ctx, params)
	case "list_install_plans":
		return h.ipTools.ListInstallPlans(ctx, params)
	case "get_install_plan":
		return h.ipTools.GetInstallPlan(ctx, params)
	default:
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Unknown method: %s", method),
			}},
			IsError: true,
		}, nil
	}
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
		Content: []types.MCPContent{{
			Type: "text",
			Text: result.String(),
		}},
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
