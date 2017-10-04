// +build !ignore_autogenerated

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1alpha1

import (
	json "encoding/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	reflect "reflect"
)

// GetGeneratedDeepCopyFuncs returns the generated funcs, since we aren't registering them.
//
// Deprecated: deepcopy registration will go away when static deepcopy is fully implemented.
func GetGeneratedDeepCopyFuncs() []conversion.GeneratedDeepCopyFunc {
	return []conversion.GeneratedDeepCopyFunc{
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*AppLink).DeepCopyInto(out.(*AppLink))
			return nil
		}, InType: reflect.TypeOf(&AppLink{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ClusterServiceVersion).DeepCopyInto(out.(*ClusterServiceVersion))
			return nil
		}, InType: reflect.TypeOf(&ClusterServiceVersion{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ClusterServiceVersionList).DeepCopyInto(out.(*ClusterServiceVersionList))
			return nil
		}, InType: reflect.TypeOf(&ClusterServiceVersionList{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ClusterServiceVersionSpec).DeepCopyInto(out.(*ClusterServiceVersionSpec))
			return nil
		}, InType: reflect.TypeOf(&ClusterServiceVersionSpec{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*ClusterServiceVersionStatus).DeepCopyInto(out.(*ClusterServiceVersionStatus))
			return nil
		}, InType: reflect.TypeOf(&ClusterServiceVersionStatus{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*CustomResourceDefinitions).DeepCopyInto(out.(*CustomResourceDefinitions))
			return nil
		}, InType: reflect.TypeOf(&CustomResourceDefinitions{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Icon).DeepCopyInto(out.(*Icon))
			return nil
		}, InType: reflect.TypeOf(&Icon{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Maintainer).DeepCopyInto(out.(*Maintainer))
			return nil
		}, InType: reflect.TypeOf(&Maintainer{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*NamedInstallStrategy).DeepCopyInto(out.(*NamedInstallStrategy))
			return nil
		}, InType: reflect.TypeOf(&NamedInstallStrategy{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*RequirementStatus).DeepCopyInto(out.(*RequirementStatus))
			return nil
		}, InType: reflect.TypeOf(&RequirementStatus{})},
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppLink) DeepCopyInto(out *AppLink) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppLink.
func (in *AppLink) DeepCopy() *AppLink {
	if in == nil {
		return nil
	}
	out := new(AppLink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterServiceVersion) DeepCopyInto(out *ClusterServiceVersion) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterServiceVersion.
func (in *ClusterServiceVersion) DeepCopy() *ClusterServiceVersion {
	if in == nil {
		return nil
	}
	out := new(ClusterServiceVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterServiceVersion) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterServiceVersionList) DeepCopyInto(out *ClusterServiceVersionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterServiceVersion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterServiceVersionList.
func (in *ClusterServiceVersionList) DeepCopy() *ClusterServiceVersionList {
	if in == nil {
		return nil
	}
	out := new(ClusterServiceVersionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterServiceVersionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterServiceVersionSpec) DeepCopyInto(out *ClusterServiceVersionSpec) {
	*out = *in
	in.InstallStrategy.DeepCopyInto(&out.InstallStrategy)
	out.Version = in.Version
	in.CustomResourceDefinitions.DeepCopyInto(&out.CustomResourceDefinitions)
	if in.Permissions != nil {
		in, out := &in.Permissions, &out.Permissions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Keywords != nil {
		in, out := &in.Keywords, &out.Keywords
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Maintainers != nil {
		in, out := &in.Maintainers, &out.Maintainers
		*out = make([]Maintainer, len(*in))
		copy(*out, *in)
	}
	if in.Links != nil {
		in, out := &in.Links, &out.Links
		*out = make([]AppLink, len(*in))
		copy(*out, *in)
	}
	out.Icon = in.Icon
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.LabelSelector)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterServiceVersionSpec.
func (in *ClusterServiceVersionSpec) DeepCopy() *ClusterServiceVersionSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterServiceVersionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterServiceVersionStatus) DeepCopyInto(out *ClusterServiceVersionStatus) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	if in.RequirementStatus != nil {
		in, out := &in.RequirementStatus, &out.RequirementStatus
		*out = make([]RequirementStatus, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterServiceVersionStatus.
func (in *ClusterServiceVersionStatus) DeepCopy() *ClusterServiceVersionStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterServiceVersionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitions) DeepCopyInto(out *CustomResourceDefinitions) {
	*out = *in
	if in.Owned != nil {
		in, out := &in.Owned, &out.Owned
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Required != nil {
		in, out := &in.Required, &out.Required
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitions.
func (in *CustomResourceDefinitions) DeepCopy() *CustomResourceDefinitions {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Icon) DeepCopyInto(out *Icon) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Icon.
func (in *Icon) DeepCopy() *Icon {
	if in == nil {
		return nil
	}
	out := new(Icon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Maintainer) DeepCopyInto(out *Maintainer) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Maintainer.
func (in *Maintainer) DeepCopy() *Maintainer {
	if in == nil {
		return nil
	}
	out := new(Maintainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamedInstallStrategy) DeepCopyInto(out *NamedInstallStrategy) {
	*out = *in
	if in.StrategySpecRaw != nil {
		in, out := &in.StrategySpecRaw, &out.StrategySpecRaw
		*out = make(json.RawMessage, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamedInstallStrategy.
func (in *NamedInstallStrategy) DeepCopy() *NamedInstallStrategy {
	if in == nil {
		return nil
	}
	out := new(NamedInstallStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequirementStatus) DeepCopyInto(out *RequirementStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequirementStatus.
func (in *RequirementStatus) DeepCopy() *RequirementStatus {
	if in == nil {
		return nil
	}
	out := new(RequirementStatus)
	in.DeepCopyInto(out)
	return out
}
