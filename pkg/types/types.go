package types

import (
	"context"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type MCPServer struct {
	Config     *rest.Config
	K8sClient  kubernetes.Interface
	OLMClient  OLMClientInterface
	Port       int
	ReadOnly   bool
	Kubeconfig string
	Toolsets   []string
}

type OLMClientInterface interface {
	ListClusterServiceVersions(ctx context.Context, namespace string) (*v1alpha1.ClusterServiceVersionList, error)
	GetClusterServiceVersion(ctx context.Context, namespace, name string) (*v1alpha1.ClusterServiceVersion, error)
	ListSubscriptions(ctx context.Context, namespace string) (*v1alpha1.SubscriptionList, error)
	GetSubscription(ctx context.Context, namespace, name string) (*v1alpha1.Subscription, error)
	ListCatalogSources(ctx context.Context, namespace string) (*v1alpha1.CatalogSourceList, error)
	GetCatalogSource(ctx context.Context, namespace, name string) (*v1alpha1.CatalogSource, error)
	ListInstallPlans(ctx context.Context, namespace string) (*v1alpha1.InstallPlanList, error)
	GetInstallPlan(ctx context.Context, namespace, name string) (*v1alpha1.InstallPlan, error)
}

// MCP Protocol types
type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id,omitempty"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool definitions for MCP
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ToolConfig struct {
	Name        string
	Description string
	Enabled     bool
}

var DefaultToolsets = map[string][]ToolConfig{
	"csv": {
		{Name: "list_csvs", Description: "List ClusterServiceVersions", Enabled: true},
		{Name: "get_csv", Description: "Get ClusterServiceVersion details", Enabled: true},
	},
	"subscription": {
		{Name: "list_subscriptions", Description: "List Subscriptions", Enabled: true},
		{Name: "get_subscription", Description: "Get Subscription details", Enabled: true},
	},
	"catalog": {
		{Name: "list_catalog_sources", Description: "List CatalogSources", Enabled: true},
		{Name: "get_catalog_source", Description: "Get CatalogSource details", Enabled: true},
	},
	"installplan": {
		{Name: "list_install_plans", Description: "List InstallPlans", Enabled: true},
		{Name: "get_install_plan", Description: "Get InstallPlan details", Enabled: true},
	},
}
