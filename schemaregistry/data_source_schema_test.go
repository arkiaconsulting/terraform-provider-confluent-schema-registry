package schemaregistry

import (
	"fmt"
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
