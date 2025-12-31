package k8s

import (
	"fmt"
	"path/filepath"
	"sync"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	clientsetOnce sync.Once
	clientset     *kubernetes.Clientset
	clientsetErr  error

	dynamicOnce   sync.Once
	dynamicClient dynamic.Interface
	dynamicErr    error
)

// GetClientset returns a singleton Kubernetes clientset
// It automatically detects in-cluster vs out-of-cluster configuration
func GetClientset(inCluster bool) (*kubernetes.Clientset, error) {
	clientsetOnce.Do(func() {
		config, err := getConfig(inCluster)
		if err != nil {
			clientsetErr = fmt.Errorf("failed to get kubernetes config: %w", err)
			return
		}

		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			clientsetErr = fmt.Errorf("failed to create kubernetes clientset: %w", err)
			return
		}
	})

	return clientset, clientsetErr
}

// GetDynamicClient returns a singleton dynamic Kubernetes client
// This is used for applying arbitrary resource types
func GetDynamicClient(inCluster bool) (dynamic.Interface, error) {
	dynamicOnce.Do(func() {
		config, err := getConfig(inCluster)
		if err != nil {
			dynamicErr = fmt.Errorf("failed to get kubernetes config: %w", err)
			return
		}

		dynamicClient, err = dynamic.NewForConfig(config)
		if err != nil {
			dynamicErr = fmt.Errorf("failed to create dynamic client: %w", err)
			return
		}
	})

	return dynamicClient, dynamicErr
}

// getConfig returns the appropriate Kubernetes configuration
func getConfig(inCluster bool) (*rest.Config, error) {
	if inCluster {
		// Use in-cluster config (running inside a pod)
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
		}
		return config, nil
	}

	// Use kubeconfig file (local development)
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Better yet, use the environment variable if present
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	return config, nil
}
