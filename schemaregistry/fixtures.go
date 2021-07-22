package schemaregistry

import (
	"bytes"
	"fmt"
	"text/template"
)

const fixtureAvro1 = `{\"type\":\"record\",\"name\":\"userAdded\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"firstName\",\"type\":\"string\"}]}`
const fixtureAvro2 = `{\"type\":\"record\",\"name\":\"userAdded\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"firstName\",\"type\":\"string\"},{\"name\":\"lastName\",\"type\":\"string\",\"default\":\"last\"}]}`
const fixtureAvro3 = `{\"type\":\"record\",\"name\":\"userAdded\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"firstName\",\"type\":\"string\"},{\"name\":\"lastName\",\"type\":\"string\"}]}`

const fixtureCreateSchema = `
	resource "schemaregistry_schema" "test" {
		subject = "%s"
		schema = "%s"
	}
`

const fixtureDataSourceSchema = `
	data "schemaregistry_schema" "test" {
		subject = schemaregistry_schema.test.subject
	}
`
const fixtureImportSchema = `
	resource "schemaregistry_schema" "import" {
		subject = "%s"
		schema = "%s"
	}
`

func fixtureDataSourceSchemaBuild(subject string, schema string) string {
	return fmt.Sprintf("%s%s", fmt.Sprintf(fixtureCreateSchema, subject, schema), fixtureDataSourceSchema)
}

type SchemaResource struct {
	ResourceName string
	Subject      string
	Schema       string
}

type SchemaWithReferences struct {
	SchemaResource
	ReferenceName string
}

type Reference struct {
	Name    string
	Subject string
	Version string
}

type schemaWithReferenceFixture struct {
	Type           string
	Referenced     []SchemaResource
	WithReferences SchemaResource
	References     []Reference
}

func fixtureResourceSchemaWithReferenceBuild(s schemaWithReferenceFixture) string {
	var buf bytes.Buffer

	err := template.Must(template.New("SchemaWithReferenceFixture").Parse(`
    {{range .Referenced}}
	resource "schemaregistry_schema" "{{.ResourceName}}" {
		subject = "{{.Subject}}"
		schema = "{{.Schema}}"
	}
    {{end}}

	resource "schemaregistry_schema" "{{.WithReferences.ResourceName}}" {
		subject = "{{.WithReferences.Subject}}"
		schema = "{{.WithReferences.Schema}}"
		{{range .References}}

		references {
			name = "{{.Name}}"
			subject = {{.Subject}}
			version = {{.Version}}
		}
		{{end}}
	}
	`)).Execute(&buf, s)

	if err != nil {
		panic(err)
	}

	return buf.String()
}
