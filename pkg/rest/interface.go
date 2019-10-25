package rest

import (
	"context"

	"github.com/alauda/kube-rest/pkg/types"

	types2 "k8s.io/apimachinery/pkg/types"
)

// Object is the object entity for a rest request.
// It knows how to get request url, parsing object and do the deepcopy.
type Object interface {
	TypeLink(segments ...string) string
	SelfLink(segments ...string) string
	Data() ([]byte, error)
	Parse(bt []byte) error
}

// ObjectList is the object list entity for a rest request.
// It knows how to get request url, and deepcopy the objects.
type ObjectList interface {
	TypeLink() string
	Parse(bt []byte) error
}

// Patch is a patch that can be applied to a rest object.
type Patch interface {
	// Type is the PatchType of the patch.
	Type() types2.PatchType
	// Data is the raw data representing the patch.
	Data(obj Object) ([]byte, error)
}

// Reader knows how to read and list objects.
type Reader interface {
	// Get retrieves an obj for the given object key from the rest object.
	// obj must be a struct pointer so that obj can be updated with the response
	// returned by the Server.
	Get(ctx context.Context, obj Object) error

	// List retrieves list of objects for a given namespace and list options. On a
	// successful call, Items field in the list will be populated with the
	// result returned from the server.
	List(ctx context.Context, obj ObjectList, option types.Option) error
}

// Writer knows how to create, delete, and update rest objects.
type Writer interface {
	// Create saves the object obj in the rest object.
	Create(ctx context.Context, obj Object, option types.Option) error

	// Delete deletes the given obj from rest object.
	Delete(ctx context.Context, obj Object, option types.Option) error

	// Update updates the given obj in the rest object. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Update(ctx context.Context, obj Object, option types.Option) error

	// Update updates the given obj in the rest object. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Patch(ctx context.Context, obj Object, patch Patch) error
}

// Client knows how to perform CRUD operations on Object
type Client interface {
	Reader
	Writer
}
