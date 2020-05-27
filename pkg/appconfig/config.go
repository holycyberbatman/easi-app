package appconfig

import "fmt"

// NewEnvironment returns an environment from a string
func NewEnvironment(config string) (Environment, error) {
	switch config {
	case LocalEnv.String():
		return LocalEnv, nil
	case TestEnv.String():
		return TestEnv, nil
	case DevEnv.String():
		return DevEnv, nil
	case ImplEnv.String():
		return ImplEnv, nil
	case ProdEnv.String():
		return ProdEnv, nil
	default:
		return "", fmt.Errorf("unknown environment: %s", config)
	}
}

// EnvironmentKey is used to access the environment from a config
const EnvironmentKey = "APP_ENV"

// Environment represents an environment
type Environment string

const (
	// LocalEnv is the local environment
	LocalEnv Environment = "local"
	// TestEnv is the environment for running tests
	TestEnv Environment = "test"
	// DevEnv is the environment for the dev deployed env
	DevEnv Environment = "dev"
	// ImplEnv is the environment for the impl deployed env
	ImplEnv Environment = "impl"
	// ProdEnv is the environment for the impl deployed env
	ProdEnv Environment = "prod"
)

// String gets the environment as a string
func (e Environment) String() string {
	switch e {
	case LocalEnv:
		return "local"
	case TestEnv:
		return "test"
	case DevEnv:
		return "dev"
	case ImplEnv:
		return "impl"
	case ProdEnv:
		return "prod"
	default:
		return ""
	}
}

// Local returns true if the environment is local
func (e Environment) Local() bool {
	if e == LocalEnv {
		return true
	}
	return false
}

// Test returns true if the environment is local
func (e Environment) Test() bool {
	if e == TestEnv {
		return true
	}
	return false
}

// Dev returns true if the environment is local
func (e Environment) Dev() bool {
	if e == TestEnv {
		return true
	}
	return false
}

// Impl returns true if the environment is local
func (e Environment) Impl() bool {
	if e == ImplEnv {
		return true
	}
	return false
}

// Prod returns true if the environment is local
func (e Environment) Prod() bool {
	if e == ProdEnv {
		return true
	}
	return false
}

// DBHostConfigKey is the Postgres hostname config key
const DBHostConfigKey = "PGHOST"

// DBPortConfigKey is the Postgres port config key
const DBPortConfigKey = "PGPORT"

// DBNameConfigKey is the Postgres database name config key
const DBNameConfigKey = "PGDATABASE"

// DBUsernameConfigKey is the Postgres username config key
const DBUsernameConfigKey = "PGUSER"

// DBPasswordConfigKey is the Postgres password config key
const DBPasswordConfigKey = "PGPASS"

// DBSSLModeConfigKey is the Postgres SSL mode config key
const DBSSLModeConfigKey = "PGSSLMODE"

// AWSSESSourceARNKey is the key for the ARN for sending email
const AWSSESSourceARNKey = "AWS_SES_SOURCE_ARN"

// AWSSESSourceKey is the key for the sender for sending email
const AWSSESSourceKey = "AWS_SES_SOURCE"

// GRTEmailKey is the key for the receiving email for the GRT
const GRTEmailKey = "GRT_EMAIL"

// ClientHostKey is the key for getting the client's domain name
const ClientHostKey = "CLIENT_HOSTNAME"

// ClientProtocolKey is the key for getting the client's protocol
const ClientProtocolKey = "CLIENT_PROTOCOL"

// EmailTemplateDirectoryKey is the key for getting the email template directory
const EmailTemplateDirectoryKey = "EMAIL_TEMPLATE_DIR"
