package types

import "net/url"

type QueryParameters map[string]string

type Options struct {
	Header url.Values
	Params QueryParameters
}
