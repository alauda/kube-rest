package http

import (
	"context"

	"github.com/alauda/kube-rest/pkg/types"

	types2 "k8s.io/apimachinery/pkg/types"
)

// Interface for http requests
type Interface interface {
	Get(ctx context.Context, absPath string) ([]byte, error)
	List(ctx context.Context, absPath string, option types.Option) ([]byte, error)
	Create(ctx context.Context, absPath string, data []byte, option types.Option) ([]byte, error)
	Update(ctx context.Context, absPath string, data []byte, option types.Option) ([]byte, error)
	Patch(ctx context.Context, absPath string, pt types2.PatchType, data []byte) ([]byte, error)
	Delete(ctx context.Context, absPath string, option types.Option) ([]byte, error)
}
