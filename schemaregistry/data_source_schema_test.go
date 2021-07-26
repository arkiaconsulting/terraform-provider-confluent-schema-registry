package schemaregistry

import (
	"fmt"
	"github.com/riferrei/srclient"
	"os"
	"strconv"
	"strings"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSchema_basic(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	subject := fmt.Sprintf("sub%s", u)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fixtureDataSourceSchemaBuild(subject, fixtureAvro1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.schemaregistry_schema.test", "id", subject),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.test", "subject", subject),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.test", "version", "1"),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.test", "schema", strings.Replace(fixtureAvro1, "\\", "", -1)),
					resource.TestCheckResourceAttrSet("data.schemaregistry_schema.test", "schema_id"),
				),
			},
		},
	})
}

func TestAccDataSourceSchemaReferences_basic(t *testing.T) {
	// GIVEN
	url, found := os.LookupEnv("SCHEMA_REGISTRY_URL")
	if !found {
		t.Fatalf("SCHEMA_REGISTRY_URL must be set for acceptance tests")
	}
	username := os.Getenv("SCHEMA_REGISTRY_USERNAME")
	password := os.Getenv("SCHEMA_REGISTRY_PASSWORD")

	client := srclient.CreateSchemaRegistryClient(url)
	if (username != "") && (password != "") {
		client.SetCredentials(username, password)
	}

	// AND
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	referencedSchemaSubject := fmt.Sprintf("referencedSub-%s", u)
	referencedSchema := strings.Replace(fixtureAvro1, "\\", "", -1)

	schemaWithReferenceSubject := fmt.Sprintf("sub-%s", u)
	schemaWithReference := `["akc.test.userAdded"]`

	references := []srclient.Reference{
		{
			Name:    "akc.test.userAdded",
			Subject: referencedSchemaSubject,
			Version: 1,
		},
	}

	// AND
	if _, err = client.CreateSchemaWithArbitrarySubject(referencedSchemaSubject, referencedSchema, srclient.Avro); err != nil {
		t.Fatalf("could not create schema for subject: %s, err: %s", referencedSchema, err)
	}

	if _, err = client.CreateSchemaWithArbitrarySubject(schemaWithReferenceSubject, schemaWithReference, srclient.Avro, references...); err != nil {
		t.Fatalf("could not create schema for subject: %s, err: %s", referencedSchemaSubject, err)
	}

	// WHEN / THEN
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "schemaregistry_schema" "schemaWithReference" {
						subject = "%s"
					}
				`, schemaWithReferenceSubject),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "id", schemaWithReferenceSubject),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "subject", schemaWithReferenceSubject),
					resource.TestCheckResourceAttrSet("data.schemaregistry_schema.schemaWithReference", "schema_id"),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "version", "1"),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "schema", schemaWithReference),

					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "references.#", "1"),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "references.0.name", references[0].Name),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "references.0.subject", references[0].Subject),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaWithReference", "references.0.version", strconv.Itoa(references[0].Version)),
				),
			},
		},
	})
}

func TestAccDataSourceSchema_atVersion(t *testing.T) {
	// GIVEN
	url, found := os.LookupEnv("SCHEMA_REGISTRY_URL")
	if !found {
		t.Fatalf("SCHEMA_REGISTRY_URL must be set for acceptance tests")
	}
	username := os.Getenv("SCHEMA_REGISTRY_USERNAME")
	password := os.Getenv("SCHEMA_REGISTRY_PASSWORD")

	client := srclient.CreateSchemaRegistryClient(url)
	if (username != "") && (password != "") {
		client.SetCredentials(username, password)
	}

	// AND
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	referencedSchemaSubject := fmt.Sprintf("referencedSub-%s", u)
	referencedSchema := strings.Replace(fixtureAvro1, "\\", "", -1)
	referencedSchemaLatest := strings.Replace(fixtureAvro2, "\\", "", -1)

	// AND
	if _, err = client.CreateSchemaWithArbitrarySubject(referencedSchemaSubject, referencedSchema, srclient.Avro); err != nil {
		t.Fatalf("could not create schema for subject: %s, err: %s", referencedSchema, err)
	}

	if _, err = client.CreateSchemaWithArbitrarySubject(referencedSchemaSubject, referencedSchemaLatest, srclient.Avro); err != nil {
		t.Fatalf("could not create schema for subject: %s, err: %s", referencedSchema, err)
	}

	// WHEN / THEN
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "schemaregistry_schema" "schemaAtVersion" {
						subject = "%s"
						version = 1
					}
				`, referencedSchemaSubject),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaAtVersion", "id", referencedSchemaSubject),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaAtVersion", "subject", referencedSchemaSubject),
					resource.TestCheckResourceAttrSet("data.schemaregistry_schema.schemaAtVersion", "schema_id"),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaAtVersion", "version", "1"),
					resource.TestCheckResourceAttr("data.schemaregistry_schema.schemaAtVersion", "schema", referencedSchema),
				),
			},
		},
	})
}