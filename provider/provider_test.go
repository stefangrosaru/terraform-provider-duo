package provider

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ProviderFactories = map[string]func() (*schema.Provider, error){
	"duo": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestAccPreCheck(t *testing.T) {
	err := accPreCheck()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func accPreCheck() error {
	if v := os.Getenv("DUO_API_HOSTNAME"); v == "" {
		return errors.New("DUO_API_HOSTNAME must be set for acceptance tests")
	}
	if v := os.Getenv("DUO_INTEGRATION_KEY"); v == "" {
		return errors.New("DUO_INTEGRATION_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("DUO_SECRET_KEY"); v == "" {
		return errors.New("DUO_SECRET_KEY must be set for acceptance tests")
	}

	return nil
}
