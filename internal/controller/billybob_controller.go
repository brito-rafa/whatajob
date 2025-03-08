/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"math/rand"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	whatajobiov1 "github.com/gramLabs/whatajob/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BillybobReconciler reconciles a Billybob object
type BillybobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=whatajob.io,resources=billybobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=whatajob.io,resources=billybobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=whatajob.io,resources=billybobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Billybob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile

func (r *BillybobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	billybob := &whatajobiov1.Billybob{}
	if err := r.Get(ctx, req.NamespacedName, billybob); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// List Pods owned by Billybob
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	ownedPods := make([]corev1.Pod, 0)
	for _, pod := range podList.Items {
		for _, ownerRef := range pod.OwnerReferences {
			if ownerRef.UID == billybob.UID && ownerRef.Controller != nil && *ownerRef.Controller {
				ownedPod := pod
				ownedPod.TypeMeta = metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}
				ownedPod.ObjectMeta = pod.ObjectMeta
				ownedPods = append(ownedPods, ownedPod)
				break
			}
		}
	}

	// Now `ownedPod` contains all Pods owned by the Billybob object

	currentPods := len(ownedPods)
	desiredPods := billybob.Spec.NumberPods

	// need to create logic to remove pods

	for i := currentPods; i < desiredPods; i++ {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      generateRandomPodName(),
				Namespace: req.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(billybob, whatajobiov1.GroupVersion.WithKind("Billybob")),
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "busybox",
						Image: "busybox",
						Args:  []string{"sleep", "3600"},
					},
				},
			},
		}

		if err := r.Create(ctx, pod); err != nil {
			continue
		}
	}

	// List Pods owned by Billybob
	newPodList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	currentOwnedPodsNames := make([]string, 0)
	for _, pod := range newPodList.Items {
		for _, ownerRef := range pod.OwnerReferences {
			if ownerRef.UID == billybob.UID && ownerRef.Controller != nil && *ownerRef.Controller {
				ownedPod := pod
				ownedPod.TypeMeta = metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}
				ownedPod.ObjectMeta = pod.ObjectMeta
				currentOwnedPodsNames = append(currentOwnedPodsNames, ownedPod.Name)
				break
			}
		}
	}

	billybob.Status.CurrentPods = len(currentOwnedPodsNames)
	billybob.Status.PodNames = currentOwnedPodsNames
	if err := r.Status().Update(ctx, billybob); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func generateRandomPodName() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return "billypod-" + string(b)
}

func (r *BillybobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&whatajobiov1.Billybob{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
