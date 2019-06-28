/*
Copyright 2019 The Kubernetes Authors.

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

package webhook

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/flect"

	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/util"
	"sigs.k8s.io/kubebuilder/pkg/scaffold/v1/resource"
)

// Webhook scaffolds a Webhook for a Resource
type Webhook struct {
	input.Input

	// Resource is the Resource to make the Webhook for
	Resource *resource.Resource

	// ResourcePackage is the package of the Resource
	ResourcePackage string

	// Plural is the plural lowercase of kind
	Plural string

	// Is the Group + "." + Domain for the Resource
	GroupDomain string

	// Is the Group + "." + Domain for the Resource
	GroupDomainWithDash string

	// If scaffold the defaulting webhook
	Defaulting bool
	// If scaffold the validating webhook
	Validating bool
}

// GetInput implements input.File
func (a *Webhook) GetInput() (input.Input, error) {

	a.ResourcePackage, a.GroupDomain = util.GetResourceInfo(a.Resource, a.Input)

	a.GroupDomainWithDash = strings.Replace(a.GroupDomain, ".", "-", -1)

	if a.Plural == "" {
		a.Plural = flect.Pluralize(strings.ToLower(a.Resource.Kind))
	}

	if a.Path == "" {
		a.Path = filepath.Join("api", a.Resource.Version,
			fmt.Sprintf("%s_webhook.go", strings.ToLower(a.Resource.Kind)))
	}
	if a.Defaulting {
		WebhookTemplate = WebhookTemplate + DefaultingWebhookTemplate
	}
	if a.Validating {
		WebhookTemplate = WebhookTemplate + ValidatingWebhookTemplate
	}

	a.TemplateBody = WebhookTemplate
	a.Input.IfExistsAction = input.Error
	return a.Input, nil
}

var (
	WebhookTemplate = `{{ .Boilerplate }}

package {{ .Resource.Version }}

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var {{ lower .Resource.Kind }}log = logf.Log.WithName("{{ lower .Resource.Kind }}-resource")

func (r *{{.Resource.Kind}}) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r)
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
`

	DefaultingWebhookTemplate = `
// +kubebuilder:webhook:path=/mutate-{{ .GroupDomainWithDash }}-{{ .Resource.Version }}-{{ lower .Resource.Kind }},mutating=true,failurePolicy=fail,groups={{ .Resource.Group }},resources={{ .Plural }},verbs=create;update,versions=v1,name=m{{ lower .Resource.Kind }}.kb.io

var _ webhook.Defaulter = &{{ .Resource.Kind }}{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *{{ .Resource.Kind }}) Default() {
	{{ lower .Resource.Kind }}log.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}
`

	ValidatingWebhookTemplate = `
// +kubebuilder:webhook:path=/validate-{{ .GroupDomainWithDash }}-{{ .Resource.Version }}-{{ lower .Resource.Kind }},mutating=false,failurePolicy=fail,groups={{ .Resource.Group }},resources={{ .Plural }},verbs=create;update,versions=v1,name=v{{ lower .Resource.Kind }}.kb.io

var _ webhook.Validator = &{{ .Resource.Kind }}{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *{{ .Resource.Kind }}) ValidateCreate() error {
	{{ lower .Resource.Kind }}log.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *{{ .Resource.Kind }}) ValidateUpdate(old runtime.Object) error {
	{{ lower .Resource.Kind }}log.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}
`
)
