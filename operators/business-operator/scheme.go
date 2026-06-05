package businessoperator

import(
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var GroupVersion = schema.GroupVersion{Group: "business.sharebite.dev", Version: "v1alpha1"}
var SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
var AddToScheme = SchemeBuilder.AddToScheme

func init () {
	SchemeBuilder.Register(&BusinessAppProfile{}, &BusinessAppProfileList{})
}