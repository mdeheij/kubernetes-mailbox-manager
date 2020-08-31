package v1

import (
	"context"

	v1 "github.com/mdeheij/kubernetes-mailbox-manager/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type MailboxInterface interface {
	List(opts metav1.ListOptions) (*v1.MailboxList, error)
	Get(name string, options metav1.GetOptions) (*v1.Mailbox, error)
	Create(*v1.Mailbox) (*v1.Mailbox, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type mailboxClient struct {
	restClient rest.Interface
	ns         string
}

func (c *mailboxClient) List(opts metav1.ListOptions) (*v1.MailboxList, error) {
	result := v1.MailboxList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mailboxes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *mailboxClient) Get(name string, opts metav1.GetOptions) (*v1.Mailbox, error) {
	result := v1.Mailbox{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mailboxes").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *mailboxClient) Create(mailbox *v1.Mailbox) (*v1.Mailbox, error) {
	result := v1.Mailbox{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("mailboxes").
		Body(mailbox).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *mailboxClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("mailboxes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
