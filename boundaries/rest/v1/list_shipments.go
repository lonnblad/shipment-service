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

const (
	defaultLimitListShipments  = 10
	maxLimitListShipments      = 100
	defaultOffsetListShipments = 0
)

// @Summary List Shipments
// @Description List Shipments
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param limit query int false "Limit" minimum(1) maximum(100) default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} listShipmentsResponse
// @Router /v1/tenants/{tenant_id}/shipments [get]
func (api *API) withListShipmentsHandler() *API {
	api.router.
		Path(pathShipments).
		Methods(http.MethodGet).
		HandlerFunc(api.listShipmentsHandler)

	return api
}

func (api *API) listShipmentsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ctx, span := trace.Tracer().Start(ctx, "v1.listShipmentsHandler")
	defer span.End()

	reqData, err := parsedListShipmentsRequest{}.parse(req)
	if err != nil {
		utils.WrapErrorAndWriteJSONResponse(w, http.StatusBadRequest, err)
	}

	span.SetAttributes(
		attribute.String("req.path.tenant_id", reqData.tenantID.String()),
		attribute.Int("req.query.limit", reqData.limit),
		attribute.Int("req.query.offset", reqData.offset),
	)

	internalShipments, err := api.logic.ListShipments(ctx, reqData.tenantID, reqData.limit, reqData.offset)
	if err != nil {
		utils.WrapErrorAndWriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	output := listShipmentsResponse{}.fromInternal(internalShipments)
	output = output.decorateWithLinks(api.publicURL, reqData)

	utils.MarshalAndWriteJSONResponse(w, http.StatusCreated, output)
}

type parsedListShipmentsRequest struct {
	tenantID uuid.UUID
	limit    int
	offset   int
}

func (parsedListShipmentsRequest) parse(req *http.Request) (_ parsedListShipmentsRequest, err error) {
	var out parsedListShipmentsRequest

	params := mux.Vars(req)

	if out.tenantID, err = uuid.Parse(params[keyTenantID]); err != nil {
		err = fmt.Errorf("could not parse tenant ID: %s, error: %w", params[keyTenantID], err)
		return
	}

	const (
		defaultLimit = defaultLimitListShipments
		maxLimit     = maxLimitListShipments
	)

	limitStr := req.URL.Query().Get("limit")
	if out.limit, err = utils.ParseLimit(limitStr, defaultLimit, maxLimit); err != nil {
		err = fmt.Errorf("could not parse limit: %w", err)
		return
	}

	const defaultOffset = defaultOffsetListShipments

	offsetStr := req.URL.Query().Get("offset")
	if out.offset, err = utils.ParseOffset(offsetStr, defaultOffset); err != nil {
		err = fmt.Errorf("could not parse offset: %w", err)
		return
	}

	return out, nil
}
