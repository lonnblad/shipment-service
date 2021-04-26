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

// @Summary Get Shipment
// @Description Get Shipment
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param shipment_id path string true "Shipment ID"
// @Success 200 {object} getShipmentResponse
// @Router /v1/tenants/{tenant_id}/shipments/{shipment_id} [get]
func (api *API) withGetShipmentHandler() *API {
	api.router.
		Path(pathShipment).
		Methods(http.MethodGet).
		HandlerFunc(api.getShipmentHandler)

	return api
}

func (api *API) getShipmentHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ctx, span := trace.Tracer().Start(ctx, "v1.getShipmentHandler")
	defer span.End()

	reqData, err := parsedGetShipmentRequest{}.parse(req)
	if err != nil {
		utils.WrapErrorAndWriteJSONResponse(w, http.StatusBadRequest, err)
	}

	span.SetAttributes(
		attribute.String("req.path.tenant_id", reqData.tenantID.String()),
		attribute.String("req.path.shipment_id", reqData.shipmentID.String()),
	)

	internalShipment, err := api.logic.GetShipment(ctx, reqData.tenantID, reqData.shipmentID)
	if err != nil {
		utils.WrapErrorAndWriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	output := getShipmentResponse{}.fromInternal(internalShipment)
	output = output.decorateWithLinks(api.publicURL)

	utils.MarshalAndWriteJSONResponse(w, http.StatusCreated, output)
}

type parsedGetShipmentRequest struct {
	tenantID   uuid.UUID
	shipmentID uuid.UUID
}

func (parsedGetShipmentRequest) parse(req *http.Request) (_ parsedGetShipmentRequest, err error) {
	var out parsedGetShipmentRequest

	params := mux.Vars(req)

	if out.tenantID, err = uuid.Parse(params[keyTenantID]); err != nil {
		err = fmt.Errorf("could not parse tenant ID: %s, error: %w", params[keyTenantID], err)
		return
	}

	if out.shipmentID, err = uuid.Parse(params[keyShipmentID]); err != nil {
		err = fmt.Errorf("could not parse shipment ID: %s, error: %w", params[keyShipmentID], err)
		return
	}

	return out, nil
}
