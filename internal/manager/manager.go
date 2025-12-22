package manager

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Taiwrash/trigra/internal/k8s"
	"github.com/Taiwrash/trigra/internal/providers"
	"github.com/Taiwrash/trigra/internal/providers/bitbucket"
	"github.com/Taiwrash/trigra/internal/providers/git"
	"github.com/Taiwrash/trigra/internal/providers/gitea"
	"github.com/Taiwrash/trigra/internal/providers/github"
	"github.com/Taiwrash/trigra/internal/providers/gitlab"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Project represents a configured sync project
type Project struct {
	Name            string
	URL             string
	Provider        providers.Provider
	WebhookSecret   string
	TargetNamespace string
	Branch          string
}

// Manager manages multiple GitOps projects
type Manager struct {
	projects  map[string]*Project
	mu        sync.RWMutex
	applier   *k8s.Applier
	clientset *kubernetes.Clientset
	dynClient dynamic.Interface
}

// NewManager creates a new project manager
func NewManager(applier *k8s.Applier, clientset *kubernetes.Clientset, dynClient dynamic.Interface) *Manager {
	return &Manager{
		projects:  make(map[string]*Project),
		applier:   applier,
		clientset: clientset,
		dynClient: dynClient,
	}
}

// Start watches for GitRepo resources
func (m *Manager) Start(ctx context.Context) error {
	log.Println("Starting GitRepo controller...")

	gvr := schema.GroupVersionResource{
		Group:    "trigra.io",
		Version:  "v1alpha1",
		Resource: "gitrepos",
	}

	factory := dynamicinformer.NewDynamicSharedInformerFactory(m.dynClient, 0)
	informer := factory.ForResource(gvr).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			m.handleUpdate(obj)
		},
		UpdateFunc: func(_, newObj interface{}) {
			m.handleUpdate(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			m.handleDelete(obj)
		},
	})

	go informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		return fmt.Errorf("failed to sync cache")
	}

	log.Println("GitRepo controller synced")
	return nil
}

func (m *Manager) handleUpdate(obj interface{}) {
	un := obj.(*unstructured.Unstructured)
	name := un.GetName()
	ns := un.GetNamespace()

	spec, _, _ := unstructured.NestedMap(un.Object, "spec")
	url, _ := spec["url"].(string)
	providerName, _ := spec["provider"].(string)
	targetNS, _ := spec["targetNamespace"].(string)
	branch, _ := spec["branch"].(string)

	if targetNS == "" {
		targetNS = ns
	}
	if branch == "" {
		branch = "main"
	}

	// Load secrets
	token := m.loadSecret(ns, spec["tokenSecretRef"])
	webhookSecret := m.loadSecret(ns, spec["webhookSecretRef"])

	var provider providers.Provider
	switch providerName {
	case "github":
		provider = github.NewProvider(token)
	case "gitlab":
		baseURL, _ := spec["baseURL"].(string)
		provider = gitlab.NewProvider(baseURL, token)
	case "gitea":
		baseURL, _ := spec["baseURL"].(string)
		provider = gitea.NewProvider(baseURL, token)
	case "bitbucket":
		user, _ := spec["bitbucketUser"].(string)
		provider = bitbucket.NewProvider(user, token)
	case "git":
		keyFile, _ := spec["sshKeyFile"].(string)
		provider = git.NewProvider(url, keyFile)
	}

	project := &Project{
		Name:            name,
		URL:             url,
		Provider:        provider,
		WebhookSecret:   webhookSecret,
		TargetNamespace: targetNS,
		Branch:          branch,
	}

	m.mu.Lock()
	m.projects[name] = project
	m.mu.Unlock()

	log.Printf("INFO: Configured project: %s (Provider: %s, URL: %s)", name, providerName, url)
}

func (m *Manager) handleDelete(obj interface{}) {
	un, ok := obj.(*unstructured.Unstructured)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return
		}
		un = tombstone.Obj.(*unstructured.Unstructured)
	}

	m.mu.Lock()
	delete(m.projects, un.GetName())
	m.mu.Unlock()

	log.Printf("INFO: Deleted project: %s", un.GetName())
}

func (m *Manager) loadSecret(ns string, ref interface{}) string {
	if ref == nil {
		return ""
	}
	mref, ok := ref.(map[string]interface{})
	if !ok {
		return ""
	}

	name, _ := mref["name"].(string)
	key, _ := mref["key"].(string)

	if name == "" || key == "" {
		return ""
	}

	secret, err := m.clientset.CoreV1().Secrets(ns).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("WARNING: Failed to load secret %s/%s: %v", ns, name, err)
		return ""
	}

	return string(secret.Data[key])
}

// GetProject returns a project by name
func (m *Manager) GetProject(name string) *Project {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.projects[name]
}

// AddStaticProject adds a project manually (legacy/env support)
func (m *Manager) AddStaticProject(p *Project) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.projects["default"] = p // Use "default" for env-based project
}
