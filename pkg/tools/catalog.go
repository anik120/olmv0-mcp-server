package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/types"
)

type CatalogTools struct {
	server *types.MCPServer
}

func NewCatalogTools(server *types.MCPServer) *CatalogTools {
	return &CatalogTools{server: server}
}

func (t *CatalogTools) ListCatalogSources(ctx context.Context, params map[string]string) (*types.MCPToolResult, error) {
	namespace := params["namespace"]
	if namespace == "" {
		namespace = "olm"
	}

	catalogs, err := t.server.OLMClient.ListCatalogSources(ctx, namespace)
	if err != nil {
		return &types.MCPToolResult{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error listing CatalogSources: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("CatalogSources in namespace '%s':\n\n", namespace))

	if len(catalogs.Items) == 0 {
		result.WriteString("No CatalogSources found.\n")
	} else {
		result.WriteString("NAME\tNAMESPACE\tSOURCE TYPE\tDISPLAY NAME\tSTATE\n")
		for _, cat := range catalogs.Items {
			result.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
				cat.Name,
				cat.Namespace,
				cat.Spec.SourceType,
				cat.Spec.DisplayName,
				cat.Status.GRPCConnectionState.LastObservedState,
			))
		}
	}

	return &types.MCPToolResult{
		Content: []types.MCPContent{{
			Type: "text",
			Text: result.String(),
		}},
	}, nil
}

func (t *CatalogTools) GetCatalogSource(ctx context.Context, params map[string]string) (*types.MCPToolResult, error) {
	namespace := params["namespace"]
	name := params["name"]

	if namespace == "" {
		namespace = "olm"
	}
	if name == "" {
		return &types.MCPToolResult{
			Content: []types.MCPContent{{
				Type: "text",
				Text: "Error: 'name' parameter is required",
			}},
			IsError: true,
		}, nil
	}

	catalog, err := t.server.OLMClient.GetCatalogSource(ctx, namespace, name)
	if err != nil {
		return &types.MCPToolResult{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error getting CatalogSource '%s': %v", name, err),
			}},
			IsError: true,
		}, nil
	}

	jsonData, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return &types.MCPToolResult{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error marshaling CatalogSource to JSON: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("CatalogSource: %s/%s\n\n", namespace, name))
	result.WriteString("Basic Info:\n")
	result.WriteString(fmt.Sprintf("  Name: %s\n", catalog.Name))
	result.WriteString(fmt.Sprintf("  Namespace: %s\n", catalog.Namespace))
	result.WriteString(fmt.Sprintf("  Display Name: %s\n", catalog.Spec.DisplayName))
	result.WriteString(fmt.Sprintf("  Source Type: %s\n", catalog.Spec.SourceType))
	result.WriteString(fmt.Sprintf("  Publisher: %s\n", catalog.Spec.Publisher))
	result.WriteString(fmt.Sprintf("  Connection State: %s\n", catalog.Status.GRPCConnectionState.LastObservedState))
	if !catalog.Status.GRPCConnectionState.LastConnectTime.IsZero() {
		result.WriteString(fmt.Sprintf("  Last Observed: %s\n", catalog.Status.GRPCConnectionState.LastConnectTime.String()))
	}
	result.WriteString("\n")

	if catalog.Spec.SourceType == "grpc" {
		result.WriteString(fmt.Sprintf("  Image: %s\n", catalog.Spec.Image))
		result.WriteString(fmt.Sprintf("  Address: %s\n", catalog.Spec.Address))
	}

	result.WriteString("\nFull JSON representation:\n")
	result.WriteString("```json\n")
	result.WriteString(string(jsonData))
	result.WriteString("\n```\n")

	return &types.MCPToolResult{
		Content: []types.MCPContent{{
			Type: "text",
			Text: result.String(),
		}},
	}, nil
}
