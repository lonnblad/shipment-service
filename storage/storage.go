package storage

import (
	"context"
	"time"
)

// ShipmentStorage is an interface for managing storage of shipments
type ShipmentStorage interface {
	StoreShipment(context.Context, Shipment) error
	GetShipment(_ context.Context, tenantID, shipmentID string) (Shipment, error)
	ListShipments(_ context.Context, tenantID string, limit, offset int) ([]Shipment, error)
}

type Shipment struct {
	ID        string
	TenantID  string
	CreatedAt time.Time
	Sender    Sender
	Receiver  Receiver
	Package   Package
}

type Sender struct {
	Name        string
	Email       string
	Address     string
	CountryCode string
}

type Receiver struct {
	Name        string
	Email       string
	Address     string
	CountryCode string
}

type Package struct {
	Weight int
	Price  int
}
