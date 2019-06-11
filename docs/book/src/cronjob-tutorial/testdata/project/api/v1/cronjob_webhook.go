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
// +kubebuilder:docs-gen:collapse=Apache License

/*
We'll start out with some imports.  You'll see below that we'll need a few more imports
than those scaffolded for us.  We'll talk about each one when we use it.
*/
package v1

import (
	"github.com/robfig/cron"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apimachineryutilvalidation "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

/*
Next, we'll setup logger for the webhooks.
*/

var log = logf.Log.WithName("cronjob-resource")

/*
Notice that we use kubebuilder markers to generate some webhook manifests.
The first and second markers are responsible for generating a mutating webhook manifest and a validating
webhook manifest respectively.

The meaning of each marker can be found [here](../TODO.md).
*/

// +kubebuilder:webhook:path=/mutate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=true,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=mcronjob.kb.io
// +kubebuilder:webhook:path=/validate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=false,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=vcronjob.kb.io

/*
We implement the `webhook.Defaulter` interface for our CronJob CRD.
`webhook.Defaulter` only has one method: `Default()`.
If `webhook.Defaulter` is implemented, a webhook server will be autowired to
serve as the mutating webhook for CronJobs.
*/

var _ webhook.Defaulter = &CronJob{}

// Default implements webhookutil.defaulter so a webhook will be registered for the type
func (c *CronJob) Default() {
	log.Info("defaulting cronjob", "namespace", c.Namespace, "name", c.Name)

	if c.Spec.ConcurrencyPolicy == "" {
		c.Spec.ConcurrencyPolicy = AllowConcurrent
	}
	if c.Spec.Suspend == nil {
		c.Spec.Suspend = new(bool)
	}
	if c.Spec.SuccessfulJobsHistoryLimit == nil {
		c.Spec.SuccessfulJobsHistoryLimit = new(int32)
		*c.Spec.SuccessfulJobsHistoryLimit = 3
	}
	if c.Spec.FailedJobsHistoryLimit == nil {
		c.Spec.FailedJobsHistoryLimit = new(int32)
		*c.Spec.FailedJobsHistoryLimit = 1
	}
}

/*
We implement the `webhook.Validator` interface for our CronJob CRD.
`webhook.Validator` interface has two methods: `ValidateCreate()` and `ValidateUpdate()`.
If `webhook.Validator` is implemented, a webhook server will be autowired to
serve as the validating webhook for CronJobs.
*/

var _ webhook.Validator = &CronJob{}

// ValidateCreate implements webhookutil.validator so a webhook will be registered for the type
func (c *CronJob) ValidateCreate() error {
	log.Info("validate create", "namespace", c.Namespace, "name", c.Name)
	return c.validateCronJob()
}

// ValidateUpdate implements webhookutil.validator so a webhook will be registered for the type
func (c *CronJob) ValidateUpdate(old runtime.Object) error {
	log.Info("validate update", "namespace", c.Namespace, "name", c.Name)
	return c.validateCronJob()
}

/*
We uses the same the validation logic for creation and update.
*/

func (c *CronJob) validateCronJob() error {
	allErrs := c.validateCronJobSpec()
	if len(c.ObjectMeta.Name) > apimachineryutilvalidation.DNS1035LabelMaxLength-11 {
		// The cronjob controller appends a 11-character suffix to the cronjob (`-$TIMESTAMP`) when
		// creating a job. The job name length limit is 63 characters.
		// Therefore cronjob names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata").Child("name"), c.Name, "must be no more than 52 characters"))
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "batch.tutorial.kubebuilder.io", Kind: "CronJob"},
		c.Name, allErrs)
}

func (c *CronJob) validateCronJobSpec() field.ErrorList {
	allErrs := field.ErrorList{}

	spec := c.Spec
	fldPath := field.NewPath("spec")

	if len(spec.Schedule) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("schedule"), ""))
	} else {
		allErrs = append(allErrs, validateScheduleFormat(spec.Schedule, fldPath.Child("schedule"))...)
	}
	if spec.StartingDeadlineSeconds != nil {
		allErrs = append(allErrs, apimachineryvalidation.ValidateNonnegativeField(int64(*spec.StartingDeadlineSeconds), fldPath.Child("startingDeadlineSeconds"))...)
	}
	allErrs = append(allErrs, validateConcurrencyPolicy(&spec.ConcurrencyPolicy, fldPath.Child("concurrencyPolicy"))...)
	//allErrs = append(allErrs, ValidateJobTemplateSpec(&spec.JobTemplate, fldPath.Child("jobTemplate"))...)

	if spec.SuccessfulJobsHistoryLimit != nil {
		// zero is a valid SuccessfulJobsHistoryLimit
		allErrs = append(allErrs, apimachineryvalidation.ValidateNonnegativeField(int64(*spec.SuccessfulJobsHistoryLimit), fldPath.Child("successfulJobsHistoryLimit"))...)
	}
	if spec.FailedJobsHistoryLimit != nil {
		// zero is a valid SuccessfulJobsHistoryLimit
		allErrs = append(allErrs, apimachineryvalidation.ValidateNonnegativeField(int64(*spec.FailedJobsHistoryLimit), fldPath.Child("failedJobsHistoryLimit"))...)
	}

	return allErrs
}

/*
Validate .spec.concurrencyPolicy presents and is valid.
*/

func validateConcurrencyPolicy(concurrencyPolicy *ConcurrencyPolicy, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	switch *concurrencyPolicy {
	case AllowConcurrent, ForbidConcurrent, ReplaceConcurrent:
		break
	case "":
		allErrs = append(allErrs, field.Required(fldPath, ""))
	default:
		validValues := []string{string(AllowConcurrent), string(ForbidConcurrent), string(ReplaceConcurrent)}
		allErrs = append(allErrs, field.NotSupported(fldPath, *concurrencyPolicy, validValues))
	}

	return allErrs
}

/*
Validate the [cron](https://en.wikipedia.org/wiki/Cron) schedule is well-formatted.
*/

func validateScheduleFormat(schedule string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if _, err := cron.ParseStandard(schedule); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, schedule, err.Error()))
	}

	return allErrs
}
