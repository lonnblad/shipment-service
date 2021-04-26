package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/lonnblad/shipment-service-backend/storage"
)

type Shipments []Shipment

type Shipment struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	CreatedAt time.Time

	Sender   Sender
	Receiver Receiver
	Package  Package
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

func (s Shipment) ToDatalayer() (dlShipment storage.Shipment) {
	dlShipment.ID = s.ID.String()
	dlShipment.TenantID = s.TenantID.String()
	dlShipment.CreatedAt = s.CreatedAt

	dlShipment.Sender = storage.Sender(s.Sender)
	dlShipment.Receiver = storage.Receiver(s.Receiver)
	dlShipment.Package = storage.Package(s.Package)

	return
}

func (s Shipment) FromDatalayer(dlShipment storage.Shipment) Shipment {
	s.ID = uuid.MustParse(dlShipment.ID)
	s.TenantID = uuid.MustParse(dlShipment.TenantID)
	s.CreatedAt = dlShipment.CreatedAt

	s.Sender = Sender(dlShipment.Sender)
	s.Receiver = Receiver(dlShipment.Receiver)
	s.Package = Package(dlShipment.Package)

	return s
}

func (s Shipments) FromDatalayer(dlShipments []storage.Shipment) Shipments {
	s = make(Shipments, len(dlShipments))

	for idx := range dlShipments {
		s[idx] = Shipment{}.FromDatalayer(dlShipments[idx])
	}

	return s
}
