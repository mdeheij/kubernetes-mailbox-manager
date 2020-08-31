package v1

import (
	v1 "github.com/mdeheij/kubernetes-mailbox-manager/api/types/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	core_v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type V1Interface interface {
	Mailboxes(namespace string) MailboxInterface
	ConfigMaps(namespace string) core_v1.ConfigMapInterface
}

type V1Client struct {
	baseClient *kubernetes.Clientset
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*V1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		return nil, err
	}

	return &V1Client{restClient: client, baseClient: clientset}, nil
}

func (c *V1Client) Mailboxes(namespace string) MailboxInterface {
	return &mailboxClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *V1Client) ConfigMaps(namespace string) core_v1.ConfigMapInterface {
	return c.baseClient.CoreV1().ConfigMaps(namespace)
}
