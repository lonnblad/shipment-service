package v1

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	swagger_http "github.com/swaggo/http-swagger"

	swagger_docs "github.com/lonnblad/shipment-service-backend/boundaries/rest/v1/generated/swagger"
	"github.com/lonnblad/shipment-service-backend/businesslogic"
)

// @title Shipment Service API
// @version v1.0.0
// @description ### Tenant ID
// @description The API requires a Tenant ID.
// @description For the purpose of testing, use this ID:
// @description **fe131811-7fcd-4942-84a2-4ce8af359da5**
type API struct {
	router    *mux.Router
	logic     *businesslogic.BusinessLogic
	publicURL url.URL
}

func New(router *mux.Router, publicURL url.URL) *API {
	api := &API{router: router, publicURL: publicURL}

	api.
		withCreateShipmentHandler().
		withListShipmentsHandler().
		withGetShipmentHandler().
		withSwagger(publicURL)

	return api
}

func (api *API) WithLogic(logic *businesslogic.BusinessLogic) *API {
	api.logic = logic
	return api
}

func (api *API) withSwagger(publicURL url.URL) *API {
	swagger_docs.SwaggerInfo.Host = publicURL.Host

	publicURL.Path = "/v1/docs/swagger/doc.json"

	api.router.PathPrefix("/docs/swagger/").
		Handler(swagger_http.Handler(
			swagger_http.URL(publicURL.String()),
		))

	publicURL.Path = "/v1/docs/swagger/index.html"

	api.router.PathPrefix("/docs/swagger").
		Handler(
			http.RedirectHandler(
				publicURL.String(),
				http.StatusPermanentRedirect,
			),
		)

	return api
}
