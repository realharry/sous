package smoke

import (
	"strings"
	"testing"
)

// TestErrors tests for error messages.
func TestErrors(t *testing.T) {

	m := newRunner(t, matrix())

	m.Run("sousdeploy-notagflag-nogittag", func(t *testing.T, f *fixture) {
		p := setupProject(t, f, f.Projects.HTTPServer())

		stderr := p.MustFail(t, "deploy", nil, "-cluster", "cluster1")

		want := "you must provide the -tag flag (no git tag found)"

		if !strings.Contains(stderr, want) {
			t.Errorf("got stderr %q; want it to contain %q", stderr, want)
		}
	})

	m.Run("sousdeploy-noclusterflag", func(t *testing.T, f *fixture) {
		p := setupProject(t, f, f.Projects.HTTPServer())

		stderr := p.MustFail(t, "deploy", nil, "-tag", "1.2.3")

		want := "you must provide the -cluster flag"

		if !strings.Contains(stderr, want) {
			t.Errorf("got stderr %q; want it to contain %q", stderr, want)
		}
	})
}
