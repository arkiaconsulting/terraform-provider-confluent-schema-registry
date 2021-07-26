terraform {
  required_providers {
    schemaregistry = {
      version = "0.6"
      source  = "github.com/arkiaconsulting/schemaregistry"
    }
  }
}

provider "schemaregistry" {
  schema_registry_url = "http://localhost:8081"
}

resource "schemaregistry_schema" "user_added" {
  subject = "MyTopic-akc.test.userAdded-value"
  schema  = file("./userAdded.avsc")
}

resource "schemaregistry_schema" "with_reference" {
  subject = "with-reference"
  schema = "[\"akc.test.userAdded\"]"

  reference {
    name = "akc.test.userAdded"
    subject = schemaregistry_schema.user_added.subject
    version = schemaregistry_schema.user_added.version
  }
}

data "schemaregistry_schema" "main" {
  subject = schemaregistry_schema.with_reference.subject
}

data "schemaregistry_schema" "user_added_v1" {
  subject = "MyTopic-akc.test.userAdded-value"
  schema  = 1
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

output "schema_references" {
  value = data.schemaregistry_schema.main.references
}
