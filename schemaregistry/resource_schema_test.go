package schemaregistry

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSchema_basic(t *testing.T) {
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
				Config: fmt.Sprintf(fixtureCreateSchema, subject, fixtureAvro1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "id", subject),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "subject", subject),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "schema_id"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "version", "1"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema", strings.Replace(fixtureAvro1, "\\", "", -1)),
				),
			},
		},
	})
}

func TestAccResourceSchema_updateCompatible(t *testing.T) {
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
				Config: fmt.Sprintf(fixtureCreateSchema, subject, fixtureAvro1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "id", subject),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "subject", subject),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "schema_id"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "version", "1"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema", strings.Replace(fixtureAvro1, "\\", "", -1)),
				),
			},
			{
				Config: fmt.Sprintf(fixtureCreateSchema, subject, fixtureAvro2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "id", subject),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "subject", subject),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "schema_id"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "version", "2"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema", strings.Replace(fixtureAvro2, "\\", "", -1)),
				),
			},
		},
	})
}

func TestAccResourceSchema_updateIncompatible(t *testing.T) {
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
				Config: fmt.Sprintf(fixtureCreateSchema, subject, fixtureAvro1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "id", subject),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "subject", subject),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "schema_id"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "version", "1"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema", strings.Replace(fixtureAvro1, "\\", "", -1)),
				),
			},
			{
				Config:      fmt.Sprintf(fixtureCreateSchema, subject, fixtureAvro3),
				ExpectError: regexp.MustCompile(`invalid "schema": incompatible`),
			},
		},
	})
}

func TestAccResourceSchema_import(t *testing.T) {
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
				Config: fmt.Sprintf(fixtureImportSchema, subject, fixtureAvro1),
			},
			{
				ResourceName:      "schemaregistry_schema.import",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := "to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}
