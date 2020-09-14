package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"

	types_v1 "github.com/mdeheij/kubernetes-mailbox-manager/api/types/v1"
	client_v1 "github.com/mdeheij/kubernetes-mailbox-manager/clientset/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type Transpiler struct {
	clientSet  client_v1.V1Interface
	lastOutput string
	store      cache.Store
	controller cache.Controller
}

func NewTranspiler(clientSet client_v1.V1Interface) *Transpiler {
	t := &Transpiler{
		clientSet: clientSet,
	}

	t.store, t.controller = cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return t.clientSet.Mailboxes(v1.NamespaceAll).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return t.clientSet.Mailboxes(v1.NamespaceAll).Watch(lo)
			},
		},
		&types_v1.Mailbox{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    t.OnAdd,
			DeleteFunc: t.OnDelete,
			UpdateFunc: t.OnUpdate,
		},
	)

	return t
}

func (t *Transpiler) Run(stopCh <-chan struct{}) {
	t.controller.Run(stopCh)
}

func (t *Transpiler) rebuild() {
	var content string
	var lines []string

	l, err := t.clientSet.Mailboxes(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		log.Println("Error listing mailboxes: ", err)
	}

	for _, m := range l.Items {
		lines = append(lines, fmt.Sprintf("%s|%s", m.Spec.EmailAddress, m.Spec.PasswordHash))
	}

	sort.Strings(lines)

	content = strings.Join(lines, "\n")
	if t.lastOutput != "" && content == t.lastOutput {
		return
	}

	t.lastOutput = content

	x := &v1.ConfigMap{
		Data: map[string]string{
			"postfix-accounts.cf": content,
		},
	}

	x.Name = "kubernetes-mailbox-manager"

	if x, err := t.clientSet.ConfigMaps("default").Get(context.Background(), x.Name, metav1.GetOptions{}); x != nil && err == nil {
		log.Println(t.clientSet.ConfigMaps("default").Create(context.Background(), x, metav1.CreateOptions{}))
	} else if err != nil {
		log.Println(err)
	}

	if _, err := t.clientSet.ConfigMaps("default").Update(context.Background(), x, metav1.UpdateOptions{}); err != nil {
		log.Println(err)
	}
}

func (t *Transpiler) OnAdd(obj interface{}) {
	t.rebuild()
}

func (t *Transpiler) OnUpdate(oldObj, newObj interface{}) {
	t.rebuild()
}

func (t *Transpiler) OnDelete(obj interface{}) {
	t.rebuild()
}
