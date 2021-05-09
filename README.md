# terraform-provider-confluent-schema-registry
A terraform provider for managing schemas in a Confluent schema registry

## Provider configuration
```
terraform {
    required_providers {
        schemaregistry = {
            source = "github.com/arkiaconsulting/schemaregistry"
        }
    }
}

provider "schemaregistry" {
    schema_registry_url = "https://xxxxxx.confluent.cloud"
    username            = "<confluent_schema_registry_key>"
    password            = "<confluent_schema_registry_password>"
}
```
_You can omit the credential details by defining the environment variables `SCHEMA_REGISTRY_URL`, `SCHEMA_REGISTRY_USERNAME`, `SCHEMA_REGISTRY_PASSWORD`_

## The schema data source
```
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
```