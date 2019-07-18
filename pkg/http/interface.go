package http

import (
	"alauda/kube-rest/pkg/types"
	"context"

	types2 "k8s.io/apimachinery/pkg/types"
)

type Interface interface {
	Get(ctx context.Context, absPath string) ([]byte, error)
	List(ctx context.Context, absPath string, options *types.Options) ([]byte, error)
	Create(ctx context.Context, absPath string, data []byte) ([]byte, error)
	Update(ctx context.Context, absPath string, data []byte) ([]byte, error)
	Patch(ctx context.Context, absPath string, pt types2.PatchType, data []byte) ([]byte, error)
	Delete(ctx context.Context, absPath string, options *types.Options) ([]byte, error)
}
