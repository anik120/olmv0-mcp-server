package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/client"
	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/server"
	"github.com/operator-framework/operator-lifecycle-manager/olmv0-mcp-server/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	port       int
	kubeconfig string
	readOnly   bool
	toolsets   []string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "olmv0-mcp-server",
		Short: "OLM v0 MCP Server - Model Context Protocol server for Operator Lifecycle Manager",
		Long: `OLM v0 MCP Server provides a Model Context Protocol interface to interact with
Operator Lifecycle Manager resources in a Kubernetes cluster.

This server exposes tools to list, get, and inspect OLM resources including:
- ClusterServiceVersions (CSVs)
- Subscriptions
- CatalogSources
- InstallPlans
- OperatorGroups

It follows the same pattern as the kubernetes-mcp-server but focuses specifically
on OLM v0 resources and operations.`,
		Run: runServer,
	}

	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "HTTP/SSE server port")
	rootCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file (default: $HOME/.kube/config)")
	rootCmd.Flags().BoolVar(&readOnly, "read-only", true, "Prevent write operations (default: true)")
	rootCmd.Flags().StringSliceVar(&toolsets, "toolsets", []string{"csv", "subscription", "catalog", "installplan"}, "Enable specific toolsets")

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	logrus.Info("Starting OLM v0 MCP Server")

	config, err := getKubeConfig(kubeconfig)
	if err != nil {
		logrus.Fatalf("Error getting kubeconfig: %v", err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Error creating Kubernetes client: %v", err)
	}

	olmClient, err := client.NewOLMClient(config)
	if err != nil {
		logrus.Fatalf("Error creating OLM client: %v", err)
	}

	mcpServer := &types.MCPServer{
		Config:     config,
		K8sClient:  k8sClient,
		OLMClient:  olmClient,
		Port:       port,
		ReadOnly:   readOnly,
		Kubeconfig: kubeconfig,
		Toolsets:   toolsets,
	}

	if err := server.StartServer(mcpServer); err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
}

func getKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		if config, err := rest.InClusterConfig(); err == nil {
			logrus.Info("Using in-cluster configuration")
			return config, nil
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting user home directory: %v", err)
		}
		kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
	}

	logrus.Infof("Using kubeconfig: %s", kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error building config from kubeconfig: %v", err)
	}

	return config, nil
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
}