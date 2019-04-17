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

package v2

import (
	"path/filepath"

	"sigs.k8s.io/kubebuilder/pkg/scaffold/input"
)

var _ input.File = &KustomizeCRD{}

// KustomizeCRD scaffolds the Kustomization file in manager folder.
type KustomizeCRD struct {
	input.Input
}

// GetInput implements input.File
func (c *KustomizeCRD) GetInput() (input.Input, error) {
	if c.Path == "" {
		c.Path = filepath.Join("config", "webhook", "kustomization.yaml")
	}
	c.TemplateBody = KustomizeCRDTemplate
	c.Input.IfExistsAction = input.Error
	return c.Input, nil
}

var KustomizeCRDTemplate = `# This kustomization.yaml is not intended to be run by itself,
# since it depends on service name and namespace that are out of this kustomize package.
# It should be run by config/default
resources:
- bases/crew_firstmate.yaml
- bases/creatures_kraken.yaml

patches:
# patches here are for enabling the conversion webhook for each CRD
#- patches/crew_firstmate_patch.yaml
#- patches/creatures_kraken_patch.yaml

# the following config is for teaching kustomize how to do kustomization for CRDs.
 configurations:
- kustomizeconfig.yaml
`
