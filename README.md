# kube-rest

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/alauda/kube-rest/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/alauda/kube-rest)](https://goreportcard.com/report/github.com/alauda/kube-rest)

Kube-Rest implement a http client for making restful request with kubernetes client-go.

## Why 
Client-go is not just a client for talking to kubernetes cluster, it is also a good rest client for go:

* multi serilizers: json, portobuf and other serilizers could be added by user
* rate-limiting support: you can specify your QPS for local client
* back-off manager: a back-off manager for unexpected network failover
* human readable/ writeable restful interfaces: the interface is easy to read and write


## How to use

For example, if we want to make a `POST` request to httpbin.org:

```golang
import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/alauda/kube-rest/pkg/config"
	"github.com/alauda/kube-rest/pkg/http"
	"github.com/alauda/kube-rest/pkg/types"
)

func main() {
	address := "http://httpbin.org"
	client, err := http.NewForConfig(config.GetConfigOrDie(address))
	if nil != err {
		log.Fatal(err)
	}

	v := url.Values{}
	v.Set("Content-Type", "application/json")
	option := &types.Options{
		Header: v,
	}

	bt, err := client.Create(context.TODO(), "post", nil, option)

	if nil != err {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("%s", bt))
}
```

Check the [examples](https://github.com/alauda/kube-rest/tree/master/exmaples/https) for more examples.
