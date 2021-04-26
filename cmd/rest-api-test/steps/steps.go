package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"

	"github.com/lonnblad/shipment-service-backend/boundaries/rest/utils"
	v1 "github.com/lonnblad/shipment-service-backend/boundaries/rest/v1"
)

const tenantID = "fe131811-7fcd-4942-84a2-4ce8af359da5"

type sharedState struct {
	body []byte
}

func RegisterSteps(s *godog.ScenarioContext) {
	var state sharedState

	s.Step(`^price equation "([^"]*)"$`, priceEquation)
	s.Step(`^"([^"]*)" price rules$`, priceRules)
	s.Step(`^"([^"]*)" validation rules$`, validationRules)
	s.Step(`^a request to create a shipment with$`, state.aRequestToCreateAShipmentWith)
	s.Step(`^the returned shipment should have$`, state.theReturnedShipmentShouldHave)
	s.Step(`^the returned error should have$`, state.theReturnedErrorShouldHave)
}

func (state *sharedState) aRequestToCreateAShipmentWith(arg1 *godog.Table) error {
	createShipmentReq := newCreateShipmentRequest()

	for _, row := range arg1.Rows {
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
			weight, err := strconv.Atoi(value)
			if err != nil {
				return err
			}

			createShipmentReq.Package.Weight = weight
		default:
			return fmt.Errorf("unsupported key: %s", key)
		}
	}

	bs, err := json.Marshal(createShipmentReq)
	if err != nil {
		return err
	}

	const url = "http://localhost:8080/v1/tenants/" + tenantID + "/shipments"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if state.body, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	return nil
}

func priceEquation(arg1 string) error {
	return nil
}

func priceRules(arg1 string, arg2 *godog.DocString) error {
	return nil
}

func validationRules(arg1 string, arg2 *godog.DocString) error {
	return nil
}

func (state *sharedState) theReturnedShipmentShouldHave(arg1 *godog.Table) error {
	var createShipmentResp v1.CreateShipmentResponse

	err := json.Unmarshal(state.body, &createShipmentResp)
	if err != nil {
		return err
	}

	for _, row := range arg1.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "package - price":
			expectedPrice := value
			actualPrice := strconv.Itoa(createShipmentResp.Shipment.Package.Price.Amount)

			if expectedPrice != actualPrice {
				return fmt.Errorf("expected price: [%s] and actual price: [%s] are not equal", expectedPrice, actualPrice)
			}
		default:
			return fmt.Errorf("unsupported key: [%s]", key)
		}
	}

	return nil
}

func (state *sharedState) theReturnedErrorShouldHave(arg1 *godog.Table) error {
	var errResp utils.ErrorResponse

	err := json.Unmarshal(state.body, &errResp)
	if err != nil {
		return err
	}

	for _, row := range arg1.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "message":
			expectedMessage := value
			actualMessage := errResp.Error.Message

			if expectedMessage != actualMessage {
				return fmt.Errorf("expected message: [%s] and actual message: [%s] are not equal", expectedMessage, actualMessage)
			}
		default:
			return fmt.Errorf("unsupported key: [%s]", key)
		}
	}

	return nil
}
