package main_test

import (
	"github.com/cucumber/godog"

	"github.com/lonnblad/shipment-service-backend/cmd/rest-api-test/steps"
)

func FeatureContext(s *godog.ScenarioContext) {
	steps.RegisterSteps(s)
}
