package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	_ "github.com/onsi/ginkgo/v2"
	"sigs.k8s.io/controller-runtime/pkg/recorder"
)

var (
	recorderProvider recorder.Provider
	somePod          *corev1.Pod
)

func main() {
	// recorderProvider is a recorder.Provider
	recorder := recorderProvider.GetEventRecorderFor("my-controller")

	// emit an event with a variable message
	mildCheese := "Wensleydale"
	recorder.Eventf(somePod, corev1.EventTypeNormal,
		"DislikesCheese", "Not even %s?", mildCheese)
}

func eventT() {
	c := ctrl.GetConfigOrDie()
	cli := clientset.NewForConfigOrDie(c)
	events_cli := cli.CoreV1().Events("default")
	ctx := context.Background()
	ev := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "lmy-test",
		},
		Reason:  "test",
		Message: "lmy-create-test",
		Type:    "Normal",
	}
	e, err := events_cli.Create(ctx, ev, metav1.CreateOptions{})
	fmt.Println(e, err)
}
