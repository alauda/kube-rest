package rest

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type patch struct {
	patchType types.PatchType
	data      []byte
}

// Type implements Patch.
func (s *patch) Type() types.PatchType {
	return s.patchType
}

// Data implements Patch.
func (s *patch) Data(obj Object) ([]byte, error) {
	return s.data, nil
}

// ConstantPatch constructs a new Patch with the given PatchType and data.
func ConstantPatch(patchType types.PatchType, data []byte) Patch {
	return &patch{patchType, data}
}

type mergeFromPatch struct {
	from runtime.Object
}

// Type implements patch.
func (s *mergeFromPatch) Type() types.PatchType {
	return types.MergePatchType
}

// Data implements Patch.
func (s *mergeFromPatch) Data(obj Object) ([]byte, error) {
	originalJSON, err := json.Marshal(s.from)
	if err != nil {
		return nil, err
	}

	modifiedJSON, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return jsonpatch.CreateMergePatch(originalJSON, modifiedJSON)
}

// MergeFrom creates a Patch that patches using the merge-patch strategy with the given object as base.
func MergeFrom(obj runtime.Object) Patch {
	return &mergeFromPatch{obj}
}

// applyPatch uses server-side apply to patch the object.
type applyPatch struct{}

// Type implements Patch.
func (p applyPatch) Type() types.PatchType {
	return types.ApplyPatchType
}

// Data implements Patch.
func (p applyPatch) Data(obj runtime.Object) ([]byte, error) {
	// NB(directxman12): we might technically want to be using an actual encoder
	// here (in case some more performant encoder is introduced) but this is
	// correct and sufficient for our uses (it's what the JSON serializer in
	// client-go does, more-or-less).
	return json.Marshal(obj)
}
