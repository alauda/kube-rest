package http

import (
	"github.com/alauda/kube-rest/pkg/types"
	"context"
	"errors"

	types2 "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

type httpClient struct {
	Client *rest.RESTClient
}

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
		req.Context(ctx)
	}
	return req.DoRaw()
}

func (c *httpClient) List(ctx context.Context, absPath string, options *types.Options) ([]byte, error) {
	req := c.Client.Get().AbsPath(absPath)
	if nil != ctx {
		req.Context(ctx)
	}
	if nil != options {
		if headers := options.Header; nil != headers {
			for k, v := range headers {
				req.SetHeader(k, v...)
			}
		}
		if params := options.Params; nil != params {
			for k, v := range options.Params {
				req.Param(k, v)
			}
		}
	}
	return req.DoRaw()
}

func (c *httpClient) Create(ctx context.Context, absPath string, outBytes []byte) ([]byte, error) {
	req := c.Client.Post().AbsPath(absPath)
	if nil != ctx {
		req.Context(ctx)
	}
	req = req.Body(outBytes)
	return req.DoRaw()
}

func (c *httpClient) Update(ctx context.Context, absPath string, outBytes []byte) ([]byte, error) {
	req := c.Client.Put().AbsPath(absPath)
	if nil != ctx {
		req.Context(ctx)
	}
	req = req.Body(outBytes)
	return req.DoRaw()
}

func (c *httpClient) Patch(ctx context.Context, absPath string, pt types2.PatchType, outBytes []byte) ([]byte, error) {
	req := c.Client.Patch(pt).AbsPath(absPath)
	if nil != ctx {
		req.Context(ctx)
	}
	req = req.Body(outBytes)
	return req.DoRaw()
}

func (c *httpClient) Delete(ctx context.Context, absPath string, options *types.Options) ([]byte, error) {
	req := c.Client.Delete().AbsPath(absPath)
	if nil != ctx {
		req.Context(ctx)
	}
	if nil != options {
		if headers := options.Header; nil != headers {
			for k, v := range headers {
				req.SetHeader(k, v...)
			}
		}
		if params := options.Params; nil != params {
			for k, v := range params {
				req.Param(k, v)
			}
		}
	}
	return req.DoRaw()
}
