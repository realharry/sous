//+build smoke

package smoke

import (
	"testing"

	sous "github.com/opentable/sous/lib"
)

func TestStress(t *testing.T) {

	pf := pfs.newParallelTestFixture(t)

	fixtureConfigs := []fixtureConfig{
		{dbPrimary: false},
		{dbPrimary: true},
	}

	pf.RunMatrix(fixtureConfigs,

		PTest{Name: "singularity-paused-req", Test: func(t *testing.T, f *TestFixture) {
			client := setupProjectSingleDockerfile(t, f, simpleServer)

			did := sous.DeploymentID{
				ManifestID: sous.ManifestID{
					Source: sous.SourceLocation{
						Repo: "github.com/user1/repo1",
					},
				},
				Cluster: "cluster1",
			}
			reqID := f.Singularity.DefaultReqID(t, did)
			flags := &sousFlags{kind: "http-service", tag: "1", cluster: "cluster1"}
			initBuildDeploy(t, client, flags, setMinimalMemAndCPUNumInst1)
			f.Singularity.PauseRequestForDeployment(t, reqID)

			f.KnownToFailHere(t)

			for {
				client.MustFail(t, "deploy", nil, "-force", "-cluster", "cluster1", "-tag", "1")
			}
		}},
	)

}
