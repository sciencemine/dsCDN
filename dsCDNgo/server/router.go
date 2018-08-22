package server

/*
 * Router is where all the api routes are defined
 */
import (
	"net/http"

	"github.com/Tkdefender88/dsCDN/logger"
	"github.com/gorilla/mux"
)

var controller = &Controller{Repository: Repository{}}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Queries     []string
	HandlerFunc http.HandlerFunc
}

//Routes defines the list of routes in the API
type Routes []Route

var routes = Routes{
	Route{
		"GetDsms",
		"GET",
		"/dsm",
		[]string{},
		controller.GetDsms,
	},
	Route{
		"GetDsm",
		"GET",
		"/dsm/{id}",
		[]string{},
		controller.GetDsm,
	},
	Route{
		"AddDsm",
		"POST",
		"/dsm",
		[]string{},
		controller.AddDsm,
	},
	Route{
		"GetCe",
		"GET",
		"/ce/{id}",
		[]string{},
		controller.GetCe,
	},
	Route{
		"GetCes",
		"GET",
		"/ce",
		[]string{},
		controller.GetCes,
	},
	Route{
		"GetAssets",
		"GET",
		"/asset",
		[]string{},
		controller.GetAssets,
	},
	Route{
		"AddCe",
		"POST",
		"/ce",
		[]string{},
		controller.AddCe,
	},
	Route{
		"GetPaths",
		"GET",
		"/path",
		[]string{},
		controller.GetPaths,
	},
	Route{
		"GetPath",
		"GET",
		"/path/{id}",
		[]string{},
		controller.GetPath,
	},
	Route{
		"AddPath",
		"POST",
		"/path",
		[]string{},
		controller.AddPath,
	},
	Route{
		"AddAsset",
		"POST",
		"/asset",
		[]string{},
		controller.AddAsset,
	},
	Route{
		"UpdatePath",
		"PUT",
		"/path/{id}",
		[]string{},
		controller.UpdatePath,
	},
	Route{
		"GetMultipartCe",
		"GET",
		"/ce/{id}/",
		[]string{"type", "{[0-9]*?}"},
		controller.GetCe,
	},
	Route{
		"UpdateDsm",
		"PUT",
		"/dsm/{id}",
		[]string{},
		controller.UpdateDsm,
	},
}

//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		if len(route.Queries) > 0 {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Queries(route.Queries...).
				Name(route.Name).
				Handler(handler)
		} else {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
	return router
}
