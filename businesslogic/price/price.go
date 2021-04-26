package price

import (
	"fmt"
	"strings"

	"github.com/pariz/gountries"

	"github.com/lonnblad/shipment-service-backend/businesslogic/models"
)

// Calculate will return a price or an error if it didn't succeed in
// calculating a price.
//
// Note. the price is returned as an integer representing the real value,
// as the base prices are all multiplies of 10 and the specified region
// multiplier only use one decimal, this is fine since the result will
// always be an integer.
//
// In a real world scenario, where there might be a need for prices with
// decimals, the returned integer should compensate so that a value like
// 100.05 is returned as 10005, with a property stating the value 100 as
// the decimal multipler.
//
// This article talks about why floating points shouldn't be used for currency.
// https://husobee.github.io/money/float/2016/09/23/never-use-floats-for-currency.html
func Calculate(s models.Shipment) (_ int, err error) {
	basePrice, err := findBasePrice(s.Package.Weight)
	if err != nil {
		return
	}

	multiplier, err := findRegionMultiplier(s.Sender.CountryCode)
	if err != nil {
		return
	}

	return basePrice * multiplier / regionMulitplierAdjustment, nil
}

type (
	// WeightError will be returned by Calculate when there isn't
	// a defined weight class for the provided weight.
	WeightClassError struct{ Weight int }

	// CountryCodeError will be returned by Calculate when there
	// isn't a defined a country for the provided country code.
	CountryCodeError struct{ CountryCode string }
)

func (we WeightClassError) Error() string {
	return fmt.Sprintf("weight: %d don't have a defined price", we.Weight)
}

func (cce CountryCodeError) Error() string {
	return fmt.Sprintf("countryCode: %s is not defined", cce.CountryCode)
}

const (
	// - Small (0 - 10kg): 100sek
	basePriceSmall        = 100
	weightLowerBoundSmall = 0
	weightUpperBoundSmall = 10

	// - Medium (10 - 25kg): 300sek
	basePriceMedium        = 300
	weightUpperBoundMedium = 25

	// - Large (26 - 50kg): 500sek
	basePriceLarge        = 500
	weightUpperBoundLarge = 50

	// - Huge (51 - 1000kg): 2000sek
	basePriceHuge        = 2000
	weightUpperBoundHuge = 1000

	// All the region multipliers are multiplied by 10
	// to remove the need of using a floating pointer.
	// - Nordic region (Sweden, Norway, Denmark, Finland), the price is multiplied by 1
	// - EU region, the price is multiplied by 1.5
	// - Outside the EU, the price is multiplied by 2.5
	regionMulitplierAdjustment = 10
	regionMultiplierNordic     = 10
	regionMultiplierEU         = 15
	regionMultiplierNonEU      = 25
)

func findBasePrice(weight int) (_ int, err error) {
	var price int

	switch {
	case weightLowerBoundSmall <= weight && weight <= weightUpperBoundSmall:
		price = basePriceSmall
	case weightUpperBoundSmall < weight && weight <= weightUpperBoundMedium:
		price = basePriceMedium
	case weightUpperBoundMedium < weight && weight <= weightUpperBoundLarge:
		price = basePriceLarge
	case weightUpperBoundLarge < weight && weight <= weightUpperBoundHuge:
		price = basePriceHuge
	default:
		err = WeightClassError{Weight: weight}
		return
	}

	return price, nil
}

var countries = gountries.New()

func findRegionMultiplier(countryCode string) (_ int, err error) {
	country, err := countries.FindCountryByAlpha(countryCode)
	if err != nil {
		err = CountryCodeError{CountryCode: countryCode}
		return
	}

	// Nordic Region
	switch strings.ToLower(country.Alpha2) {
	case "se", "no", "dk", "fi":
		return regionMultiplierNordic, nil
	}

	// EU Region
	if country.EuMember {
		return regionMultiplierEU, nil
	}

	// Outside the EU
	return regionMultiplierNonEU, nil
}
