package rest

import (
	"context"

	"github.com/alauda/kube-rest/pkg/http"
	"github.com/alauda/kube-rest/pkg/types"

	"k8s.io/client-go/rest"
)

var _ Client = &client{}

type client struct {
	Client http.Interface
}

// Create implements client.Client
func (c *client) Create(ctx context.Context, obj Object, option types.Option) error {
	data, err := obj.Data()
	if nil != err {
		return err
	}
	data, err = c.Client.Create(ctx, obj.TypeLink(), data, option)
	if nil != err {
		return err
	}
	return obj.Parse(data)
}

// Update implements client.Client
func (c *client) Update(ctx context.Context, obj Object, option types.Option) error {
	data, err := obj.Data()
	if nil != err {
		return err
	}
	data, err = c.Client.Update(ctx, obj.SelfLink(), data, option)
	if nil != err {
		return err
	}
	return obj.Parse(data)
}

func (c *client) Get(ctx context.Context, obj Object) error {
	var bt []byte
	var err error
	if bt, err = c.Client.Get(ctx, obj.SelfLink()); nil != err {
		return err
	}
	return obj.Parse(bt)
}

func (c *client) List(ctx context.Context, obj ObjectList, option types.Option) error {
	var bt []byte
	var err error
	if bt, err = c.Client.List(ctx, obj.TypeLink(), option); nil != err {
		return err
	}
	return obj.Parse(bt)
}

func (c *client) Delete(ctx context.Context, obj Object, option types.Option) error {
	_, err := c.Client.Delete(ctx, obj.SelfLink(), option)
	return err
}

func (c *client) Patch(ctx context.Context, obj Object, patch Patch) error {
	var bt []byte
	var err error
	bt, err = patch.Data(obj)
	if nil != err {
		return err
	}
	if bt, err = c.Client.Patch(ctx, obj.TypeLink(), patch.Type(), bt); nil != err {
		return err
	}
	return obj.Parse(bt)
}

// NewForConfig creates a new rest client
func NewForConfig(cfg *rest.Config) (Client, error) {
	restClient, err := http.NewForConfig(cfg)
	if nil != err {
		return nil, err
	}
	return &client{Client: restClient}, nil
}
