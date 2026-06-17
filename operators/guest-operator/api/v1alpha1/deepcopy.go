package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (in *GuestAppProfile) DeepCopyInto(out *GuestAppProfile) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *GuestAppProfile) DeepCopy() *GuestAppProfile {
	if in == nil {
		return nil
	}
	out := new(GuestAppProfile)
	in.DeepCopyInto(out)
	return out
}

func (in *GuestAppProfile) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *GuestAppProfileList) DeepCopyInto(out *GuestAppProfileList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GuestAppProfile, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *GuestAppProfileList) DeepCopy() *GuestAppProfileList {
	if in == nil {
		return nil
	}
	out := new(GuestAppProfileList)
	in.DeepCopyInto(out)
	return out
}

func (in *GuestAppProfileList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *GuestAppProfileStatus) DeepCopyInto(out *GuestAppProfileStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}
