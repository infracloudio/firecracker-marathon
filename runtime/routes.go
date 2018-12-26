package runtime

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

func (r *runtimeAPI) Routes() []*Route {
	return []*Route{
		//		{verb: "GET", path: path("upload", APIVersion), fn: g.uploadFunction},
		{verb: "POST", path: path("instance/start", APIVersion), fn: r.startInstance},
		{verb: "GET", path: path("instance/get/{id}", APIVersion), fn: r.getInstance},
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
