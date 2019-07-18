package http

import (
	"alauda/kube-rest/pkg/config"
	"alauda/kube-rest/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	types2 "k8s.io/apimachinery/pkg/types"
)

func getJSON(key, val string) []byte {
	return []byte(fmt.Sprintf(`{%q: %q}`, key, val))
}

func getJSONObjectList(obj map[string]string) []byte {
	bt, _ := json.Marshal(obj)
	return bt
}

func getClientServer(h func(http.ResponseWriter, *http.Request)) (Interface, *httptest.Server, error) {
	svr := httptest.NewServer(http.HandlerFunc(h))
	cfg, err := config.GetDefaultConfig(svr.URL)
	if nil != err {
		return nil, nil, err
	}
	client, err := NewForConfig(cfg)
	if nil != err {
		svr.Close()
		return nil, nil, err
	}
	return client, svr, nil
}

func TestList(t *testing.T) {
	cases := []struct {
		method string
		name   string
		path   string
		resp   []byte
		want   []byte
	}{
		{
			name: "normal_list",
			path: "/test/",
			resp: getJSON("a", "b"),
			want: getJSON("a", "b"),
		},
		{
			name: "filtered_list",
			path: "/test/?filter=a",
			resp: getJSONObjectList(map[string]string{
				"a": "b",
				"c": "d",
			}),
			want: getJSONObjectList(map[string]string{"a": "b"}),
		},
	}

	for _, c := range cases {
		var err error
		var path *url.URL
		if path, err = url.Parse(c.path); nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}
		cli, srv, err := getClientServer(func(w http.ResponseWriter, r *http.Request) {
			if "GET" != r.Method {
				t.Errorf("List(%q) got HTTP method %s. wanted GET", c.name, r.Method)
			}
			if r.URL.Path != path.Path {
				t.Errorf("List(%q) got path %s. wanted %s", c.name, r.URL.Path, c.path)
			}
			filter, ok := r.URL.Query()["filter"]
			if ok {
				obj := make(map[string]string)
				res := make(map[string]string)
				if err := json.Unmarshal(c.resp, &obj); nil != err {
					t.Errorf("unexpected error when filtering result: %v", err)
				}
				for _, key := range filter {
					if val, ok := obj[key]; ok {
						res[key] = val
					}
				}
				c.resp = getJSONObjectList(res)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(c.resp)
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		params := types.QueryParameters{}

		for k, values := range path.Query() {
			params[k] = values[0]
		}

		got, err := cli.List(context.TODO(), path.Path, &types.Options{Params: params})

		if nil != err {
			t.Errorf("unexpected error when listing %q: %v", c.name, err)
			continue
		}

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("List(%q) want: %s\ngot: %s", c.name, c.want, got)
		}
	}

}

func TestDelete(t *testing.T) {
	cases := []struct {
		method string
		name   string
		path   string
		resp   []byte
		want   []byte
	}{
		{
			name: "normal_delete",
			path: "/test/",
		},
		{
			name: "filtered_delete",
			path: "/test/?filter=a",
		},
	}

	for _, c := range cases {
		var err error
		var path *url.URL
		if path, err = url.Parse(c.path); nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}
		cli, srv, err := getClientServer(func(w http.ResponseWriter, r *http.Request) {
			if "DELETE" != r.Method {
				t.Errorf("Delete(%q) got HTTP method %s. wanted DELETE", c.name, r.Method)
			}
			if r.URL.Path != path.Path {
				t.Errorf("Delete(%q) got path %s. wanted %s", c.name, r.URL.Path, c.path)
			}
			if !reflect.DeepEqual(r.URL.Query(), path.Query()) {
				t.Errorf("Delete(%q) got query %s. wanted %s", c.name, r.URL.Query(), path.Query())
			}
			w.Header().Set("Content-Type", "application/json")
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		params := types.QueryParameters{}

		for k, values := range path.Query() {
			params[k] = values[0]
		}

		_, err = cli.Delete(context.TODO(), path.Path, &types.Options{Params: params})

		if nil != err {
			t.Errorf("unexpected error when deleting %q: %v", c.name, err)
			continue
		}

	}

}

func TestPatch(t *testing.T) {
	cases := []struct {
		method string
		name   string
		path   string
		patch  []byte
		want   []byte
	}{
		{
			name:  "normal_delete",
			path:  "/test/",
			patch: getJSON("a", "b"),
			want:  getJSON("a", "b"),
		},
	}

	for _, c := range cases {
		var err error
		var path *url.URL
		if path, err = url.Parse(c.path); nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}
		cli, srv, err := getClientServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Patch(%q) got HTTP method %s. wanted PATCH", c.name, r.Method)
			}

			if r.URL.Path != path.Path {
				t.Errorf("Patch(%q) got path %s. wanted %s", c.name, r.URL.Path, path.Path)
			}

			content := r.Header.Get("Content-Type")
			if content != string(types2.StrategicMergePatchType) {
				t.Errorf("Patch(%q) got Content-Type %s. wanted %s", c.name, content, types2.StrategicMergePatchType)
			}

			w.Header().Set("Content-Type", "application/json")

			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Patch(%q) unexpected error reading body: %v", c.name, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(data)
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		params := types.QueryParameters{}

		for k, values := range path.Query() {
			params[k] = values[0]
		}

		got, err := cli.Patch(context.TODO(), path.Path, types2.StrategicMergePatchType, c.patch)

		if nil != err {
			t.Errorf("unexpected error when patching %q: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Patch(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}

func TestCreate(t *testing.T) {
	cases := []struct {
		method string
		name   string
		path   string
		data   []byte
		want   []byte
	}{
		{
			name: "normal_delete",
			path: "/test/",
			data: getJSON("a", "b"),
			want: getJSON("a", "b"),
		},
	}

	for _, c := range cases {
		var err error
		var path *url.URL
		if path, err = url.Parse(c.path); nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}
		cli, srv, err := getClientServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Create(%q) got HTTP method %s. wanted POST", c.name, r.Method)
			}

			if r.URL.Path != path.Path {
				t.Errorf("Create(%q) got path %s. wanted %s", c.name, r.URL.Path, path.Path)
			}

			w.Header().Set("Content-Type", "application/json")

			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Create(%q) unexpected error reading body: %v", c.name, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(data)
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		params := types.QueryParameters{}

		for k, values := range path.Query() {
			params[k] = values[0]
		}

		got, err := cli.Create(context.TODO(), path.Path, c.data)

		if nil != err {
			t.Errorf("unexpected error when creating %q: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Create(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		method string
		name   string
		path   string
		data   []byte
		want   []byte
	}{
		{
			name: "normal_delete",
			path: "/test/",
			data: getJSON("a", "b"),
			want: getJSON("a", "b"),
		},
	}

	for _, c := range cases {
		var err error
		var path *url.URL
		if path, err = url.Parse(c.path); nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}
		cli, srv, err := getClientServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Errorf("Update(%q) got HTTP method %s. wanted PUT", c.name, r.Method)
			}

			if r.URL.Path != path.Path {
				t.Errorf("Update(%q) got path %s. wanted %s", c.name, r.URL.Path, path.Path)
			}

			w.Header().Set("Content-Type", "application/json")
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Update(%q) unexpected error reading body: %v", c.name, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(data)
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		params := types.QueryParameters{}

		for k, values := range path.Query() {
			params[k] = values[0]
		}

		got, err := cli.Update(context.TODO(), path.Path, c.data)

		if nil != err {
			t.Errorf("unexpected error when updating %q: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Update(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		method string
		name   string
		path   string
		resp   []byte
		want   []byte
	}{
		{
			name: "normal_get",
			path: "/test/",
			resp: getJSON("a", "b"),
			want: getJSON("a", "b"),
		},
	}

	for _, c := range cases {
		var err error
		var path *url.URL
		if path, err = url.Parse(c.path); nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}
		cli, srv, err := getClientServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Get(%q) got HTTP method %s. wanted GET", c.name, r.Method)
			}

			if r.URL.Path != c.path {
				t.Errorf("Get(%q) got path %s. wanted %s", c.name, r.URL.Path, path.Path)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(c.resp)
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		params := types.QueryParameters{}

		for k, values := range path.Query() {
			params[k] = values[0]
		}

		got, err := cli.Get(context.TODO(), path.Path)

		if nil != err {
			t.Errorf("unexpected error when gett %q: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Get(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}
