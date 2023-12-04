package lock

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/coordination/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Locker interface {
	Lock() error
	UnLock() error
}

type LeaseLock struct {
	Namespace    string
	Name         string
	Client       *kubernetes.Clientset
	LeaseSeconds int32
	Identity     string
}

func (l *LeaseLock) getOrCreateLease(ctx context.Context) (*v1.Lease, error) {
	cli := l.Client.CoordinationV1().Leases(l.Namespace)
	lease, err := cli.Get(ctx, l.Name, metav1.GetOptions{})
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}
	if apierrors.IsNotFound(err) {
		lease, err = cli.Create(ctx, l.newLease(), metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			return cli.Get(ctx, l.Name, metav1.GetOptions{})
		}
		return lease, err
	}
	return lease, nil
}

func (l *LeaseLock) newLease() *v1.Lease {
	now := metav1.NowMicro()
	trans := int32(0)
	return &v1.Lease{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: l.Namespace,
			Name:      l.Name,
		},
		Spec: v1.LeaseSpec{
			HolderIdentity:       &l.Identity,
			LeaseDurationSeconds: &l.LeaseSeconds,
			AcquireTime:          &now,
			RenewTime:            &now,
			LeaseTransitions:     &trans,
		},
	}
}

func (l *LeaseLock) getLease(ctx context.Context) (*v1.Lease, error) {
	return l.Client.CoordinationV1().Leases(l.Namespace).Get(ctx, l.Name, metav1.GetOptions{})
}

func (l *LeaseLock) updateLease(ctx context.Context, lease *v1.Lease) (*v1.Lease, error) {
	return l.Client.CoordinationV1().Leases(l.Namespace).Update(ctx, lease, metav1.UpdateOptions{})
}

func (l *LeaseLock) Lock(ctx context.Context, timeout int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	lease, err := l.getOrCreateLease(ctx)
	if err != nil {
		return err
	}

	if l.AcquireLease(ctx, lease) {
		return nil
	}

	watchHdlr, err := l.Client.CoordinationV1().Leases(l.Namespace).Watch(
		ctx, metav1.ListOptions{Watch: true, FieldSelector: fmt.Sprintf("metadata.name=%s", l.Name)})
	defer watchHdlr.Stop()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			fmt.Println("ctx done with error: ", err)
			return err
		case e := <-watchHdlr.ResultChan():
			fmt.Println("event", e)
			eLease := e.Object.(*v1.Lease)
			fmt.Println("event lease", eLease.Name, eLease.Spec.RenewTime, eLease.Spec.LeaseDurationSeconds)
			if (e.Type == watch.Added || e.Type == watch.Modified) && l.AcquireLease(ctx, eLease) {
				return nil
			}
		}
	}
}

func (l *LeaseLock) LeaseExpire(lease *v1.Lease) bool {
	if lease.Spec.HolderIdentity == nil {
		return true
	}
	now := metav1.NowMicro().UTC()
	duration := time.Duration(*lease.Spec.LeaseDurationSeconds)
	return lease.Spec.RenewTime.UTC().Add(duration * time.Second).Before(now)
}

func (l *LeaseLock) AcquireLease(ctx context.Context, lease *v1.Lease) bool {
	fmt.Println("acquire lease: ", lease.Name, l.Identity)
	now := metav1.NowMicro()
	if !l.LeaseExpire(lease) && *lease.Spec.HolderIdentity != l.Identity {
		fmt.Println("lease not expire and current holder not id required", lease.Name, l.Identity, lease.Spec)
		return false
	}

	if lease.Spec.HolderIdentity != nil && *lease.Spec.HolderIdentity == l.Identity {
		lease.Spec.RenewTime = &now
		lease.Spec.LeaseDurationSeconds = &l.LeaseSeconds
	} else {
		lease.Spec.HolderIdentity = &l.Identity
		lease.Spec.AcquireTime = &now
		lease.Spec.RenewTime = &now
		var trans int32
		if lease.Spec.LeaseTransitions != nil {
			trans = *lease.Spec.LeaseTransitions + 1
		}
		lease.Spec.LeaseTransitions = &trans
		lease.Spec.LeaseDurationSeconds = &l.LeaseSeconds
	}

	lease, err := l.updateLease(ctx, lease)
	fmt.Println("acquire lease: ", lease.Name, l.Identity, err)
	return err == nil
}

func (l *LeaseLock) UnLock() error {
	ctx := context.Background()
	lease, err := l.getLease(ctx)
	if err != nil {
		return client.IgnoreNotFound(err)
	}
	if *lease.Spec.HolderIdentity != l.Identity {
		return nil
	}

	lease.Spec.HolderIdentity = nil
	lease.Spec.LeaseDurationSeconds = nil
	lease.Spec.AcquireTime = nil
	lease.Spec.RenewTime = nil
	_, err = l.updateLease(ctx, lease)
	return err
}

func (l *LeaseLock) Clear(ctx context.Context) error {
	return nil
}
