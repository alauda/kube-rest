package config

import (
	"net/url"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

const (
	DefaultAcceptContentType = ""
	DefaultContentType       = ""
	DefaultUserAgent         = "kube-rest"
	DefaultApiPath           = ""
	DefaultClientBurst       = 30.0
	DefaultClientQPS         = 20.0
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func GetDefaultConfig(server string) (*rest.Config, error) {
	_, err := url.Parse(server)
	if nil != err {
		return nil, err
	}
	config := &rest.Config{
		Host: server,
		ContentConfig: rest.ContentConfig{
			AcceptContentTypes:   DefaultAcceptContentType,
			ContentType:          DefaultContentType,
			NegotiatedSerializer: Codecs,
			GroupVersion:         &schema.GroupVersion{},
		},
		UserAgent: DefaultUserAgent,
		APIPath:   DefaultApiPath,
		Burst:     DefaultClientBurst,
		QPS:       DefaultClientQPS,
	}
	return config, nil
}

func GetConfigOrDie(server string) *rest.Config {
	cfg, err := GetDefaultConfig(server)
	if nil != err {
		klog.Fatalf("unable to get kubeconfig, err=%s", err.Error())
		return nil
	}
	return cfg
}
