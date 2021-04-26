package businesslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/lonnblad/shipment-service-backend/businesslogic/models"
	"github.com/lonnblad/shipment-service-backend/businesslogic/price"
	"github.com/lonnblad/shipment-service-backend/storage"
	"github.com/lonnblad/shipment-service-backend/trace"
)

type BusinessLogic struct {
	storage storage.ShipmentStorage
}

// New will take a pointer the ShipmentStorage and return a new BusinessLogic instance.
func New(storage storage.ShipmentStorage) *BusinessLogic {
	return &BusinessLogic{storage: storage}
}

func (bl *BusinessLogic) CreateShipment(ctx context.Context, shipment models.Shipment) (_ models.Shipment, err error) {
	ctx, span := trace.Tracer().Start(ctx, "businesslogic.CreateShipment")
	defer span.End()

	span.SetAttributes(
		attribute.String("shipment.tenant_id", shipment.TenantID.String()),
		attribute.String("shipment.sender.country_code", shipment.Sender.CountryCode),
	)

	if err = shipment.Validate(); err != nil {
		err = fmt.Errorf("shipment was invalid: %w", err)
		return
	}

	shipment.ID = uuid.New()
	shipment.CreatedAt = time.Now()

	span.SetAttributes(
		attribute.String("shipment.id", shipment.ID.String()),
		attribute.String("shipment.created_at", shipment.CreatedAt.Format(time.RFC3339)),
		attribute.Int("shipment.package.weight", shipment.Package.Weight),
	)

	shipment.Package.Price, err = price.Calculate(shipment)
	if err != nil {
		err = fmt.Errorf("could not calculate the price of the shipment: %w", err)
		return
	}

	span.SetAttributes(
		attribute.Int("shipment.package.price", shipment.Package.Price),
	)

	dlShipment := shipment.ToDatalayer()
	if err = bl.storage.StoreShipment(ctx, dlShipment); err != nil {
		err = fmt.Errorf("could not create shipment in storage: %w", err)
		return
	}

	return shipment, nil
}

func (bl *BusinessLogic) ListShipments(ctx context.Context, tenantID uuid.UUID, limit, offset int) (_ models.Shipments, err error) {
	ctx, span := trace.Tracer().Start(ctx, "businesslogic.ListShipments")
	defer span.End()

	span.SetAttributes(
		attribute.String("tenant_id", tenantID.String()),
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	)

	dlShipments, err := bl.storage.ListShipments(ctx, tenantID.String(), limit, offset)
	if err != nil {
		err = fmt.Errorf("could not list shipments: %w", err)
		return
	}

	return models.Shipments{}.FromDatalayer(dlShipments), nil
}

func (bl *BusinessLogic) GetShipment(ctx context.Context, tenantID, shipmentID uuid.UUID) (_ models.Shipment, err error) {
	ctx, span := trace.Tracer().Start(ctx, "businesslogic.GetShipment")
	defer span.End()

	span.SetAttributes(
		attribute.String("tenant_id", tenantID.String()),
		attribute.String("shipment_id", shipmentID.String()),
	)

	dlShipment, err := bl.storage.GetShipment(ctx, tenantID.String(), shipmentID.String())
	if err != nil {
		err = fmt.Errorf("could not get shipment: %w", err)
		return
	}

	return models.Shipment{}.FromDatalayer(dlShipment), nil
}
