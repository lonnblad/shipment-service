package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"

	"github.com/lonnblad/shipment-service-backend/boundaries/rest/utils"
	"github.com/lonnblad/shipment-service-backend/trace"
)

// @Summary Create Shipment
// @Description Create Shipment
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param body body CreateShipmentRequest true "Shipment Data"
// @Success 200 {object} CreateShipmentResponse
// @Router /v1/tenants/{tenant_id}/shipments [post]
func (api *API) withCreateShipmentHandler() *API {
	api.router.
		Path(pathShipments).
		Methods(http.MethodPost).
		HandlerFunc(api.createShipmentHandler)

	return api
}

func (api *API) createShipmentHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ctx, span := trace.Tracer().Start(ctx, "v1.createShipmentHandler")
	defer span.End()

	reqData, err := parsedCreateShipmentRequest{}.parse(req)
	if err != nil {
		utils.WrapErrorAndWriteJSONResponse(w, http.StatusBadRequest, err)
	}

	span.SetAttributes(
		attribute.String("req.path.tenant_id", reqData.tenantID.String()),
	)

	internalShipment := reqData.body.toInternal(reqData.tenantID)

	internalShipment, err = api.logic.CreateShipment(ctx, internalShipment)
	if err != nil {
		utils.WrapErrorAndWriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	output := CreateShipmentResponse{}.fromInternal(internalShipment)
	output = output.decorateWithLinks(api.publicURL)

	utils.MarshalAndWriteJSONResponse(w, http.StatusCreated, output)
}

type parsedCreateShipmentRequest struct {
	tenantID uuid.UUID
	body     CreateShipmentRequest
}

func (parsedCreateShipmentRequest) parse(req *http.Request) (_ parsedCreateShipmentRequest, err error) {
	var out parsedCreateShipmentRequest

	params := mux.Vars(req)

	if out.tenantID, err = uuid.Parse(params[keyTenantID]); err != nil {
		err = fmt.Errorf("could not parse tenant ID: %s, error: %w", params[keyTenantID], err)
		return
	}

	if err = utils.UnmarshalRequest(req.Body, &out.body); err != nil {
		err = fmt.Errorf("could not parse request body: %w", err)
		return
	}

	return out, nil
}
