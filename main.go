package main

import (
	"flag"
	"log"

	v1 "github.com/mdeheij/kubernetes-mailbox-manager/api/types/v1"
	clientV1 "github.com/mdeheij/kubernetes-mailbox-manager/clientset/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	if err := v1.AddToScheme(scheme.Scheme); err != nil {
		log.Fatalln("Unable to add to scheme: ", err)
	}

	clientSet, err := clientV1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	NewTranspiler(clientSet).Run(stop)
}
