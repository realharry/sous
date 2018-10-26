package sous

import (
	"testing"

	"github.com/samsalisbury/semv"
	"github.com/stretchr/testify/assert"
)

func TestClusterMap(t *testing.T) {
	assert := assert.New(t)

	s := State{
		Defs: Defs{
			Clusters: Clusters{
				"one": &Cluster{},
				"two": &Cluster{},
			},
		},
	}

	m := s.ClusterMap()
	assert.Len(m, 2)
	assert.Contains(m, "one")
	assert.Contains(m, "two")
}

func TestState_Validate(t *testing.T) {

	mid := MustParseManifestID("github.com/user/repo")

	// TODO: Expand the definition of "valid". At the time of initial writing,
	// only the individual manifests are validated without reference to
	// definitions in State.Defs; they should additionally be validated against
	// these definitions.

	validState := &State{
		Manifests: NewManifestsFromMap(map[ManifestID]*Manifest{
			mid: &Manifest{
				Owners: []string{"owners"},
				Source: mid.Source,
				Kind:   ManifestKindService,
				Deployments: DeploySpecs{
					"some-cluster": DeploySpec{
						DeployConfig: DeployConfig{
							Resources: Resources{
								"cpus":   "1",
								"memory": "256",
								"ports":  "1",
							},
							NumInstances: 3,
						},
						Version: semv.MustParse("1"),
					},
				},
			},
		}),
		Defs: Defs{
			Clusters: Clusters{
				"some-cluster": {
					Startup: Startup{
						SkipCheck: true,
					},
				},
			},
		},
	}

	flaws := validState.Validate()
	if len(flaws) != 0 {
		for _, f := range flaws {
			t.Error(f)
		}
		t.Fatalf("got %d flaws; want 0", len(flaws))
	}

}
