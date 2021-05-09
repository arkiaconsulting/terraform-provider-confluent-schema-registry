terraform {
  required_providers {
    schemaregistry = {
      version = "0.1"
      source  = "github.com/arkiaconsulting/schemaregistry"
    }
  }
}

provider "schemaregistry" {
}

data "schemaregistry_schema" "main" {
  subject = "Akc-key"
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
