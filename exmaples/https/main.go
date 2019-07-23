package main

import (
	"alauda/kube-rest/pkg/config"
	"alauda/kube-rest/pkg/rest"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

const (
	ListeningAddress = ":8443"
	ServerAddress    = "https://localhost:8443"
	CertFile         = "./cert/server/server.crt"
	KeyFile          = "./cert/server/server.key"
)

func StartServer(stop chan struct{}, logger *log.Logger) {

	s := &http.Server{
		Addr:           ListeningAddress,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/rest", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte(`{"items":[]}`))
	})

	go func() {
		<-stop
		if err := s.Shutdown(context.TODO()); nil != err {
			logger.Printf("HTTP server Shutdown: %v", err)
		} else {
			logger.Printf("HTTP server Shutdown")
		}
		close(stop)
	}()

	logger.Printf("About to listen on 8443. Go to %s", ServerAddress)
	logger.Fatal(s.ListenAndServeTLS(CertFile, KeyFile))
}

var _ rest.Object = &Rest{}
var _ rest.ObjectList = &RestList{}

type RestList struct {
	Items []*Rest `json:"items"`
}

func (r *RestList) AbsPath() string {
	return "/rest"
}

func (r *RestList) Parse(bt []byte) error {
	return json.Unmarshal(bt, &r)
}

type Rest struct {
	Name string `json:"name"`
}

func (r *Rest) AbsPath() string {
	return "/rest"
}

func (r *Rest) AbsObjPath() string {
	return path.Join("/rest", r.Name)
}

func (r *Rest) Data() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Rest) Parse(bt []byte) error {
	return json.Unmarshal(bt, &r)
}

func main() {

	logger := log.New(os.Stdout, "INFO: ", log.Lshortfile)

	logger.Println("Hello, world")

	stop := make(chan struct{})

	defer func() {
		stop <- struct{}{}
	}()

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		<-sig
		stop <- struct{}{}
	}()

	go StartServer(stop, logger)

	cfg := config.GetConfigOrDie(ServerAddress)

	cfg.TLSClientConfig.CAFile = CertFile

	cli, err := rest.NewForConfig(cfg)

	if nil != err {
		logger.Fatal(err)
	}

	obj := &Rest{}

	err = cli.Create(context.TODO(), obj)
	if nil != err {
		log.Fatal(err)
	} else {
		logger.Printf("Create success, obj=%v", obj)
	}

	objList := &RestList{}

	err = cli.List(context.TODO(), objList, nil)
	if nil != err {
		log.Fatal(err)
	} else {
		logger.Printf("List success, obj=%v", objList)
	}

}
