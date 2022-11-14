/*
Copyright 2018, 2019 the Velero contributors.

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

package plugin

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// RestorePlugin is a restore item action plugin for Velero
type RestorePlugin struct {
	log logrus.FieldLogger
}

// NewRestorePlugin instantiates a RestorePlugin.
func NewRestorePlugin(log logrus.FieldLogger) *RestorePlugin {
	return &RestorePlugin{log: log}
}

// AppliesTo returns information about which resources this action should be invoked for.
// The IncludedResources and ExcludedResources slices can include both resources
// and resources with group names. These work: "ingresses", "ingresses.extensions".
// A RestoreItemAction's Execute function will only be invoked on items that match the returned
// selector. A zero-valued ResourceSelector matches all resources.
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"pods"},
	}, nil
}

// Execute allows the RestorePlugin to perform arbitrary logic with the item being restored,
// in this case, setting a custom annotation on the item being restored.
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.log.Infof("Executing Anotations Prune Plugin plugin for Restore %s", input.Restore.Name)

	metadata, err := meta.Accessor(input.Item)
	if err != nil {
		return &velero.RestoreItemActionExecuteOutput{}, err
	}

	pod := new(v1.Pod)
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(input.Item.UnstructuredContent(), pod); err != nil {
		p.log.Error("Error converting item to pod schema", err)
		return &velero.RestoreItemActionExecuteOutput{}, errors.WithStack(err)
	}

	Annotations := metadata.GetAnnotations()
	if Annotations == nil {
		Annotations = make(map[string]string)
	}

	PruneList := []string{
		"ovn.kubernetes.io/allocated",
		"ovn.kubernetes.io/cidr",
		"ovn.kubernetes.io/gateway",
		"ovn.kubernetes.io/ip_address",
		"ovn.kubernetes.io/ip_pool",
		"ovn.kubernetes.io/logical_switch",
		"ovn.kubernetes.io/mac_address",
		"ovn.kubernetes.io/network_type",
		"ovn.kubernetes.io/pod_nic_type",
		"ovn.kubernetes.io/provider_network",
		"ovn.kubernetes.io/routed",
		"ovn.kubernetes.io/vlan_id",
	}

	// Remove blacklisted annotations
	for annotation := range Annotations {
		for _, blacklisted := range PruneList {
			if annotation == blacklisted {
				delete(Annotations, annotation)
				p.log.Infof("Removed annotation %s from pod", annotation, pod.Name)
				break
			}
		}
	}

	// convert back and return the mapped result
	res, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pod)
	if err != nil {
		p.log.Errorf("Error converting item back to unstructured schema")
		return &velero.RestoreItemActionExecuteOutput{}, errors.WithStack(err)
	}
	return velero.NewRestoreItemActionExecuteOutput(&unstructured.Unstructured{Object: res}), nil
}
