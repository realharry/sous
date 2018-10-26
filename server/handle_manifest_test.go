package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/opentable/sous/lib"
	"github.com/opentable/sous/util/logging"
	"github.com/opentable/sous/util/restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryValuesToManifestIDHappyPath(t *testing.T) {
	assert := assert.New(t)

	pq := func(s string) restful.QueryValues {
		v, _ := url.ParseQuery(s)
		return restful.QueryValues{Values: v}
	}
	ev := func(x interface{}, e error) error {
		return e
	}
	mid := func(v sous.ManifestID, e error) sous.ManifestID {
		return v
	}

	assert.NoError(ev(manifestIDFromValues(pq("repo=gh1"))))
	assert.NoError(ev(manifestIDFromValues(pq("repo=gh1&offset=o1"))))
	assert.NoError(ev(manifestIDFromValues(pq("repo=gh1&offset=o1&flavor=f1"))))

	assert.Equal(
		mid(manifestIDFromValues(pq("repo=gh1"))),
		sous.ManifestID{Source: sous.SourceLocation{Repo: "gh1"}})

	assert.Equal(
		mid(manifestIDFromValues(pq("repo=gh1&offset=o1"))),
		sous.ManifestID{Source: sous.SourceLocation{Repo: "gh1", Dir: "o1"}})

	assert.Equal(
		mid(manifestIDFromValues(pq("repo=gh1&offset=o1&flavor=f1"))),
		sous.ManifestID{Source: sous.SourceLocation{Repo: "gh1", Dir: "o1"}, Flavor: "f1"})
}

func TestQueryValuesToManifestIDSadPath(t *testing.T) {
	assert := assert.New(t)

	pq := func(s string) restful.QueryValues {
		v, _ := url.ParseQuery(s)
		return restful.QueryValues{Values: v}
	}
	ev := func(x interface{}, e error) error {
		return e
	}

	assert.Error(ev(manifestIDFromValues(pq(""))))
	assert.Error(ev(manifestIDFromValues(pq("repo=gh1&repo=gh2"))))
	assert.Error(ev(manifestIDFromValues(pq("repo=gh1&offset=o1&offset=o2"))))
	assert.Error(ev(manifestIDFromValues(pq("repo=gh1&offset=o1&flavor=f1&flavor=f2"))))
}

func TestHandlesManifestGetNotKnown(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	q, err := url.ParseQuery("repo=gh")
	require.NoError(err)

	th := &GETManifestHandler{
		State:       sous.NewState(),
		QueryValues: restful.QueryValues{Values: q},
	}
	_, status := th.Exchange()
	assert.Equal(404, status)
}

func TestHandlesManifestGetBadURL(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	q, err := url.ParseQuery("repo=gh&repo=gh")
	require.NoError(err)
	state := sous.NewState()
	state.Manifests.Add(&sous.Manifest{
		Owners: []string{"some", "owners"},
		Source: sous.SourceLocation{Repo: "gh"},
	})

	th := &GETManifestHandler{
		State:       state,
		QueryValues: restful.QueryValues{Values: q},
	}
	_, status := th.Exchange()
	assert.Equal(404, status)

}

func TestHandlesManifestGet(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	q, err := url.ParseQuery("repo=gh")
	require.NoError(err)
	state := sous.NewState()
	state.Manifests.Add(&sous.Manifest{
		Owners: []string{"some", "owners"},
		Source: sous.SourceLocation{Repo: "gh"},
	})

	th := &GETManifestHandler{
		State:       state,
		QueryValues: restful.QueryValues{Values: q},
	}
	_, status := th.Exchange()
	assert.Equal(status, 200)

}

func TestHandlesManifestPut(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	q, err := url.ParseQuery("repo=gh")
	require.NoError(err)
	state := sous.NewState()
	state.Manifests.Add(&sous.Manifest{
		Owners: []string{"sam", "judson"},
		Source: sous.SourceLocation{Repo: "gh"},
		Kind:   sous.ManifestKindService,
	})
	writer := &sous.DummyStateManager{State: state}

	uripath := "certainly/i/am/healthy"

	manifest := &sous.Manifest{
		Source: sous.SourceLocation{Repo: "gh"},
		Owners: []string{"sam", "judson"},
		Kind:   sous.ManifestKindService,
		Deployments: sous.DeploySpecs{
			"ci": sous.DeploySpec{
				DeployConfig: sous.DeployConfig{
					Resources: sous.Resources{
						"cpus":   "0.1",
						"memory": "100",
						"ports":  "1",
					},
					Startup: sous.Startup{
						CheckReadyProtocol: "HTTP",
						CheckReadyURIPath:  uripath,
					},
					SingularityRequestID: "custom-sing-req-id",
				},
			},
		},
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.Encode(manifest)
	req, err := http.NewRequest("PUT", "", buf)
	require.NoError(err)

	log, _ := logging.NewLogSinkSpy()

	th := &PUTManifestHandler{
		Request:     req,
		StateWriter: writer,
		State:       state,
		QueryValues: restful.QueryValues{Values: q},
		LogSink:     log,
	}

	data, status := th.Exchange()
	assert.Equal(status, 200)
	require.IsType(&sous.Manifest{}, data)
	assert.Len(data.(*sous.Manifest).Owners, 2)
	assert.Equal(data.(*sous.Manifest).Owners[1], "judson")
	assert.Equal(data.(*sous.Manifest).Deployments["ci"].Startup.CheckReadyURIPath, uripath)

	changed, found := state.Manifests.Get(sous.ManifestID{Source: sous.SourceLocation{Repo: "gh"}})
	require.True(found)
	assert.Len(changed.Owners, 2)
	assert.Equal(changed.Owners[1], "judson")
	assert.Equal(changed.Deployments["ci"].SingularityRequestID, "custom-sing-req-id")

}
