package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/types"
)

type SubscriptionTools struct {
	server *types.MCPServer
}

func NewSubscriptionTools(server *types.MCPServer) *SubscriptionTools {
	return &SubscriptionTools{server: server}
}

func (t *SubscriptionTools) ListSubscriptions(ctx context.Context, params map[string]string) (*types.MCPResponse, error) {
	namespace := params["namespace"]
	if namespace == "" {
		namespace = "default"
	}

	subscriptions, err := t.server.OLMClient.ListSubscriptions(ctx, namespace)
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error listing Subscriptions: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Subscriptions in namespace '%s':\n\n", namespace))

	if len(subscriptions.Items) == 0 {
		result.WriteString("No Subscriptions found.\n")
	} else {
		result.WriteString("NAME\tNAMESPACE\tPACKAGE\tCHANNEL\tSOURCE\tINSTALLED CSV\n")
		for _, sub := range subscriptions.Items {
			result.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n",
				sub.Name,
				sub.Namespace,
				sub.Spec.Package,
				sub.Spec.Channel,
				sub.Spec.CatalogSource,
				sub.Status.InstalledCSV,
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

func (t *SubscriptionTools) GetSubscription(ctx context.Context, params map[string]string) (*types.MCPResponse, error) {
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

	subscription, err := t.server.OLMClient.GetSubscription(ctx, namespace, name)
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error getting Subscription '%s': %v", name, err),
			}},
			IsError: true,
		}, nil
	}

	jsonData, err := json.MarshalIndent(subscription, "", "  ")
	if err != nil {
		return &types.MCPResponse{
			Content: []types.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Error marshaling Subscription to JSON: %v", err),
			}},
			IsError: true,
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Subscription: %s/%s\n\n", namespace, name))
	result.WriteString("Basic Info:\n")
	result.WriteString(fmt.Sprintf("  Name: %s\n", subscription.Name))
	result.WriteString(fmt.Sprintf("  Namespace: %s\n", subscription.Namespace))
	result.WriteString(fmt.Sprintf("  Package: %s\n", subscription.Spec.Package))
	result.WriteString(fmt.Sprintf("  Channel: %s\n", subscription.Spec.Channel))
	result.WriteString(fmt.Sprintf("  Source: %s\n", subscription.Spec.CatalogSource))
	result.WriteString(fmt.Sprintf("  Source Namespace: %s\n", subscription.Spec.CatalogSourceNamespace))
	result.WriteString(fmt.Sprintf("  Installed CSV: %s\n", subscription.Status.InstalledCSV))
	result.WriteString(fmt.Sprintf("  Current CSV: %s\n\n", subscription.Status.CurrentCSV))

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