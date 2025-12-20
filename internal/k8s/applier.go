// Package k8s provides mechanisms for applying Kubernetes resources using dynamic clients.
package k8s

import (
	"context"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// Applier handles applying Kubernetes resources
type Applier struct {
	dynamicClient dynamic.Interface
	mapper        meta.RESTMapper
}

// NewApplier creates a new resource applier
func NewApplier(inCluster bool) (*Applier, error) {
	dynamicClient, err := GetDynamicClient(inCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get dynamic client: %w", err)
	}

	// Get config for discovery client
	config, err := getConfig(inCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Create a RESTMapper to discover resource types
	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(config)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	return &Applier{
		dynamicClient: dynamicClient,
		mapper:        mapper,
	}, nil
}

// ApplyResource applies a Kubernetes resource from YAML content
// It supports all Kubernetes resource types (Deployments, Services, ConfigMaps, Secrets, etc.)
func (a *Applier) ApplyResource(ctx context.Context, yamlContent []byte, namespace string) error {
	// Decode YAML to unstructured object
	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode(yamlContent, nil, obj)
	if err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}

	// Find the GVR (Group/Version/Resource) for this resource type
	mapping, err := a.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("failed to get REST mapping for %s: %w", gvk.Kind, err)
	}

	// Get the appropriate dynamic client for this resource
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// Use namespace from object if set, otherwise use provided namespace
		ns := obj.GetNamespace()
		if ns == "" {
			ns = namespace
			obj.SetNamespace(ns)
		}
		dr = a.dynamicClient.Resource(mapping.Resource).Namespace(ns)
	} else {
		// Cluster-scoped resource
		dr = a.dynamicClient.Resource(mapping.Resource)
	}

	// Try to get the existing resource
	existing, err := dr.Get(ctx, obj.GetName(), metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Resource doesn't exist, create it
			log.Printf("Creating %s/%s in namespace %s", gvk.Kind, obj.GetName(), obj.GetNamespace())
			_, err = dr.Create(ctx, obj, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create %s/%s: %w", gvk.Kind, obj.GetName(), err)
			}
			log.Printf("Successfully created %s/%s", gvk.Kind, obj.GetName())
			return nil
		}
		return fmt.Errorf("failed to get existing resource: %w", err)
	}

	// Resource exists, update it
	// Preserve resourceVersion for optimistic concurrency
	obj.SetResourceVersion(existing.GetResourceVersion())

	log.Printf("Updating %s/%s in namespace %s", gvk.Kind, obj.GetName(), obj.GetNamespace())
	_, err = dr.Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update %s/%s: %w", gvk.Kind, obj.GetName(), err)
	}

	log.Printf("Successfully updated %s/%s", gvk.Kind, obj.GetName())
	return nil
}

// ApplyMultipleResources applies multiple YAML documents from a single file
// Documents should be separated by "---"
func (a *Applier) ApplyMultipleResources(ctx context.Context, yamlContent []byte, namespace string) error {
	// Split YAML documents
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	// Simple split by --- (this is a basic implementation)
	// For production, consider using a proper YAML stream parser
	docs := splitYAMLDocuments(yamlContent)

	for i, doc := range docs {
		if len(doc) == 0 {
			continue
		}

		obj := &unstructured.Unstructured{}
		_, _, err := decoder.Decode(doc, nil, obj)
		if err != nil {
			log.Printf("Warning: failed to decode document %d: %v", i, err)
			continue
		}

		if err := a.ApplyResource(ctx, doc, namespace); err != nil {
			return fmt.Errorf("failed to apply document %d: %w", i, err)
		}
	}

	return nil
}

// splitYAMLDocuments splits a YAML file into individual documents
func splitYAMLDocuments(content []byte) [][]byte {
	// This is a simple implementation
	// For production use, consider a more robust YAML parser
	separator := []byte("\n---\n")
	docs := [][]byte{}

	start := 0
	for {
		idx := indexOf(content[start:], separator)
		if idx == -1 {
			// Last document
			if start < len(content) {
				docs = append(docs, content[start:])
			}
			break
		}

		docs = append(docs, content[start:start+idx])
		start += idx + len(separator)
	}

	return docs
}

// indexOf finds the index of a byte slice within another
func indexOf(haystack, needle []byte) int {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		found := true
		for j := 0; j < len(needle); j++ {
			if haystack[i+j] != needle[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}
