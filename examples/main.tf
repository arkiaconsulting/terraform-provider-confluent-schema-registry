terraform {
  required_providers {
    schemaregistry = {
      version = "0.3"
      source  = "github.com/arkiaconsulting/schemaregistry"
    }
  }
}

provider "schemaregistry" {
}

resource "schemaregistry_schema" "user_added" {
  subject = "MyTopic-akc.test.userAdded"
  schema  = file("./userAdded.avsc")
}

data "schemaregistry_schema" "main" {
  subject = schemaregistry_schema.user_added.subject
}

output "schema_id" {
  value = data.schemaregistry_schema.main.id
}

output "schema_version" {
  value = data.schemaregistry_schema.main.version
}

output "schema_string" {
  value = data.schemaregistry_schema.main.schema
}

output "schema_schema_id" {
  value = data.schemaregistry_schema.main.schema_id
}
