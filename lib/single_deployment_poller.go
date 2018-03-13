package sous

import (
	"fmt"
	"time"

	"github.com/opentable/sous/util/logging"
	"github.com/opentable/sous/util/restful"
	"github.com/pkg/errors"
)

type (
	SingleDeploymentPoller struct {
		restful.HTTPClient
		ClusterName, URL         string
		locationFilter, idFilter *ResolveFilter
		User                     User
		httpErrorCount           int
		logs                     logging.LogSink
	}
)

// PollTimeout is the pause between each polling request to /status.
const SDPollTimeout = 1 * time.Second

func NewSingleDeploymentPoller(clusterName, serverURL string, baseFilter *ResolveFilter, user User, logs logging.LogSink) (*SingleDeploymentPoller, error) {
	cl, err := restful.NewClient(serverURL, logs.Child("http"))
	if err != nil {
		return nil, err
	}

	loc := *baseFilter
	loc.Cluster = ResolveFieldMatcher{}
	loc.Tag = ResolveFieldMatcher{}
	loc.Revision = ResolveFieldMatcher{}

	id := *baseFilter
	id.Cluster = ResolveFieldMatcher{}

	return &SingleDeploymentPoller{
		ClusterName:    clusterName,
		URL:            serverURL,
		HTTPClient:     cl,
		locationFilter: &loc,
		idFilter:       &id,
		User:           user,
		logs:           logs.Child(clusterName),
	}, nil
}

// start issues a new /status request, reporting the state as computed.
// c.f. pollOnce.
func (sdp *SingleDeploymentPoller) start(rs chan pollResult, done chan struct{}) {
	rs <- pollResult{url: sdp.URL, stat: ResolveNotPolled}
	pollResult := sdp.pollOnce()
	rs <- pollResult
	ticker := time.NewTicker(SDPollTimeout)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			latest := sdp.pollOnce()
			rs <- latest
		case <-done:
			return
		}
	}
}

func (sdp *SingleDeploymentPoller) result(rs ResolveState, data *statusData, err error) pollResult {
	resolveID := "<none in progress>"
	if data.InProgress != nil {
		resolveID = data.InProgress.Started.String()
	}
	return pollResult{url: sdp.URL, stat: rs, resolveID: resolveID, err: err}
}

func (sdp *SingleDeploymentPoller) pollOnce() pollResult {
	data := &statusData{}
	if _, err := sdp.Retrieve("./status", nil, data, sdp.User.HTTPHeaders()); err != nil {
		reportDebugSubPollerMessage(fmt.Sprintf("%s: error on GET /status: %s", sdp.ClusterName, errors.Cause(err)), sdp.logs)
		reportDebugSubPollerMessage(fmt.Sprintf("%s: %T %+v", sdp.ClusterName, errors.Cause(err), err), sdp.logs)
		sdp.httpErrorCount++
		if sdp.httpErrorCount > 10 {
			return sdp.result(
				ResolveHTTPFailed,
				data,
				fmt.Errorf("more than 10 HTTP errors, giving up; latest error: %s", err),
			)
		}
		return sdp.result(ResolveErredHTTP, data, err)
	}
	sdp.httpErrorCount = 0

	// This serves to maintain backwards compatibility.
	// XXX One day, remove it.
	if data.Completed != nil && len(data.Completed.Intended) == 0 {
		data.Completed.Intended = data.Deployments
	}
	if data.InProgress != nil && len(data.InProgress.Intended) == 0 {
		data.InProgress.Intended = data.Deployments
	}

	currentState, err := sdp.computeState(sdp.stateFeatures("in-progress", data.InProgress))

	if currentState == ResolveNotStarted ||
		currentState == ResolveNotVersion ||
		currentState == ResolvePendingRequest {
		state, err := sdp.computeState(sdp.stateFeatures("completed", data.Completed))
		return sdp.result(state, data, err)
	}

	return sdp.result(currentState, data, err)
}

func (sdp *SingleDeploymentPoller) stateFeatures(kind string, rezState *ResolveStatus) (*Deployment, *DiffResolution) {
	current := diffResolutionForSDP(rezState, sdp.locationFilter)
	srvIntent := serverIntentSDP(rezState, sdp.locationFilter)
	reportDebugSubPollerMessage(fmt.Sprintf("%s reports %s intent to resolve [%v]", sdp.URL, kind, srvIntent), sdp.logs)
	reportDebugSubPollerMessage(fmt.Sprintf("%s reports %s rez: %v", sdp.URL, kind, current), sdp.logs)

	return srvIntent, current
}

func diffResolutionForSDP(rstat *ResolveStatus, rf *ResolveFilter) *DiffResolution {
	if rstat == nil {
		reportDebugSubPollerMessage(fmt.Sprintf("Status was nil - no match for %s", rf), logging.Log)
		return nil
	}
	rezs := rstat.Log
	for _, rez := range rezs {
		reportDebugSubPollerMessage(fmt.Sprintf("Checking resolution for: %#v(%[1]T)", rez.ManifestID), logging.Log)
		if rf.FilterManifestID(rez.ManifestID) {
			reportDebugSubPollerMessage(fmt.Sprintf("Matching intent for %s: %#v", rf, rez), logging.Log)
			return &rez
		}
	}
	reportDebugSubPollerMessage(fmt.Sprintf("No match for %s in %d entries", rf, len(rezs)), logging.Log)
	return nil
}

func serverIntentSDP(rstat *ResolveStatus, rf *ResolveFilter) *Deployment {
	reportDebugSubPollerMessage(fmt.Sprintf("Filtering with %q", rf), logging.Log)
	if rstat == nil {
		reportDebugSubPollerMessage("Nil resolve status!", logging.Log)
		return nil
	}
	reportDebugSubPollerMessage(fmt.Sprintf("Filtering %s", rstat.Intended), logging.Log)

	var dep *Deployment
	for _, d := range rstat.Intended {
		if rf.FilterDeployment(d) {
			if dep != nil {
				reportDebugSubPollerMessage(fmt.Sprintf("With %s we didn't match exactly one deployment.", rf), logging.Log)
				return nil
			}
			dep = d
		}
	}
	reportDebugSubPollerMessage(fmt.Sprintf("Filtering found %s", dep), logging.Log)
	return dep
}

// computeState takes the servers intended deployment, and the stable and
// current DiffResolutions and computes the state of resolution for the
// deployment based on that data.
func (sdp *SingleDeploymentPoller) computeState(srvIntent *Deployment, current *DiffResolution) (ResolveState, error) {
	// In there's no intent for the deployment in the current resolution, we
	// haven't started on it yet. Remember that we've already determined that the
	// most-recent GDM does have the deployment scheduled for this cluster, so it
	// should be picked up in the next cycle.
	if srvIntent == nil {
		return ResolveNotStarted, nil
	}

	// This is a nuanced distinction from the above: the cluster is in the
	// process of resolving a different version than what we're watching for.
	// Again, if it weren't in the freshest GDM, we wouldn't have gotten here.
	// Next cycle! (note that in both cases, we're likely to poll again several
	// times before that cycle starts.)
	if !sdp.idFilter.FilterDeployment(srvIntent) {
		return ResolveNotVersion, nil
	}

	// If there's no DiffResolution yet for our Deployment, then we're still
	// waiting for a relatively recent change to the GDM to be processed. I think
	// this could only happen in the first attempt to resolve a recent change to
	// the GDM, and only before the cluster has gotten a DiffResolution recorded.
	if current == nil {
		return ResolvePendingRequest, nil
	}

	if current.Error != nil {
		// Certain errors in resolution may clear on their own. (Generally
		// speaking, these are HTTP errors from Singularity which we hope/assume
		// will become successes with enough persistence - i.e. on the next
		// resolution cycle, Singularity will e.g. have finished a pending->running
		// transition and be ready to receive a new Deploy)
		if IsTransientResolveError(current.Error) {
			reportDebugSubPollerMessage(fmt.Sprintf("%s: received resolver error %s, retrying", sdp.ClusterName, current.Error), sdp.logs)
			return ResolveErredRez, current.Error
		}
		// Other errors are unlikely to clear by themselves. In this case, log the
		// error for operator action, and consider this sdppoller done as failed.
		reportDebugSubPollerMessage(fmt.Sprintf("%#v", current), sdp.logs)
		reportDebugSubPollerMessage(fmt.Sprintf("%#v", current.Error), sdp.logs)

		subject := ""
		if sdp.locationFilter == nil {
			subject = "<no filter defined>"
		} else {
			sourceLocation, ok := sdp.locationFilter.SourceLocation()
			if ok {
				subject = sourceLocation.String()
			} else {
				subject = sdp.locationFilter.String()
			}
		}

		reportSubPollerMessage(fmt.Sprintf("Deployment of %s to %s failed: %s", subject, sdp.ClusterName, current.Error.String), sdp.logs)
		return ResolveFailed, current.Error
	}

	// In the case where the GDM and ADS deployments are the same, the /status
	// will be described as "unchanged." The upshot is that the most current
	// intend to deploy matches this cluster's current resolver's intend to
	// deploy, and that that matches the deploy that's running. Success!
	if current.Desc == StableDiff {
		return ResolveComplete, nil
	}

	if current.Desc == ComingDiff {
		return ResolveTasksStarting, nil
	}

	return ResolveInProgress, nil
}
