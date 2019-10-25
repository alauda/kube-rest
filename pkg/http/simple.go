package http

import (
	"context"
	"errors"

	"github.com/alauda/kube-rest/pkg/types"

	types2 "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

type httpClient struct {
	Client *rest.RESTClient
}

// NewForConfig returns http client interface
func NewForConfig(cfg *rest.Config) (Interface, error) {
	if nil == cfg {
		return nil, errors.New("nil rest config")
	}
	restCli, err := rest.RESTClientFor(cfg)
	if nil != err {
		return nil, err
	}
	return &httpClient{Client: restCli}, nil
}

func (c *httpClient) Get(ctx context.Context, absPath string) ([]byte, error) {
	req := c.Client.Get().AbsPath(absPath)
	if nil != ctx {
		req = req.Context(ctx)
	}
	return req.DoRaw()
}

func (c *httpClient) List(ctx context.Context, absPath string, option types.Option) ([]byte, error) {
	req := c.Client.Get().AbsPath(absPath)
	if nil != ctx {
		req = req.Context(ctx)
	}
	if nil != option {
		req = option.ApplyToRequest(req)
	}
	return req.DoRaw()
}

func (c *httpClient) Create(ctx context.Context, absPath string, outBytes []byte, option types.Option) ([]byte, error) {
	req := c.Client.Post().AbsPath(absPath)
	if nil != ctx {
		req = req.Context(ctx)
	}
	if nil != option {
		req = option.ApplyToRequest(req)
	}
	req.Body(outBytes)
	return req.DoRaw()
}

func (c *httpClient) Update(ctx context.Context, absPath string, outBytes []byte, option types.Option) ([]byte, error) {
	req := c.Client.Put().AbsPath(absPath)
	if nil != ctx {
		req = req.Context(ctx)
	}
	if nil != option {
		req = option.ApplyToRequest(req)
	}
	req.Body(outBytes)
	return req.DoRaw()
}

func (c *httpClient) Patch(ctx context.Context, absPath string, pt types2.PatchType, outBytes []byte) ([]byte, error) {
	req := c.Client.Patch(pt).AbsPath(absPath)
	if nil != ctx {
		req = req.Context(ctx)
	}
	req.Body(outBytes)
	return req.DoRaw()
}

func (c *httpClient) Delete(ctx context.Context, absPath string, option types.Option) ([]byte, error) {
	req := c.Client.Delete().AbsPath(absPath)
	if nil != ctx {
		req = req.Context(ctx)
	}
	if nil != option {
		req = option.ApplyToRequest(req)
	}
	return req.DoRaw()
}
