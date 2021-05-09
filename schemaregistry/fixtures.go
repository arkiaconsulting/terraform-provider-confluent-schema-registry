package schemaregistry

import "fmt"

const fixtureAvro1 = `{\"type\":\"record\",\"name\":\"user\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"firstName\",\"type\":\"string\"}]}`
const fixtureAvro2 = `{\"type\":\"record\",\"name\":\"user\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"firstName\",\"type\":\"string\"},{\"name\":\"lastName\",\"type\":\"string\",\"default\":\"last\"}]}`

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

func fixtureDataSourceSchemaBuild(subject string, schema string) string {
	return fmt.Sprintf("%s%s", fmt.Sprintf(fixtureCreateSchema, subject, schema), fixtureDataSourceSchema)
}
