package server

type (
	// AdvisoryDeploymentsResource dispatches /servers
	AdvisoryDeploymentsResource struct {
		context ComponentLocator
	}

	// PUTAdvisoryDeploymentsHandler handles GET for /servers
	PUTAdvisoryDeploymentsHandler struct {
		Config      *config.Config
		writer      sous.ClusterWriter
		clusterName string
	}
)

func newAdvisoryDeploymentResource(context ComponentLocator) *AdvisoryDeploymentsResource {
	return &AdvisoryDeploymentsResource{context: context}
}

func (adr *AdvisoryDeploymentsResource) Put(w http.ResponseWriter, req *http.Request, parms httprouter.Params) restful.Exchanger {
}

func (padh *PUTAdvisoryDeploymentsHandler) Exchange() (interface{}, int) {
}
