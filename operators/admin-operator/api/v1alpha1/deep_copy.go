package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (in *AdminAppProfile) DeepCopyInto(out *AdminAppProfile) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
}
func (in *AdminAppProfile) DeepCopy() *AdminAppProfile {
	if in == nil {
		return nil
	}
	out := new(AdminAppProfile)
	in.DeepCopyInto(out)
	return out
}
func (in *AdminAppProfile) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *AdminAppProfileList) DeepCopyInto(out *AdminAppProfileList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AdminAppProfile, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}
func (in *AdminAppProfileList) DeepCopy() *AdminAppProfileList {
	if in == nil {
		return nil
	}
	out := new(AdminAppProfileList)
	in.DeepCopyInto(out)
	return out
}
func (in *AdminAppProfileList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *AdminAppProfileStatus) DeepCopyInto(out *AdminAppProfileStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}
