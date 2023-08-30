package serializer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openim-sigs/component-base/pkg/api/equality"
	metav1 "github.com/openim-sigs/component-base/pkg/apis/meta/v1"
	"github.com/openim-sigs/component-base/pkg/runtime"
	"github.com/openim-sigs/component-base/pkg/runtime/schema"
)

type FakeV1Obj struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

func (*FakeV1Obj) DeepCopyObject() runtime.Object {
	panic("not supported")
}

type FakeV2DifferentObj struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

func (*FakeV2DifferentObj) DeepCopyObject() runtime.Object {
	panic("not supported")
}
func TestSparse(t *testing.T) {
	v1 := schema.GroupVersion{Group: "mygroup", Version: "v1"}
	v2 := schema.GroupVersion{Group: "mygroup", Version: "v2"}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(v1, &FakeV1Obj{})
	scheme.AddKnownTypes(v2, &FakeV2DifferentObj{})
	codecs := NewCodecFactory(scheme)

	srcObj1 := &FakeV1Obj{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}
	srcObj2 := &FakeV2DifferentObj{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}

	encoder := codecs.LegacyCodec(v2, v1)
	decoder := codecs.UniversalDecoder(v2, v1)

	srcObj1Bytes, err := runtime.Encode(encoder, srcObj1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(srcObj1Bytes))
	srcObj2Bytes, err := runtime.Encode(encoder, srcObj2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(srcObj2Bytes))

	uncastDstObj1, err := runtime.Decode(decoder, srcObj1Bytes)
	if err != nil {
		t.Fatal(err)
	}
	uncastDstObj2, err := runtime.Decode(decoder, srcObj2Bytes)
	if err != nil {
		t.Fatal(err)
	}

	// clear typemeta
	uncastDstObj1.(*FakeV1Obj).TypeMeta = metav1.TypeMeta{}
	uncastDstObj2.(*FakeV2DifferentObj).TypeMeta = metav1.TypeMeta{}

	if !equality.Semantic.DeepEqual(srcObj1, uncastDstObj1) {
		t.Fatal(cmp.Diff(srcObj1, uncastDstObj1))
	}
	if !equality.Semantic.DeepEqual(srcObj2, uncastDstObj2) {
		t.Fatal(cmp.Diff(srcObj2, uncastDstObj2))
	}
}
