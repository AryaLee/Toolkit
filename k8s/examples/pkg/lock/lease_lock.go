package lock

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog"
)

type LeaseLock struct {
	Name      string
	Namespace string
	ID        string
	Func      RunFunc
	Client    *clientset.Clientset
}

type RunFunc func(ctx context.Context)

func (l *LeaseLock) Run(ctx context.Context) {
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      l.Name,
			Namespace: l.Namespace,
		},
		Client: l.Client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: l.ID,
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   5 * time.Second,
		RetryPeriod:     3 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				klog.Infof("leader started: %s", l.ID)
				l.Func(ctx)
				cancel()
			},
			OnStoppedLeading: func() {
				klog.Infof("leader lost: %s", l.ID)
				cancel()
			},
			OnNewLeader: func(identity string) {
				if identity == l.ID {
					klog.Infof("leader new: %s", l.ID)
					return
				}
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})
}
