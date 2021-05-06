package steps

import (
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
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

func decorateWithValues(createShipmentReq v1.CreateShipmentRequest, values *godog.Table) (_ v1.CreateShipmentRequest, err error) {
	for _, row := range values.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "sender - name":
			createShipmentReq.Sender.Name = value
		case "sender - email":
			createShipmentReq.Sender.Email = value
		case "sender - address":
			createShipmentReq.Sender.Address = value
		case "sender - country code":
			createShipmentReq.Sender.CountryCode = value
		case "receiver - name":
			createShipmentReq.Receiver.Name = value
		case "receiver - email":
			createShipmentReq.Receiver.Email = value
		case "receiver - address":
			createShipmentReq.Receiver.Address = value
		case "receiver - country code":
			createShipmentReq.Receiver.CountryCode = value
		case "package - weight":
			createShipmentReq.Package.Weight, err = strconv.Atoi(value)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("unsupported key: %s", key)
			return
		}
	}

	return createShipmentReq, nil
}
