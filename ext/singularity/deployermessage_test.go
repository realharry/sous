package singularity

import (
	"testing"

	"github.com/opentable/sous/util/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployerMessage(t *testing.T) {
	logger, control := logging.NewLogSinkSpy()
	pair := baseDeployablePair()

	reportDeployerMessage("test", pair, nil, nil, logging.InformationLevel, logger)

	logCalls := control.CallsTo("LogMessage")
	require.Len(t, logCalls, 1)
	assert.Equal(t, logCalls[0].PassedArgs().Get(0), logging.InformationLevel)

	//weak check on WriteToConsole
	consoleCalls := control.CallsTo("Console")
	require.Len(t, consoleCalls, 1)
}
