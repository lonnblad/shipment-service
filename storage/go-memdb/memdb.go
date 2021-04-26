package memdb

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-memdb"
	"go.opentelemetry.io/otel/attribute"

	"github.com/lonnblad/shipment-service-backend/storage"
	"github.com/lonnblad/shipment-service-backend/trace"
)

var _ storage.ShipmentStorage = &ShipmentStorage{}

const (
	writeMode = true
	readMode  = false

	tableShipments                   = "shipment"
	tableShipmentsIndexKeyTenant     = "tenant"
	tableShipmentsIndexFieldTenant   = "TenantID"
	tableShipmentsIndexKeyShipment   = "id"
	tableShipmentsIndexFieldShipment = "ID"
)

// Create the DB schema
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		tableShipments: {
			Name: tableShipments,
			Indexes: map[string]*memdb.IndexSchema{
				tableShipmentsIndexKeyShipment: {
					Name:   tableShipmentsIndexKeyShipment,
					Unique: true,
					Indexer: &memdb.CompoundIndex{
						Indexes: []memdb.Indexer{
							&memdb.UUIDFieldIndex{Field: tableShipmentsIndexFieldTenant},
							&memdb.UUIDFieldIndex{Field: tableShipmentsIndexFieldShipment},
						},
					},
				},
				tableShipmentsIndexKeyTenant: {
					Name:    tableShipmentsIndexKeyTenant,
					Unique:  false,
					Indexer: &memdb.UUIDFieldIndex{Field: tableShipmentsIndexFieldTenant},
				},
			},
		},
	},
}

// ShipmentStorage implements storage.ShipmentStorage
type ShipmentStorage struct {
	db *memdb.MemDB
}

// NewShipmentStorage will return a pointer to a new in-mem ShipmentStorage
func NewShipmentStorage() (_ *ShipmentStorage, err error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		err = fmt.Errorf("failed to create a new memdb: %w", err)
		return
	}

	return &ShipmentStorage{db: db}, nil
}

func (s *ShipmentStorage) StoreShipment(ctx context.Context, shipment storage.Shipment) error {
	_, span := trace.Tracer().Start(ctx, "memdb.StoreShipment")
	defer span.End()

	span.SetAttributes(
		attribute.String("shipment.tenant_id", shipment.TenantID),
		attribute.String("shipment.id", shipment.ID),
	)

	txn := s.db.Txn(writeMode)
	defer txn.Commit()

	err := txn.Insert(tableShipments, shipment)
	if err != nil {
		txn.Abort()
		return fmt.Errorf("failed to insert shipment: %w", err)
	}

	return nil
}

func (s *ShipmentStorage) GetShipment(ctx context.Context, tenantID, shipmentID string) (_ storage.Shipment, err error) {
	_, span := trace.Tracer().Start(ctx, "memdb.GetShipment")
	defer span.End()

	span.SetAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.String("shipment_id", shipmentID),
	)

	txn := s.db.Txn(readMode)

	obj, err := txn.First(tableShipments, tableShipmentsIndexKeyShipment, tenantID, shipmentID)
	if err != nil {
		err = fmt.Errorf("could not look up shipment: %w", err)
		return
	}

	if obj == nil {
		err = fmt.Errorf("could not find up shipment: %w", err)
		return
	}

	return obj.(storage.Shipment), nil
}

func (s *ShipmentStorage) ListShipments(ctx context.Context, tenantID string, limit, offset int) (_ []storage.Shipment, err error) {
	_, span := trace.Tracer().Start(ctx, "memdb.ListShipments")
	defer span.End()

	span.SetAttributes(
		attribute.String("tenant_id", tenantID),
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	)

	txn := s.db.Txn(readMode)

	it, err := txn.Get(tableShipments, tableShipmentsIndexKeyTenant, tenantID)
	if err != nil {
		err = fmt.Errorf("could not look up shipments: %w", err)
		return
	}

	shipments := make([]storage.Shipment, 0, limit)

	var limitCounter = 0
	if limitCounter == limit {
		return shipments, nil
	}

	var offsetCounter = 0

	for obj := it.Next(); obj != nil; obj = it.Next() {
		if offsetCounter++; offsetCounter <= offset {
			continue
		}

		shipment := obj.(storage.Shipment)
		shipments = append(shipments, shipment)

		if limitCounter++; limitCounter == limit {
			break
		}
	}

	return shipments, nil
}
