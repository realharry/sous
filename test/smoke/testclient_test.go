//+build smoke

package smoke

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/opentable/sous/config"
	"github.com/opentable/sous/ext/docker"
	sous "github.com/opentable/sous/lib"
	"github.com/opentable/sous/util/filemap"
	"github.com/opentable/sous/util/yaml"
)

type TestClient struct {
	Bin
	// Config is set after calling Configure()
	Config config.Config
	// Fixture is the test fixture this client belongs to.
	Fixture *TestFixture
}

func makeClient(fixture *TestFixture, baseDir, sousBin string) *TestClient {
	baseDir = path.Join(baseDir, "client1")
	c := &TestClient{
		Bin: Bin{
			BaseDir:   baseDir,
			BinPath:   sousBin,
			ConfigDir: path.Join(baseDir, "config"),
			LogDir:    path.Join(baseDir, "logs"),
			Env:       map[string]string{},
		},
		Fixture: fixture,
	}
	c.Bin.MassageArgs = c.insertClusterSuffix
	return c
}

func (c *TestClient) Configure(server, dockerReg, userEmail string) error {
	user := strings.Split(userEmail, "@")
	conf := config.Config{
		Server: server,
		Docker: docker.Config{
			RegistryHost: dockerReg,
		},
		User: sous.User{
			Name:  user[0],
			Email: userEmail,
		},
	}
	conf.PollIntervalForClient = 600

	clientDebug := os.Getenv("SOUS_CLIENT_DEBUG") == "true"

	if clientDebug {
		conf.Logging.Basic.Level = "ExtraDebug1"
		conf.Logging.Basic.DisableConsole = false
		conf.Logging.Basic.ExtraConsole = true
	}

	c.Config = conf

	configYAML, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	c.Env["SOUS_CONFIG_DIR"] = c.Bin.ConfigDir

	return c.Bin.Configure(filemap.FileMap{
		"config.yaml": string(configYAML),
	})
}

func (c *TestClient) insertClusterSuffix(t *testing.T, args []string) []string {
	t.Helper()
	for i, s := range args {
		if s == "-cluster" && len(args) > i+1 {
			args[i+1] = c.Fixture.IsolatedClusterName(args[i+1])
		}
		if s == "-tag" && len(args) > i+1 {
			args[i+1] = c.Fixture.IsolatedVersionTag(t, args[i+1])
		}
	}
	return args
}

// TransformManifest applies each of transforms in order to the retrieved
// manifest, then calls 'sous manifest set' to apply them. Any failure is fatal.
func (c *TestClient) TransformManifest(t *testing.T, flags *sousFlags, transforms ...ManifestTransform) {
	t.Helper()
	flags = flags.ManifestIDFlags()
	manifest := c.MustRun(t, "manifest get", flags.ManifestIDFlags())
	var m sous.Manifest
	if err := yaml.Unmarshal([]byte(manifest), &m); err != nil {
		t.Fatalf("manifest get returned invalid YAML: %s\nInvalid YAML was:\n%s", err, manifest)
	}
	for _, f := range transforms {
		m = f(m)
	}
	manifestBytes, err := yaml.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal updated manifest: %s\nInvalid manifest was:\n% #v", err, m)
	}
	manifestSetCmd := c.configureCommand(t, "manifest set", flags)
	defer manifestSetCmd.Cancel()
	manifestSetCmd.Cmd.Stdin = ioutil.NopCloser(bytes.NewReader(manifestBytes))
	if err := manifestSetCmd.runWithTimeout(3 * time.Minute); err != nil {
		t.Fatalf("manifest set failed: %s; output:\n%s", err, manifestSetCmd.executed.Combined)
	}
}

func (c *TestClient) setSingularityRequestID(t *testing.T, clusterName, singReqID string) ManifestTransform {
	return func(m sous.Manifest) sous.Manifest {
		clusterName := c.Fixture.IsolatedClusterName(clusterName)
		d, ok := m.Deployments[clusterName]
		if !ok {
			t.Fatalf("no deployment for %q", clusterName)
		}
		c := d.DeployConfig
		c.SingularityRequestID = singReqID
		d.DeployConfig = c
		m.Deployments[clusterName] = d
		m.Deployments = setMemAndCPUForAll(m.Deployments)
		return m
	}
}
