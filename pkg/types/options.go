package types

import (
	"net/url"

	"k8s.io/client-go/rest"
)

// QueryParameters ...
type QueryParameters map[string]string

// Option provide optional entities that are not frequently changed while making requests
type Option interface {
	ApplyToRequest(req *rest.Request) *rest.Request
}

// Options ...
type Options struct {
	Header url.Values
	Params QueryParameters
}

// ApplyToRequest apply options to rest request
func (options *Options) ApplyToRequest(req *rest.Request) *rest.Request {
	if nil != options {
		if headers := options.Header; nil != headers {
			for k, v := range headers {
				req = req.SetHeader(k, v...)
			}
		}
		if params := options.Params; nil != params {
			for k, v := range params {
				req = req.Param(k, v)
			}
		}
	}
	return req
}
