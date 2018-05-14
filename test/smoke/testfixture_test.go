//+build smoke

package smoke

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/opentable/sous/dev_support/sous_qa_setup/desc"
	sous "github.com/opentable/sous/lib"
)

type TestFixture struct {
	EnvDesc     desc.EnvDesc
	Cluster     TestBunchOfSousServers
	Client      *TestClient
	BaseDir     string
	Singularity *Singularity
	// ClusterSuffix is used to add a suffix to each generated cluster name.
	// This can be used to segregate parallel tests.
	ClusterSuffix string
	Parent        *ParallelTestFixture
	TestName      string
}

var clientCommandTimeoutStr = os.Getenv("SMOKE_TEST_CLIENT_TIMEOUT")
var defaultClientCommandTimeout = 4 * time.Minute

func newTestFixture(t *testing.T, parent *ParallelTestFixture, nextAddr func() string) TestFixture {
	t.Helper()
	t.Parallel()
	if testing.Short() {
		t.Skipf("-short flag present")
	}
	sousBin := getSousBin(t)
	envDesc := getEnvDesc(t)
	baseDir := getDataDir(t)

	clusterSuffix := strings.Replace(t.Name(), "/", "_", -1)
	fmt.Fprintf(os.Stdout, "Cluster suffix: %s", clusterSuffix)

	singularity := NewSingularity(envDesc.SingularityURL())
	singularity.ClusterSuffix = clusterSuffix

	state := sous.StateFixture(sous.StateFixtureOpts{
		ClusterCount:  3,
		ManifestCount: 3,
		ClusterSuffix: clusterSuffix,
	})

	addURLsToState(state, envDesc)

	c, err := newBunchOfSousServers(t, state, baseDir, nextAddr)
	if err != nil {
		t.Fatalf("setting up test cluster: %s", err)
	}

	if err := c.Configure(envDesc); err != nil {
		t.Fatalf("configuring test cluster: %s", err)
	}

	if err := c.Start(t, sousBin); err != nil {
		t.Fatalf("starting test cluster: %s", err)
	}

	clientCommandTimeout := defaultClientCommandTimeout
	if clientCommandTimeoutStr != "" {
		clientCommandTimeout, err = time.ParseDuration(clientCommandTimeoutStr)
		if err != nil {
			t.Fatalf("Parsing SMOKE_TEST_CLIENT_TIMEOUT='%s': %s", clientCommandTimeoutStr, err)
		}
	}

	client := makeClient(baseDir, sousBin, clientCommandTimeout)
	primaryServer := "http://" + c.Instances[0].Addr
	if err := client.Configure(primaryServer, envDesc.RegistryName()); err != nil {
		t.Fatal(err)
	}

	return TestFixture{
		Cluster:       *c,
		Client:        client,
		BaseDir:       baseDir,
		Singularity:   singularity,
		ClusterSuffix: clusterSuffix,
		Parent:        parent,
		TestName:      t.Name(),
	}
}

func (f *TestFixture) Stop(t *testing.T) {
	t.Helper()
	f.Cluster.Stop(t)
}

func (f *TestFixture) ReportSuccess(t *testing.T) {
	t.Helper()
	f.Parent.recordTestPassed(t)
}
