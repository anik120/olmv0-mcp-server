package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/types"
)

type InstallPlanTools struct {
	server *types.MCPServer
}

func NewInstallPlanTools(server *types.MCPServer) *InstallPlanTools {
	return &InstallPlanTools{server: server}
}

func (t *InstallPlanTools) ListInstallPlans(ctx context.Context, params map[string]string) (*types.MCPResponse, error) {
	namespace := params["namespace"]
	if namespace == "" {
		namespace = "default"
	}

	installPlans, err := t.server.OLMClient.ListInstallPlans(ctx, namespace)
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error listing InstallPlans: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("InstallPlans in namespace '%s':\n\n", namespace))

	if len(installPlans.Items) == 0 {
		result.WriteString("No InstallPlans found.\n")
	} else {
		result.WriteString("NAME\tNAMESPACE\tAPPROVAL\tAPPROVED\tPHASE\n")
		for _, ip := range installPlans.Items {
			result.WriteString(fmt.Sprintf("%s\t%s\t%s\t%t\t%s\n",
				ip.Name,
				ip.Namespace,
				ip.Spec.Approval,
				ip.Spec.Approved,
				ip.Status.Phase,
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

func (t *InstallPlanTools) GetInstallPlan(ctx context.Context, params map[string]string) (*types.MCPResponse, error) {
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

	installPlan, err := t.server.OLMClient.GetInstallPlan(ctx, namespace, name)
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error getting InstallPlan '%s': %v", name, err),
			}},
			IsError: true,
		}, nil
	}

	jsonData, err := json.MarshalIndent(installPlan, "", "  ")
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error marshaling InstallPlan to JSON: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("InstallPlan: %s/%s\n\n", namespace, name))
	result.WriteString("Basic Info:\n")
	result.WriteString(fmt.Sprintf("  Name: %s\n", installPlan.Name))
	result.WriteString(fmt.Sprintf("  Namespace: %s\n", installPlan.Namespace))
	result.WriteString(fmt.Sprintf("  Approval: %s\n", installPlan.Spec.Approval))
	result.WriteString(fmt.Sprintf("  Approved: %t\n", installPlan.Spec.Approved))
	result.WriteString(fmt.Sprintf("  Phase: %s\n", installPlan.Status.Phase))
	result.WriteString("\n")

	if len(installPlan.Spec.ClusterServiceVersionNames) > 0 {
		result.WriteString("ClusterServiceVersions to install:\n")
		for _, csvName := range installPlan.Spec.ClusterServiceVersionNames {
			result.WriteString(fmt.Sprintf("  - %s\n", csvName))
		}
		result.WriteString("\n")
	}

	if len(installPlan.Status.Plan) > 0 {
		result.WriteString("Planned Resources:\n")
		for _, step := range installPlan.Status.Plan {
			result.WriteString(fmt.Sprintf("  - %s: %s/%s (Status: %s)\n",
				step.Resource.Kind,
				step.Resource.Name,
				step.Resource.Manifest,
				step.Status,
			))
		}
		result.WriteString("\n")
	}

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