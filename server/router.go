package server

import (
	"fmt"
	"net/http"
)

// Route is a specification and  handler for a REST endpoint.
type Route struct {
	verb string
	path string
	fn   func(http.ResponseWriter, *http.Request)
}

const APIVersion = "v1"

func (g *gatewayAPI) Routes() []*Route {
	return []*Route{
		{verb: "GET", path: path("upload", APIVersion), fn: g.uploadFunction},
		{verb: "GET", path: path("execute/{name}", APIVersion), fn: g.execute},
	}
}

func getVersion(route, version string) string {
	path := "/" + version + "/" + route
	fmt.Println("Path - ", path)
	return "/" + version + "/" + route
}

func path(route, version string) string {
	return getVersion(route, version)
}

// 	{verb: "GET", path: "/startVM", fn: },
