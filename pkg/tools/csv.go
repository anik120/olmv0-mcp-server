package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/types"
)

type CSVTools struct {
	server *types.MCPServer
}

func NewCSVTools(server *types.MCPServer) *CSVTools {
	return &CSVTools{server: server}
}

func (t *CSVTools) ListCSVs(ctx context.Context, params map[string]string) (*types.MCPResponse, error) {
	namespace := params["namespace"]
	if namespace == "" {
		namespace = "default"
	}

	csvs, err := t.server.OLMClient.ListClusterServiceVersions(ctx, namespace)
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error listing ClusterServiceVersions: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("ClusterServiceVersions in namespace '%s':\n\n", namespace))

	if len(csvs.Items) == 0 {
		result.WriteString("No ClusterServiceVersions found.\n")
	} else {
		result.WriteString("NAME\tNAMESPACE\tPHASE\tVERSION\tREPLACES\n")
		for _, csv := range csvs.Items {
			result.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
				csv.Name,
				csv.Namespace,
				csv.Status.Phase,
				csv.Spec.Version.String(),
				csv.Spec.Replaces,
			))
		}
	}

	return &types.MCPResponse{
		Content: []types.MCPContent{{
			Type: "text",
			Text: result.String(),
		}},
	}, nil
}

func (t *CSVTools) GetCSV(ctx context.Context, params map[string]string) (*types.MCPResponse, error) {
	namespace := params["namespace"]
	name := params["name"]

	if namespace == "" {
		namespace = "default"
	}
	if name == "" {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: "Error: 'name' parameter is required",
			}},
			IsError: true,
		}, nil
	}

	csv, err := t.server.OLMClient.GetClusterServiceVersion(ctx, namespace, name)
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error getting ClusterServiceVersion '%s': %v", name, err),
			}},
			IsError: true,
		}, nil
	}

	jsonData, err := json.MarshalIndent(csv, "", "  ")
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error marshaling CSV to JSON: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("ClusterServiceVersion: %s/%s\n\n", namespace, name))
	result.WriteString("Basic Info:\n")
	result.WriteString(fmt.Sprintf("  Name: %s\n", csv.Name))
	result.WriteString(fmt.Sprintf("  Namespace: %s\n", csv.Namespace))
	result.WriteString(fmt.Sprintf("  Phase: %s\n", csv.Status.Phase))
	result.WriteString(fmt.Sprintf("  Version: %s\n", csv.Spec.Version.String()))
	result.WriteString(fmt.Sprintf("  Replaces: %s\n", csv.Spec.Replaces))
	result.WriteString(fmt.Sprintf("  Display Name: %s\n", csv.Spec.DisplayName))
	result.WriteString(fmt.Sprintf("  Description: %s\n\n", csv.Spec.Description))

	result.WriteString("Full JSON representation:\n")
	result.WriteString("```json\n")
	result.WriteString(string(jsonData))
	result.WriteString("\n```\n")

	return &types.MCPResponse{
		Content: []types.MCPContent{{
			Type: "text",
			Text: result.String(),
		}},
	}, nil
}