/*
Copyright 2024.

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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	scalingv1 "example.com/pod-scaler/api/v1"
)

// PodScalerReconciler reconciles a PodScaler object
type PodScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=scaling.example.com,resources=podscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scaling.example.com,resources=podscalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=scaling.example.com,resources=podscalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *PodScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// ログをPodScalerに紐づける
	log := r.Log.WithValues("podscaler", req.NamespacedName)

	// PodScaler型
	var podScaler scalingv1.PodScaler
	// 指定した名前空間と名前に基づいてPodScalerを取得
	if err := r.Get(ctx, req.NamespacedName, &podScaler); err != nil {
		log.Error(err, "unable to fetch PodScaler")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 対象のPodをリスト
	var pods corev1.PodList
	labelSelector := labels.SelectorFromSet(podScaler.Spec.Selector)
	if err := r.List(ctx, &pods, &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: labelSelector,
	}); err != nil {
		log.Error(err, "unable to list pods")
		return ctrl.Result{}, err
	}

	// Podの数を調整
	currentCount := len(pods.Items)
	desiredCount := podScaler.Spec.Count

	if currentCount < desiredCount {
		for i := 0; i < (desiredCount - currentCount); i++ {
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "scaled-pod-",
					Namespace:    req.Namespace,
					Labels:       podScaler.Spec.Selector,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			}
			if err := r.Create(ctx, pod); err != nil {
				log.Error(err, "unable to create Pod")
				return ctrl.Result{}, err
			}
		}
	} else if currentCount > desiredCount {
		// Podを削除
		for i := 0; i < (currentCount - desiredCount); i++ {
			pod := &pods.Items[i]
			if err := r.Delete(ctx, pod); err != nil {
				log.Error(err, "unable to delete Pod")
				return ctrl.Result{}, err
			}
		}
	}
	log.Info("Reconciliation complete", "currentCount", currentCount, "desiredCount", desiredCount)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalingv1.PodScaler{}).
		Named("podscaler").
		Complete(r)
}
