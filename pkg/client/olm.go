package client

import (
	"context"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OLMClient struct {
	config *rest.Config
	client rest.Interface
}

func NewOLMClient(config *rest.Config) (*OLMClient, error) {
	v1alpha1.AddToScheme(scheme.Scheme)

	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	client, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}

	return &OLMClient{
		config: config,
		client: client,
	}, nil
}

func (c *OLMClient) ListClusterServiceVersions(ctx context.Context, namespace string) (*v1alpha1.ClusterServiceVersionList, error) {
	result := &v1alpha1.ClusterServiceVersionList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("clusterserviceversions").
		Do(ctx).
		Into(result)
	return result, err
}

func (c *OLMClient) GetClusterServiceVersion(ctx context.Context, namespace, name string) (*v1alpha1.ClusterServiceVersion, error) {
	result := &v1alpha1.ClusterServiceVersion{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("clusterserviceversions").
		Name(name).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *OLMClient) ListSubscriptions(ctx context.Context, namespace string) (*v1alpha1.SubscriptionList, error) {
	result := &v1alpha1.SubscriptionList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("subscriptions").
		Do(ctx).
		Into(result)
	return result, err
}

func (c *OLMClient) GetSubscription(ctx context.Context, namespace, name string) (*v1alpha1.Subscription, error) {
	result := &v1alpha1.Subscription{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("subscriptions").
		Name(name).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *OLMClient) ListCatalogSources(ctx context.Context, namespace string) (*v1alpha1.CatalogSourceList, error) {
	result := &v1alpha1.CatalogSourceList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("catalogsources").
		Do(ctx).
		Into(result)
	return result, err
}

func (c *OLMClient) GetCatalogSource(ctx context.Context, namespace, name string) (*v1alpha1.CatalogSource, error) {
	result := &v1alpha1.CatalogSource{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("catalogsources").
		Name(name).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *OLMClient) ListInstallPlans(ctx context.Context, namespace string) (*v1alpha1.InstallPlanList, error) {
	result := &v1alpha1.InstallPlanList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("installplans").
		Do(ctx).
		Into(result)
	return result, err
}

func (c *OLMClient) GetInstallPlan(ctx context.Context, namespace, name string) (*v1alpha1.InstallPlan, error) {
	result := &v1alpha1.InstallPlan{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("installplans").
		Name(name).
		Do(ctx).
		Into(result)
	return result, err
}
