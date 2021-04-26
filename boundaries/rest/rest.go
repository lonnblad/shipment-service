package rest

import (
	"fmt"
	"net/http"
	"net/url"

	health "github.com/InVisionApp/go-health/v2"
	health_handlers "github.com/InVisionApp/go-health/v2/handlers"
	"github.com/gorilla/mux"
	otel_mux "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/lonnblad/shipment-service-backend/boundaries/rest/utils"
	v1 "github.com/lonnblad/shipment-service-backend/boundaries/rest/v1"
	"github.com/lonnblad/shipment-service-backend/businesslogic"
	"github.com/lonnblad/shipment-service-backend/config"
)

type API struct {
	server *http.Server
	router *mux.Router
	err    error
}

// New will take a pointer to the businesslogic and return a pointer to a REST API boundary.
func New(publicURL url.URL, logic *businesslogic.BusinessLogic) (*API, error) {
	api := API{router: mux.NewRouter()}

	api.router.NotFoundHandler = utils.HandlerNotFound()
	api.router.MethodNotAllowedHandler = utils.HandlerMethodNotAllowed()

	api.
		withMiddleware().
		withHealthServer()

	v1SubRouter := api.router.PathPrefix("/v1").Subrouter()
	v1.New(v1SubRouter, publicURL).WithLogic(logic)

	return &api, api.err
}

func (api *API) withHealthServer() *API {
	healthCheck := health.New()

	if err := healthCheck.Start(); err != nil {
		api.err = fmt.Errorf("failed to start healthchecks: %w", err)
		return api
	}

	api.router.
		Path("/healthy/status").
		Handler(health_handlers.NewJSONHandlerFunc(healthCheck, map[string]interface{}{}))

	api.router.
		Path("/healthy").
		Handler(health_handlers.NewBasicHandlerFunc(healthCheck))

	return api
}

// func (api *API) withServiceDocs() *API {
// 	handler := service_docs.Handler()
// 	api.newEndpoint().Insecure().
// 		PathPrefix("/docs/service/").Handler(handler)
// 	api.newEndpoint().Insecure().
// 		PathPrefix("/docs/service").Handler(handler)

// 	return api
// }

// func (api *API) withProblemsDocs() *API {
// 	handler := problems_docs.Handler()
// 	api.newEndpoint().Insecure().
// 		PathPrefix("/problems/").Handler(handler)
// 	api.newEndpoint().Insecure().
// 		PathPrefix("/problems").Handler(handler)

// 	return api
// }

func (api *API) withMiddleware() *API {
	api.router.Use(
		utils.MiddlewareRecovery,
		otel_mux.Middleware(config.GetServiceName()),
	)

	return api
}
