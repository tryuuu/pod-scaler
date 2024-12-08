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

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	scalingv1 "example.com/pod-scaler/api/v1"
)

// nolint:unused
// log is for logging in this package.
var podscalerlog = logf.Log.WithName("podscaler-resource")

// SetupPodScalerWebhookWithManager registers the webhook for PodScaler in the manager.
func SetupPodScalerWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&scalingv1.PodScaler{}).
		WithValidator(&PodScalerCustomValidator{}).
		WithDefaulter(&PodScalerCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-scaling-example-com-v1-podscaler,mutating=true,failurePolicy=fail,sideEffects=None,groups=scaling.example.com,resources=podscalers,verbs=create;update,versions=v1,name=mpodscaler-v1.kb.io,admissionReviewVersions=v1

// PodScalerCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind PodScaler when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PodScalerCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &PodScalerCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind PodScaler.
func (d *PodScalerCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	podscaler, ok := obj.(*scalingv1.PodScaler)
	if !ok {
		return fmt.Errorf("expected an PodScaler object but got %T", obj)
	}

	// Mutating Admission Webhookのデフォルト値を設定(spec.count)
	if podscaler.Spec.Count < 1 {
		podscaler.Spec.Count = 5
	}

	podscalerlog.Info("Defaulting for PodScaler", "name", podscaler.GetName())
	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-scaling-example-com-v1-podscaler,mutating=false,failurePolicy=fail,sideEffects=None,groups=scaling.example.com,resources=podscalers,verbs=create;update,versions=v1,name=vpodscaler-v1.kb.io,admissionReviewVersions=v1

// PodScalerCustomValidator struct is responsible for validating the PodScaler resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type PodScalerCustomValidator struct {
	//TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &PodScalerCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type PodScaler.
func (v *PodScalerCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	podscaler, ok := obj.(*scalingv1.PodScaler)
	if !ok {
		return nil, fmt.Errorf("expected a PodScaler object but got %T", obj)
	}

	// validating admission webhookの設定
	if podscaler.Spec.Count < 1 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	if len(podscaler.Spec.Selector) == 0 {
		return nil, fmt.Errorf("selector must be specified")
	}

	podscalerlog.Info("Validation for PodScaler upon creation", "name", podscaler.GetName())

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type PodScaler.
func (v *PodScalerCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	podscaler, ok := newObj.(*scalingv1.PodScaler)

	// validating admission webhookの設定
	if podscaler.Spec.Count < 1 {
		return nil, fmt.Errorf("count must be greater than 0")
	}

	if len(podscaler.Spec.Selector) == 0 {
		return nil, fmt.Errorf("selector must be specified")
	}

	if !ok {
		return nil, fmt.Errorf("expected a PodScaler object for the newObj but got %T", newObj)
	}
	podscalerlog.Info("Validation for PodScaler upon update", "name", podscaler.GetName())

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type PodScaler.
func (v *PodScalerCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	podscaler, ok := obj.(*scalingv1.PodScaler)
	if !ok {
		return nil, fmt.Errorf("expected a PodScaler object but got %T", obj)
	}
	podscalerlog.Info("Validation for PodScaler upon deletion", "name", podscaler.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
