package price_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lonnblad/shipment-service-backend/businesslogic/models"
	"github.com/lonnblad/shipment-service-backend/businesslogic/price"
)

type testCase struct {
	name          string
	shipment      models.Shipment
	expectedPrice int
	expectedError error
}

func (tc testCase) Name() string {
	return tc.name
}

func Test_PriceCalculation(t *testing.T) {
	for _, tc := range createTestCases() {
		t.Run(tc.Name(), func(t *testing.T) {
			t.Parallel()

			actualPrice, actualError := price.Calculate(tc.shipment)

			assert.Equal(t, tc.expectedError, actualError)
			assert.Equal(t, tc.expectedPrice, actualPrice)
		})
	}
}

func createTestCases() (tcs []testCase) {
	tcs = append(tcs, newTestCase("Nordic/Small", "SE", 10, 100, nil))
	tcs = append(tcs, newTestCase("Nordic/Medium", "SE", 25, 300, nil))
	tcs = append(tcs, newTestCase("Nordic/Large", "SE", 50, 500, nil))
	tcs = append(tcs, newTestCase("Nordic/Huge", "SE", 1000, 2000, nil))

	tcs = append(tcs, newTestCase("EU/Small", "DE", 10, 150, nil))
	tcs = append(tcs, newTestCase("EU/Medium", "DE", 25, 450, nil))
	tcs = append(tcs, newTestCase("EU/Large", "DE", 50, 750, nil))
	tcs = append(tcs, newTestCase("EU/Huge", "DE", 1000, 3000, nil))

	tcs = append(tcs, newTestCase("Outside_EU/Small", "US", 10, 250, nil))
	tcs = append(tcs, newTestCase("Outside_EU/Medium", "US", 25, 750, nil))
	tcs = append(tcs, newTestCase("Outside_EU/Large", "US", 50, 1250, nil))
	tcs = append(tcs, newTestCase("Outside_EU/Huge", "US", 1000, 5000, nil))

	tcs = append(tcs, newTestCase("Bad_Weight", "SE", -1, 0, price.WeightClassError{Weight: -1}))
	tcs = append(tcs, newTestCase("Bad_CountryCode", "XX", 0, 0, price.CountryCodeError{CountryCode: "XX"}))

	return tcs
}

func newTestCase(name, senderCountryCode string, weight, expectedPrice int, expectedErr error) testCase {
	tc := testCase{}

	tc.name = name
	tc.shipment.Sender.CountryCode = senderCountryCode
	tc.shipment.Package.Weight = weight
	tc.expectedError = expectedErr
	tc.expectedPrice = expectedPrice

	return tc
}
