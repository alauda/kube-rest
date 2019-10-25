package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"reflect"
	"testing"

	"github.com/alauda/kube-rest/pkg/config"
	"github.com/alauda/kube-rest/pkg/types"

	types2 "k8s.io/apimachinery/pkg/types"
)

var _ Object = &testObj{}

var defaultOptions = &types.Options{Header: url.Values{"Content-Type": []string{"application/json"}}}

type testObj struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (t *testObj) TypeLink(segments ...string) string {
	return "/test"
}

func (t *testObj) SelfLink(segments ...string) string {
	return path.Join("/test", t.Name)
}

func (t *testObj) Data() ([]byte, error) {
	return json.Marshal(t)
}

func (t *testObj) Parse(bt []byte) error {
	clone := new(testObj)
	if err := json.Unmarshal(bt, clone); nil != err {
		return err
	}
	*t = *clone
	return nil
}

type testObjList struct {
	Items []testObj `json:"items"`
}

func (t *testObjList) TypeLink() string {
	return "/test"
}

func (t *testObjList) Parse(bt []byte) error {
	clone := new(testObjList)
	if err := json.Unmarshal(bt, clone); nil != err {
		return err
	}
	*t = *clone
	return nil
}

func getJSON(name, id string) []byte {
	return []byte(fmt.Sprintf(`{"name":%q,"id":%q}`, name, id))
}

func getJSONList(items ...[]byte) []byte {
	json := fmt.Sprintf(`{"items": [%s]}`, bytes.Join(items, []byte(",")))
	return []byte(json)
}

func getClientServer(h func(http.ResponseWriter, *http.Request)) (Client, *httptest.Server, error) {
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

func TestGet(t *testing.T) {
	cases := []struct {
		name string
		path string
		resp []byte
		want Object
	}{
		{
			name: "normal_get",
			path: "/test/a",
			resp: getJSON("a", "b"),
			want: &testObj{"a", "b"},
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

			if r.URL.Path != path.Path {
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

		got := &testObj{Name: "a"}

		err = cli.Get(context.TODO(), got)

		if nil != err {
			t.Errorf("unexpected error when get %q: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Get(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name string
		path string
		resp []byte
		want ObjectList
	}{
		{
			name: "normal_get",
			path: "/test?filter=a",
			resp: getJSONList(getJSON("a", "b")),
			want: &testObjList{
				Items: []testObj{
					{Name: "a", ID: "b"},
				},
			},
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
				t.Errorf("List(%q) got HTTP method %s. wanted GET", c.name, r.Method)
			}

			if r.URL.Path != path.Path {
				t.Errorf("List(%q) got path %s. wanted %s", c.name, r.URL.Path, path.Path)
			}

			if !reflect.DeepEqual(r.URL.Query(), path.Query()) {
				t.Errorf("List(%q) got query %v. wanted %v", c.name, r.URL.Query(), path.Query())
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(c.resp)
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		got := &testObjList{}
		err = cli.List(context.TODO(), got, &types.Options{Params: map[string]string{"filter": "a"}})

		if nil != err {
			t.Errorf("unexpected error when listing %q: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("List(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name string
		path string
		resp []byte
		want Object
	}{
		{
			name: "normal_delete",
			path: "/test/a?dryRun=false",
			want: &testObj{Name: "a", ID: "b"},
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
			if r.Method != "DELETE" {
				t.Errorf("Delete(%q) got HTTP method %s. wanted DELETE", c.name, r.Method)
			}

			if r.URL.Path != path.Path {
				t.Errorf("Delete(%q) got path %s. wanted %s", c.name, r.URL.Path, path.Path)
			}

			if !reflect.DeepEqual(r.URL.Query(), path.Query()) {
				t.Errorf("Delete(%q) got query %v. wanted %v", c.name, r.URL.Query(), path.Query())
			}

			w.Header().Set("Content-Type", "application/json")
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		got := &testObj{Name: "a"}
		err = cli.Delete(context.TODO(), got, &types.Options{Params: map[string]string{"dryRun": "false"}})

		if nil != err {
			t.Errorf("unexpected error when deleting %q: %v", c.name, err)
			continue
		}
	}
}

func TestCreate(t *testing.T) {
	cases := []struct {
		name string
		path string
		resp []byte
		want Object
	}{
		{
			name: "normal_create",
			path: "/test",
			resp: getJSON("a", "b"),
			want: &testObj{Name: "a", ID: "b"},
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

		got := &testObj{}
		err = cli.Create(context.TODO(), got, defaultOptions)

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
		name string
		path string
		resp []byte
		want Object
	}{
		{
			name: "normal_upate",
			path: "/test/a",
			resp: getJSON("a", "b1"),
			want: &testObj{Name: "a", ID: "b1"},
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

			if !reflect.DeepEqual(c.resp, data) {
				t.Errorf("Update(%q) got data %s. wanted %s", c.name, data, c.resp)
			}
			w.Write(c.resp)
		})

		if nil != err {
			t.Errorf("unexpected error when creating client: %v", err)
			continue
		}

		defer srv.Close()

		got := &testObj{Name: "a", ID: "b1"}
		err = cli.Update(context.TODO(), got, defaultOptions)

		if nil != err {
			t.Errorf("unexpected error when updating %q: %v", c.name, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Update(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}

func TestPatch(t *testing.T) {
	cases := []struct {
		name  string
		path  string
		patch []byte
		resp  []byte
		want  Object
	}{
		{
			name:  "normal_patch",
			path:  "/test",
			patch: []byte(`{"id":"b1"}`),
			resp:  getJSON("a", "b1"),
			want:  &testObj{Name: "a", ID: "b1"},
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

			if !reflect.DeepEqual(c.patch, data) {
				t.Errorf("Patch(%q) got data %s. wanted %s", c.name, data, c.patch)
			}
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

		got := &testObj{}
		err = cli.Patch(context.TODO(), got, ConstantPatch(types2.StrategicMergePatchType, c.patch))

		if nil != err {
			t.Errorf("unexpected error when patching %q: %v", c.name, err)
			continue
		}

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Patch(%q) want: %v\ngot: %v", c.name, c.want, got)
		}
	}
}
