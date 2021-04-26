package models

import (
	"fmt"
	"regexp"

	"github.com/badoux/checkmail"
	"github.com/pariz/gountries"
)

const (
	lengthCountryCodeAlpha2 = 2
	maxLengthName           = 30
	maxLengthAddress        = 100
	minPackageWeight        = 0
	maxPackageWeight        = 1000
)

var (
	regexpNameValidatorNumbers = regexp.MustCompile("[0-9]+")
	countries                  = gountries.New()
)

// Validate ...
func (s Shipment) Validate() error {
	if err := s.Sender.validate(); err != nil {
		return fmt.Errorf("failed to validate sender: %w", err)
	}

	if err := s.Receiver.validate(); err != nil {
		return fmt.Errorf("failed to validate receiver: %w", err)
	}

	if err := s.Package.validate(); err != nil {
		return fmt.Errorf("failed to validate package: %w", err)
	}

	return nil
}

func (s Sender) validate() error {
	if err := validateName(s.Name); err != nil {
		return fmt.Errorf("sender name is invalid: %w", err)
	}

	if err := checkmail.ValidateFormat(s.Email); err != nil {
		return fmt.Errorf("sender email: %s is not valid: %w", s.Email, err)
	}

	if len(s.Address) > maxLengthAddress {
		return fmt.Errorf("sender address: %s is longer than max length: %d", s.Address, maxLengthAddress)
	}

	if err := validateCountry(s.CountryCode); err != nil {
		return fmt.Errorf("sender country code is invalid: %w", err)
	}

	return nil
}

func (r Receiver) validate() error {
	if err := validateName(r.Name); err != nil {
		return fmt.Errorf("receiver name is invalid: %w", err)
	}

	if err := checkmail.ValidateFormat(r.Email); err != nil {
		return fmt.Errorf("receiver email: %s is not valid: %w", r.Email, err)
	}

	if len(r.Address) > maxLengthAddress {
		return fmt.Errorf("receiver address: %s is longer than max length: %d", r.Address, maxLengthAddress)
	}

	if err := validateCountry(r.CountryCode); err != nil {
		return fmt.Errorf("receiver country code is invalid: %w", err)
	}

	return nil
}

func (p Package) validate() error {
	if p.Weight < minPackageWeight {
		return fmt.Errorf("package weight: %d can't be below minimum: %d", p.Weight, minPackageWeight)
	}

	if p.Weight > maxPackageWeight {
		return fmt.Errorf("package weight: %d can't be above maximum: %d", p.Weight, maxPackageWeight)
	}

	return nil
}

func validateName(name string) error {
	if len(name) > maxLengthName {
		return fmt.Errorf("name: %s is longer than the max length: %d", name, maxLengthName)
	}

	if regexpNameValidatorNumbers.MatchString(name) {
		return fmt.Errorf("name is not valid, it contains numbers: %s", name)
	}

	return nil
}

func validateCountry(countryCode string) error {
	if len(countryCode) > lengthCountryCodeAlpha2 {
		return fmt.Errorf("country code: %s is not of length: %d", countryCode, lengthCountryCodeAlpha2)
	}

	if _, err := countries.FindCountryByAlpha(countryCode); err != nil {
		return fmt.Errorf("could not find country by code: %s, error: %w", countryCode, err)
	}

	return nil
}
