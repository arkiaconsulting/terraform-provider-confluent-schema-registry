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

func TestAccResourceSchemaReferences_basic(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		name    string
		fixture schemaWithReferenceFixture
		check   []resource.TestCheckFunc
	}{
		{
			name: "single reference",
			fixture: schemaWithReferenceFixture{
				Referenced: []SchemaResource{
					{
						ResourceName: "referencedSchema",
						Schema:       fixtureAvro1,
						Subject:      fmt.Sprintf("referencedSub-%s", u),
					},
				},
				WithReferences: SchemaResource{
					ResourceName: "schemaWithReference",
					Schema:       `[\"akc.test.userAdded\"]`,
					Subject:      fmt.Sprintf("sub%s", u),
				},
				References: []Reference{
					{
						Name:    "akc.test.userAdded",
						Subject: "schemaregistry_schema.referencedSchema.subject",
						Version: "schemaregistry_schema.referencedSchema.version",
					},
				},
			},
			check: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "id", fmt.Sprintf("referencedSub-%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "subject", fmt.Sprintf("referencedSub-%s", u)),
				resource.TestCheckResourceAttrSet("schemaregistry_schema.referencedSchema", "schema_id"),
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "version", "1"),
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "schema", strings.Replace(fixtureAvro1, "\\", "", -1)),

				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "id", fmt.Sprintf("sub%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "subject", fmt.Sprintf("sub%s", u)),
				resource.TestCheckResourceAttrSet("schemaregistry_schema.schemaWithReference", "schema_id"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "version", "1"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "schema", strings.Replace(`[\"akc.test.userAdded\"]`, "\\", "", -1)),

				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.#", "1"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.name", "akc.test.userAdded"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.subject", fmt.Sprintf("referencedSub-%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.version", "1"),
			},
		},
		{
			name: "multiple references",
			fixture: schemaWithReferenceFixture{
				Referenced: []SchemaResource{
					{
						ResourceName: "referencedSchema",
						Schema:       fixtureAvro1,
						Subject:      fmt.Sprintf("referencedSub-%s", u),
					},
					{
						ResourceName: "otherReferencedSchema",
						Schema:       `{\"type\":\"record\",\"name\":\"other\",\"namespace\":\"foo.bar\",\"fields\":[{\"name\":\"foo\",\"type\":\"string\"}]}`,
						Subject:      fmt.Sprintf("otherReferencedSub-%s", u),
					},
				},
				WithReferences: SchemaResource{
					ResourceName: "schemaWithReference",
					Schema:       `[\"akc.test.userAdded\"]`,
					Subject:      fmt.Sprintf("sub%s", u),
				},
				References: []Reference{
					{
						Name:    "akc.test.userAdded",
						Subject: "schemaregistry_schema.referencedSchema.subject",
						Version: "schemaregistry_schema.referencedSchema.version",
					},
					{
						Name:    "foo.bar.other",
						Subject: "schemaregistry_schema.otherReferencedSchema.subject",
						Version: "schemaregistry_schema.otherReferencedSchema.version",
					},
				},
			},
			check: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "id", fmt.Sprintf("referencedSub-%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "subject", fmt.Sprintf("referencedSub-%s", u)),
				resource.TestCheckResourceAttrSet("schemaregistry_schema.referencedSchema", "schema_id"),
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "version", "1"),
				resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "schema", strings.Replace(fixtureAvro1, "\\", "", -1)),

				resource.TestCheckResourceAttr("schemaregistry_schema.otherReferencedSchema", "id", fmt.Sprintf("otherReferencedSub-%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.otherReferencedSchema", "subject", fmt.Sprintf("otherReferencedSub-%s", u)),
				resource.TestCheckResourceAttrSet("schemaregistry_schema.otherReferencedSchema", "schema_id"),
				resource.TestCheckResourceAttr("schemaregistry_schema.otherReferencedSchema", "version", "1"),
				resource.TestCheckResourceAttr("schemaregistry_schema.otherReferencedSchema", "schema", strings.Replace(`{\"type\":\"record\",\"name\":\"other\",\"namespace\":\"foo.bar\",\"fields\":[{\"name\":\"foo\",\"type\":\"string\"}]}`, "\\", "", -1)),

				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "id", fmt.Sprintf("sub%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "subject", fmt.Sprintf("sub%s", u)),
				resource.TestCheckResourceAttrSet("schemaregistry_schema.schemaWithReference", "schema_id"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "version", "1"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "schema", strings.Replace(`[\"akc.test.userAdded\"]`, "\\", "", -1)),

				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.#", "2"),

				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.name", "akc.test.userAdded"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.subject", fmt.Sprintf("referencedSub-%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.version", "1"),

				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.1.name", "foo.bar.other"),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.1.subject", fmt.Sprintf("otherReferencedSub-%s", u)),
				resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.1.version", "1"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			config := fixtureResourceSchemaWithReferenceBuild(tc.fixture)
			resource.Test(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				PreCheck:          func() { testAccPreCheck(t) },
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: resource.ComposeTestCheckFunc(tc.check...),
					},
				},
			})

		})
	}
}

func TestAccResourceSchemaReferences_validateSchema(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	referencedSubject := fmt.Sprintf("referencedSub-%s", u)
	schemaWithReferenceSubject := fmt.Sprintf("sub-%s", u)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fixtureResourceSchemaWithReferenceBuild(schemaWithReferenceFixture{
					Referenced: []SchemaResource{
						{
							ResourceName: "referencedSchema",
							Schema:       fixtureAvro1,
							Subject:      referencedSubject,
						},
					},
					WithReferences: SchemaResource{
						ResourceName: "schemaWithReference",
						Schema:       `[\"invalid.schema\"]`,
						Subject:      schemaWithReferenceSubject,
					},
					References: []Reference{
						{
							Name:    "akc.test.userAdded",
							Subject: "schemaregistry_schema.referencedSchema.subject",
							Version: "schemaregistry_schema.referencedSchema.version",
						},
					},
				}),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(`Invalid schema ["invalid.schema"] with refs`)),
			},
		},
	})
}

func TestAccResourceSchemaReferences_updateCompatible(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	referencedSubject := fmt.Sprintf("referencedSub-%s", u)
	schemaWithReferenceSubject := fmt.Sprintf("sub-%s", u)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fixtureResourceSchemaWithReferenceBuild(schemaWithReferenceFixture{
					Referenced: []SchemaResource{
						{
							ResourceName: "referencedSchema",
							Schema:       fixtureAvro1,
							Subject:      referencedSubject,
						},
					},
					WithReferences: SchemaResource{
						ResourceName: "schemaWithReference",
						Schema:       `[\"akc.test.userAdded\"]`,
						Subject:      schemaWithReferenceSubject,
					},
					References: []Reference{
						{
							Name:    "akc.test.userAdded",
							Subject: "schemaregistry_schema.referencedSchema.subject",
							Version: "schemaregistry_schema.referencedSchema.version",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "id", referencedSubject),
					resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "version", "1"),

					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "id", schemaWithReferenceSubject),
					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "version", "1"),

					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.#", "1"),
					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.version", "1"),
				),
			},
			{
				Config: fixtureResourceSchemaWithReferenceBuild(schemaWithReferenceFixture{
					Referenced: []SchemaResource{
						{
							ResourceName: "referencedSchema",
							Schema:       fixtureAvro2,
							Subject:      referencedSubject,
						},
					},
					WithReferences: SchemaResource{
						ResourceName: "schemaWithReference",
						Schema:       `[\"akc.test.userAdded\"]`,
						Subject:      schemaWithReferenceSubject,
					},
					References: []Reference{
						{
							Name:    "akc.test.userAdded",
							Subject: "schemaregistry_schema.referencedSchema.subject",
							Version: "schemaregistry_schema.referencedSchema.version",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "id", referencedSubject),
					resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "version", "2"),

					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "id", schemaWithReferenceSubject),
					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "version", "2"),

					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.#", "1"),
					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.version", "2"),
				),
			},
		},
	})
}

func TestAccResourceSchemaReferences_updateIncompatible(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	referencedSubject := fmt.Sprintf("referencedSub-%s", u)
	schemaWithReferenceSubject := fmt.Sprintf("sub-%s", u)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fixtureResourceSchemaWithReferenceBuild(schemaWithReferenceFixture{
					Referenced: []SchemaResource{
						{
							ResourceName: "referencedSchema",
							Schema:       fixtureAvro1,
							Subject:      referencedSubject,
						},
					},
					WithReferences: SchemaResource{
						ResourceName: "schemaWithReference",
						Schema:       `[\"akc.test.userAdded\"]`,
						Subject:      schemaWithReferenceSubject,
					},
					References: []Reference{
						{
							Name:    "akc.test.userAdded",
							Subject: "schemaregistry_schema.referencedSchema.subject",
							Version: "schemaregistry_schema.referencedSchema.version",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "id", referencedSubject),
					resource.TestCheckResourceAttr("schemaregistry_schema.referencedSchema", "version", "1"),

					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "id", schemaWithReferenceSubject),
					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "version", "1"),

					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.#", "1"),
					resource.TestCheckResourceAttr("schemaregistry_schema.schemaWithReference", "references.0.version", "1"),
				),
			},
			{
				Config: fixtureResourceSchemaWithReferenceBuild(schemaWithReferenceFixture{
					Referenced: []SchemaResource{
						{
							ResourceName: "referencedSchema",
							Schema:       fixtureAvro1,
							Subject:      referencedSubject,
						},
					},
					WithReferences: SchemaResource{
						ResourceName: "schemaWithReference",
						Schema:       `[\"akc.test.incompatible\"]`,
						Subject:      schemaWithReferenceSubject,
					},
					References: []Reference{
						{
							Name:    "akc.test.userAdded",
							Subject: "schemaregistry_schema.referencedSchema.subject",
							Version: "schemaregistry_schema.referencedSchema.version",
						},
					},
				}),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(`Invalid schema ["akc.test.incompatible"] with refs`)),
			},
		},
	})
}

func TestAccResourceSchemaReferences_import(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	referencedSubject := fmt.Sprintf("referencedSub-%s", u)
	schemaWithReferenceSubject := fmt.Sprintf("sub-%s", u)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fixtureResourceSchemaWithReferenceBuild(schemaWithReferenceFixture{
					Referenced: []SchemaResource{
						{
							ResourceName: "referencedSchema",
							Schema:       fixtureAvro1,
							Subject:      referencedSubject,
						},
					},
					WithReferences: SchemaResource{
						ResourceName: "schemaWithReference",
						Schema:       `[\"akc.test.userAdded\"]`,
						Subject:      schemaWithReferenceSubject,
					},
					References: []Reference{
						{
							Name:    "akc.test.userAdded",
							Subject: "schemaregistry_schema.referencedSchema.subject",
							Version: "schemaregistry_schema.referencedSchema.version",
						},
					},
				}),
			},
			{
				ResourceName:      "schemaregistry_schema.referencedSchema",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "schemaregistry_schema.schemaWithReference",
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
