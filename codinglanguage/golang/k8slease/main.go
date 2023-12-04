package main

import (
	"context"
	"fmt"
	"time"

	"example.com/aryaLee/golang/k8slease/lock"
	clientset "k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	c := ctrl.GetConfigOrDie()
	cli := clientset.NewForConfigOrDie(c)
	llock := lock.LeaseLock{
		Namespace:    "default",
		Name:         "vpc-sample-subnet",
		Client:       cli,
		LeaseSeconds: 30,
		Identity:     "lmy-test",
	}
	fmt.Println("vim-go", llock)

	ctx := context.Background()
	err := llock.Lock(ctx, 90)
	fmt.Println("lock err: ", llock, err)

	time.AfterFunc(10*time.Second, func() {
		err = llock.UnLock()
		fmt.Println("unlock err: ", llock, err)
	})

	llock2 := lock.LeaseLock{
		Namespace:    "default",
		Name:         "vpc-sample-subnet",
		Client:       cli,
		LeaseSeconds: 30,
		Identity:     "lmy-test1",
	}
	ctx = context.Background()
	err = llock2.Lock(ctx, 90)
	fmt.Println("lock err: ", llock2, err)

	// t := metav1.NowMicro()
	// fmt.Println("time since renew to now:", t.Sub(lease.Spec.RenewTime.Time))
	// // id := "lmy-test"
	// // lease.Spec.HolderIdentity = &id
	// // lease.Spec.AcquireTime = &t
	// // lease.Spec.RenewTime = &t
	// // leasen, err := cli.Update(ctx, lease, metav1.UpdateOptions{})
	// // fmt.Println(err, leasen)
	// mu := sync.Mutex{}
	// mu.TryLock()
}
