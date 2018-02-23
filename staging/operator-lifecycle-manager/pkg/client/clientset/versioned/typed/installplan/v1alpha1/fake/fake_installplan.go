/*
Copyright 2018 The Kubernetes Authors.

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

package fake

import (
	v1alpha1 "github.com/coreos-inc/alm/pkg/apis/installplan/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeInstallPlans implements InstallPlanInterface
type FakeInstallPlans struct {
	Fake *FakeInstallplanV1alpha1
	ns   string
}

var installplansResource = schema.GroupVersionResource{Group: "installplan", Version: "v1alpha1", Resource: "installplans"}

var installplansKind = schema.GroupVersionKind{Group: "installplan", Version: "v1alpha1", Kind: "InstallPlan"}

// Get takes name of the installPlan, and returns the corresponding installPlan object, and an error if there is any.
func (c *FakeInstallPlans) Get(name string, options v1.GetOptions) (result *v1alpha1.InstallPlan, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(installplansResource, c.ns, name), &v1alpha1.InstallPlan{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstallPlan), err
}

// List takes label and field selectors, and returns the list of InstallPlans that match those selectors.
func (c *FakeInstallPlans) List(opts v1.ListOptions) (result *v1alpha1.InstallPlanList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(installplansResource, installplansKind, c.ns, opts), &v1alpha1.InstallPlanList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.InstallPlanList{}
	for _, item := range obj.(*v1alpha1.InstallPlanList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested installPlans.
func (c *FakeInstallPlans) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(installplansResource, c.ns, opts))

}

// Create takes the representation of a installPlan and creates it.  Returns the server's representation of the installPlan, and an error, if there is any.
func (c *FakeInstallPlans) Create(installPlan *v1alpha1.InstallPlan) (result *v1alpha1.InstallPlan, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(installplansResource, c.ns, installPlan), &v1alpha1.InstallPlan{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstallPlan), err
}

// Update takes the representation of a installPlan and updates it. Returns the server's representation of the installPlan, and an error, if there is any.
func (c *FakeInstallPlans) Update(installPlan *v1alpha1.InstallPlan) (result *v1alpha1.InstallPlan, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(installplansResource, c.ns, installPlan), &v1alpha1.InstallPlan{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstallPlan), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeInstallPlans) UpdateStatus(installPlan *v1alpha1.InstallPlan) (*v1alpha1.InstallPlan, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(installplansResource, "status", c.ns, installPlan), &v1alpha1.InstallPlan{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstallPlan), err
}

// Delete takes name of the installPlan and deletes it. Returns an error if one occurs.
func (c *FakeInstallPlans) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(installplansResource, c.ns, name), &v1alpha1.InstallPlan{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeInstallPlans) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(installplansResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.InstallPlanList{})
	return err
}

// Patch applies the patch and returns the patched installPlan.
func (c *FakeInstallPlans) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.InstallPlan, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(installplansResource, c.ns, name, data, subresources...), &v1alpha1.InstallPlan{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.InstallPlan), err
}
