package v1

import (
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/lonnblad/shipment-service-backend/boundaries/rest/utils"
	"github.com/lonnblad/shipment-service-backend/businesslogic/models"
)

const (
	keyTenantID   = "tenant_id"
	keyShipmentID = "shipment_id"

	pathShipments = "/tenants/{" + keyTenantID + ":" + utils.RegexpUUID + "}/shipments"
	pathShipment  = pathShipments + "/{" + keyShipmentID + ":" + utils.RegexpUUID + "}"
)

type CreateShipmentRequest struct {
	Sender struct {
		Name        string `json:"name" example:"User Example A"`
		Email       string `json:"email" format:"email"`
		Address     string `json:"address" example:"Apt. Example 1A"`
		CountryCode string `json:"countryCode" example:"SE"`
	} `json:"sender"`

	Receiver struct {
		Name        string `json:"name" example:"User Example B"`
		Email       string `json:"email" format:"email"`
		Address     string `json:"address" example:"Apt. Example 1B"`
		CountryCode string `json:"countryCode" example:"DE"`
	} `json:"receiver"`

	Package struct {
		Weight int `json:"weight" example:"10"`
	} `json:"package"`
}

func (s CreateShipmentRequest) toInternal(tenantID uuid.UUID) models.Shipment {
	var internal models.Shipment

	internal.TenantID = tenantID

	internal.Sender = models.Sender(s.Sender)
	internal.Receiver = models.Receiver(s.Receiver)
	internal.Package.Weight = s.Package.Weight

	return internal
}

type CreateShipmentResponse getShipmentResponse

func (s CreateShipmentResponse) fromInternal(internal models.Shipment) (out CreateShipmentResponse) {
	out.Shipment = shipment{}.fromInternal(internal)
	return
}

func (s CreateShipmentResponse) decorateWithLinks(url url.URL) CreateShipmentResponse {
	s.Links = make([]link, 1)

	url.Path = "/v1/tenants/" + s.Shipment.TenantID.String() + "/shipments/" + s.Shipment.ID.String()
	s.Links[0] = link{Rel: "self", Href: url.String()}

	return s
}

type shipment struct {
	CreateShipmentRequest
	ID        uuid.UUID `json:"id" format:"uuid"`
	TenantID  uuid.UUID `json:"tenantId" format:"uuid"`
	CreatedAt time.Time `json:"createdAt" format:"date-time"`

	Package struct {
		Weight int      `json:"weight"`
		Price  currency `json:"price"`
	} `json:"package"`
}

type currency struct {
	Amount            int    `json:"amount"`
	DecimalMultiplier int    `json:"decimalMultiplier"`
	Currency          string `json:"string"`
}

func (s shipment) fromInternal(internal models.Shipment) shipment {
	s.ID = internal.ID
	s.TenantID = internal.TenantID
	s.CreatedAt = internal.CreatedAt

	s.Sender.Name = internal.Sender.Name
	s.Sender.Email = internal.Sender.Email
	s.Sender.Address = internal.Sender.Address
	s.Sender.CountryCode = internal.Sender.CountryCode

	s.Receiver.Name = internal.Receiver.Name
	s.Receiver.Email = internal.Receiver.Email
	s.Receiver.Address = internal.Receiver.Address
	s.Receiver.CountryCode = internal.Receiver.CountryCode

	s.Package.Weight = internal.Package.Weight
	s.Package.Price.Amount = internal.Package.Price
	s.Package.Price.DecimalMultiplier = 1
	s.Package.Price.Currency = "SEK"

	return s
}

type listShipmentsResponse struct {
	Shipments []getShipmentResponse `json:"shipments"`
	Metadata  struct {
		Total int `json:"total"`
	} `json:"metadata"`
	Links []link `json:"links"`
}

func (r listShipmentsResponse) fromInternal(shipments models.Shipments) listShipmentsResponse {
	r.Shipments = make([]getShipmentResponse, len(shipments))

	for idx, internal := range shipments {
		r.Shipments[idx] = getShipmentResponse{}.fromInternal(internal)
	}

	return r
}

func (r listShipmentsResponse) decorateWithLinks(url url.URL, req parsedListShipmentsRequest) listShipmentsResponse {
	r.Links = make([]link, 2)

	self := url
	self.Path = "/v1/tenants/" + req.tenantID.String() + "/shipments"
	selfQuery := self.Query()
	selfQuery.Add("limit", strconv.Itoa(req.limit))
	selfQuery.Add("offset", strconv.Itoa(req.offset))
	self.RawQuery = selfQuery.Encode()
	r.Links[0] = link{Rel: "self", Href: self.String()}

	next := url
	next.Path = "/v1/tenants/" + req.tenantID.String() + "/shipments"
	nextQuery := next.Query()
	nextQuery.Add("limit", strconv.Itoa(req.limit))
	nextQuery.Add("offset", strconv.Itoa(req.offset+len(r.Shipments)))
	next.RawQuery = nextQuery.Encode()
	r.Links[1] = link{Rel: "next", Href: next.String()}

	for idx := range r.Shipments {
		r.Shipments[idx] = r.Shipments[idx].decorateWithLinks(url)
	}

	return r
}

type getShipmentResponse struct {
	Shipment shipment `json:"shipment"`
	Links    []link   `json:"links"`
}

func (s getShipmentResponse) fromInternal(internal models.Shipment) (out getShipmentResponse) {
	out.Shipment = shipment{}.fromInternal(internal)
	return
}

func (s getShipmentResponse) decorateWithLinks(url url.URL) getShipmentResponse {
	s.Links = make([]link, 1)

	url.Path = "/v1/tenants/" + s.Shipment.TenantID.String() + "/shipments/" + s.Shipment.ID.String()
	s.Links[0] = link{Rel: "self", Href: url.String()}

	return s
}

type link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
