package gateway

import (
	"net/http"
)

// Route is a specification and  handler for a REST endpoint.
type Route struct {
	verb string
	path string
	fn   func(http.ResponseWriter, *http.Request)
}

func (g *gatewayAPI) Routes() []*Route {
	return []*Route{
		{verb: "POST", path: path("function", APIVersion), fn: g.uploadFunction},
		{verb: "PUT", path: path("function/{id}", APIVersion), fn: g.updateFunction},
		{verb: "GET", path: path("function/{id}", APIVersion), fn: g.getFunction},
		{verb: "GET", path: path("function/execute/{id}", APIVersion), fn: g.executeFunction},
		{verb: "DELETE", path: path("function/{id}", APIVersion), fn: g.deleteFunction},
	}
}

func path(route, version string) string {
	return "/" + version + "/" + route
}