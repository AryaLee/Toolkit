package common

import (
	"flag"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/Users/meiyunli/.kube/config", "absolute path to the kubeconfig file")
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewClient() *clientset.Clientset {
	config, err := buildConfig(kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	return clientset.NewForConfigOrDie(config)
}
