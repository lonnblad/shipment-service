package config

type Environment string

var (
	EnvironmentLocal   Environment = "local"
	EnvironmentSandbox Environment = "sandbox"
	EnvironmentStaging Environment = "staging"
	EnvironmentProd    Environment = "production"
)

func (s Environment) IsLocal() bool {
	return s == EnvironmentLocal
}

func (s Environment) IsStaging() bool {
	return s == EnvironmentStaging
}

func (s Environment) IsProd() bool {
	return s == EnvironmentProd
}

func (s Environment) String() string {
	return string(s)
}
