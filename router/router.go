package router

import (
	"github.com/rancher/go-rancher/api"

	"github.com/gorilla/mux"
	"github.com/rancher/v2-api/model"
	"github.com/rancher/v2-api/server"
)

func New(s *server.Server) *mux.Router {
	schemas := model.NewSchema()
	router := mux.NewRouter().StrictSlash(true)

	// API framework routes
	router.Methods("GET").Path("/").Handler(api.VersionsHandler(schemas, "v1", "v2"))
	router.Methods("GET").Path("/v2/schemas").Handler(api.SchemasHandler(schemas))
	router.Methods("GET").Path("/v2/schemas/{id}").Handler(api.SchemaHandler(schemas))
	router.Methods("GET").Path("/v2").Handler(api.VersionHandler(schemas, "v2"))

	f := s.HandlerFunc
	router.Methods("GET").Path("/v2/services").Handler(f(schemas, s.ServiceList))
	router.Methods("GET").Path("/v2/service").Handler(f(schemas, s.ServiceList))

	router.Methods("GET").Path("/v2/stacks").Handler(f(schemas, s.StackList))
	router.Methods("GET").Path("/v2/stack").Handler(f(schemas, s.StackList))
	router.Methods("GET").Path("/v2/stacks/{id}").Handler(f(schemas, s.StackByID))
	router.Methods("GET").Path("/v2/stack/{id}").Handler(f(schemas, s.StackByID))
	router.Methods("POST").Path("/v2/stacks").Handler(f(schemas, s.StackCreate))

	router.Methods("GET").Path("/v2/nodes").Handler(f(schemas, s.NodeList))
	router.Methods("GET").Path("/v2/node").Handler(f(schemas, s.NodeList))
	router.Methods("GET").Path("/v2/nodes/{id}").Handler(f(schemas, s.NodeByID))
	router.Methods("GET").Path("/v2/node/{id}").Handler(f(schemas, s.NodeByID))
	router.Methods("POST").Path("/v2/nodes").Handler(f(schemas, s.NodeCreate))

	return router
}
