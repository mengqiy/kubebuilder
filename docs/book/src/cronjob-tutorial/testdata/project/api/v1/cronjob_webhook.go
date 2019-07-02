/*

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
// +kubebuilder:gen-docs:collapse=Apache License

package v1

import (
	"github.com/robfig/cron"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

/*
Next, we'll setup a logger for the webhooks.
*/

var cronjoblog = logf.Log.WithName("cronjob-resource")

func (r *CronJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

/*
Notice that we use kubebuilder markers to generate webhook manifests.
This markers is responsible for generating a mutating webhook manifest.

The meaning of each marker can be found [here](../TODO.md).
*/

// +kubebuilder:webhook:path=/mutate-batch-tutorial-kubebuilder-io-tutorial-kubebuilder-io-v1-cronjob,mutating=true,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=mcronjob.kb.io

/*
We use the `webhook.Defaulter` interface to set defaults to our CRD.
A webhook will automatically be served that calls this defaulting.

The Default method is expected to mutate the receiver, setting the defaults.
*/

var _ webhook.Defaulter = &CronJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CronJob) Default() {
	cronjoblog.Info("default", "name", r.Name)

	if r.Spec.ConcurrencyPolicy == "" {
		r.Spec.ConcurrencyPolicy = AllowConcurrent
	}
	if r.Spec.Suspend == nil {
		r.Spec.Suspend = new(bool)
	}
	if r.Spec.SuccessfulJobsHistoryLimit == nil {
		r.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*r.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if r.Spec.FailedJobsHistoryLimit == nil {
		r.Spec.FailedJobsHistoryLimit = new(int32)
		*r.Spec.FailedJobsHistoryLimit = 1
	}
}

/*
Notice that we use kubebuilder markers to generate webhook manifests.
This markers is responsible for generating a validating webhook manifest.

The meaning of each marker can be found [here](../TODO.md).
*/

// +kubebuilder:webhook:path=/validate-batch-tutorial-kubebuilder-io-tutorial-kubebuilder-io-v1-cronjob,mutating=false,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=vcronjob.kb.io

/*
To validate our CRD beyond what's possible with declarative validation.
Generally, declarative validation should be sufficient, but sometimes more
advanced use cases call for complex validation.

For instance, we'll see below that we use this to validate a well-formed cron
schedule without making up a long regular expression.

If `webhook.Validator` interface is implemented, a webhook will automatically be
served that calls the validation.

The ValidateCreate and ValidateUpdate methods are expected to validate that its
receiver upon creation and update respectively. We separate out ValidateCreate
from ValidateUpdate to allow behavior like making certain fields immutable, so
that they can only be set on creation.
Here, however, we just use the same shared validation.
*/

var _ webhook.Validator = &CronJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateCreate() error {
	cronjoblog.Info("validate create", "name", r.Name)

	return r.validateCronJob()
}

/*
The ValidateUpdate method is expected to validate the receiver upon update.
If not, it return an error that's used as a validation message.
The receiver is the new object and the argument is the old object.
*/

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CronJob) ValidateUpdate(old runtime.Object) error {
	cronjoblog.Info("validate update", "name", r.Name)

	return r.validateCronJob()
}

/*
We validate the name and the spec of the CronJob.
*/

func (c *CronJob) validateCronJob() error {
	allErrs := c.validateCronJobSpec()
	if err := c.validateCronJobName(); err != nil {
		allErrs = append(allErrs, err)
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "batch.tutorial.kubebuilder.io", Kind: "CronJob"},
		c.Name, allErrs)
}

/*
Some fields are declaratively validated by openapi schema.
You can find kubebuilder validation markers (prefixed
with `// +kubebuilder:validation`) in the [API](api-design.md)
You can find all of the kubebuilder supported markers for
declaring validation in [here](../TODO.md).
*/

func (c *CronJob) validateCronJobSpec() field.ErrorList {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateScheduleFormat(
		c.Spec.Schedule,
		field.NewPath("spec").Child("schedule"))...)
	return allErrs
}

/*
Validating the length of a string field can be done declaratively by
the validation schema.

But the `ObjectMeta.Name` field is defined in a shared package under
the apimachinery repo, so we can't declaratively validate it using
the validation schema.
*/

func (c *CronJob) validateCronJobName() *field.Error {
	if len(c.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// The job name length is 63 character like all Kubernetes objects
		// (which must fit in a DNS subdomain). The cronjob controller appends
		// a 11-character suffix to the cronjob (`-$TIMESTAMP`) when creating
		// a job. The job name length limit is 63 characters. Therefore cronjob
		// names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		return field.Invalid(field.NewPath("metadata").Child("name"), c.Name, "must be no more than 52 characters")
	}
	return nil
}

// +kubebuilder:docs-gen:collapse=Validate object name

/*
We'll need to validate the [cron](https://en.wikipedia.org/wiki/Cron) schedule
is well-formatted.
*/

func validateScheduleFormat(schedule string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if _, err := cron.ParseStandard(schedule); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, schedule, err.Error()))
	}

	return allErrs
}
