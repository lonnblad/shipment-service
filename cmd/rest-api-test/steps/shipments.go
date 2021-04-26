package steps

import (
	v1 "github.com/lonnblad/shipment-service-backend/boundaries/rest/v1"
)

func newCreateShipmentRequest() v1.CreateShipmentRequest {
	var req v1.CreateShipmentRequest

	req.Sender.Name = "User Example A"
	req.Sender.Email = "user@example.com"
	req.Sender.Address = "Apt. Example 1A"
	req.Sender.CountryCode = "SE"

	req.Receiver.Name = "User Example B"
	req.Receiver.Email = "user@example.com"
	req.Receiver.Address = "Apt. Example 1B"
	req.Receiver.CountryCode = "DE"

	req.Package.Weight = 10

	return req
}
