package types

import "net/url"

type QueryParameters map[string]string

// Options is an aggregation of options that can be applied to a rest object.
// It's optional when making one request
//type Options interface {
//	// Headers is the headers of the request.
//	Headers() url.Values
//	// QueryParameters is the query parameters. Valid query parameters will be added to the end of request url.
//	QueryParameters() QueryParameters
//}

type Options struct {
	Header url.Values
	Params QueryParameters
}
