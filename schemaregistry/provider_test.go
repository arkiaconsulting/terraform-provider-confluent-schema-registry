package schemaregistry

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider = getProvider()
var testAccProviders = testAccProvidersFactory(testAccProvider)

func testAccProvidersFactory(provider *schema.Provider) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"schemaregistry": func() (*schema.Provider, error) {
			return provider, nil
		},
	}
}

func getProvider() *schema.Provider {
	return Provider()
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	log.Println("[INFO] TestProvider_impl")
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	log.Println("[INFO] testAccPreCheck")

	if v := os.Getenv("SCHEMA_REGISTRY_URL"); v == "" {
		t.Fatal("SCHEMA_REGISTRY_URL must be set for acceptance tests")
	}

	if v := os.Getenv("SCHEMA_REGISTRY_USERNAME"); v == "" {
		t.Fatal("SCHEMA_REGISTRY_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("SCHEMA_REGISTRY_PASSWORD"); v == "" {
		t.Fatal("SCHEMA_REGISTRY_PASSWORD must be set for acceptance tests")
	}
}
