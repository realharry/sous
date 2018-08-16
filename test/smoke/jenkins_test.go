//+build smoke

package smoke

import "testing"

func TestJenkins(t *testing.T) {

	m := newRunner(t, matrix().FixedDimension("project", "simple"))

	m.Run("jenkins", func(t *testing.T, f *fixture) {
		client := setupProject(t, f, f.Projects.HTTPServer())

		flags := &sousFlags{kind: "http-service", tag: "1.2.3", cluster: "cluster1", repo: "github.com/user1/repo1"}

		client.MustRun(t, "init", flags.SousInitFlags())

		client.MustRun(t, "jenkins", nil, "-cluster", "cluster1")
	})
}
