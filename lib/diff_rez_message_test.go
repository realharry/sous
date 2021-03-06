package sous

import (
	"fmt"
	"testing"

	"github.com/opentable/sous/util/logging"
)

func TestDiffRezMessage(t *testing.T) {
	msg := &diffRezMessage{
		callerInfo: logging.GetCallerInfo(),
		resolution: &DiffResolution{
			DeploymentID: DeploymentID{
				ManifestID: ManifestID{
					Flavor: "", Source: SourceLocation{Repo: "github.com/opentable/example", Dir: ""},
				},
				Cluster: "test-cluster",
			},
			Desc:  ModifyDiff,
			Error: WrapResolveError(fmt.Errorf("dumb test error")),
		},
	}

	logging.AssertMessageFields(t, msg, logging.StandardVariableFields, map[string]interface{}{
		"@loglov3-otl":                 logging.SousDiffResolution,
		"sous-resolution-errortype":    "*errors.errorString",
		"sous-resolution-errormessage": "dumb test error",
		"sous-deployment-id":           "test-cluster:github.com/opentable/example",
		"sous-manifest-id":             "github.com/opentable/example",
		"sous-resolution-description":  "updated",
	})
}
