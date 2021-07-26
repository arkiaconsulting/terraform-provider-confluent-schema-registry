# terraform-provider-confluent-schema-registry
A terraform provider for managing schemas in a Confluent schema registry

## Provider configuration
```
terraform {
    required_providers {
        schemaregistry = {
            source = "arkiaconsulting/confluent-schema-registry"
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

## The schema resource
```
resource "schemaregistry_schema" "main" {
  subject = "<subject_name>"
  schema  = file("<avro_schema_file>")
}
```

## The schema resource with references

Schema registry references can be used to allow [putting Several Event Types in the Same Topic](https://www.confluent.io/blog/multiple-event-types-in-the-same-kafka-topic/).
Please refer to Confluent [Schema Registry API Reference](https://docs.confluent.io/platform/current/schema-registry/develop/api.html) for details.

### Upgrade reference version to the latest event schema version

Reference the event schema `resource` from a schema with reference wil upgrade a reference alongside with its referenced schema.

```
resource "schemaregistry_schema" "referenced_event" {
  subject = "referenced_event_subject"
  schema  = "{\"type\":\"record\",\"name\":\"event\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"foo\",\"type\":\"string\"}]}"
}

resource "schemaregistry_schema" "other_referenced_event" {
  subject = "other_referenced_event_subject"
  schema  = "{\"type\":\"record\",\"name\":\"other_event\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"bar\",\"type\":\"string\"}]}"
}

resource "schemaregistry_schema" "with_reference" {
  subject = "with_reference_subject"
  schema = "[\"akc.test.event\", \"akc.test.other_event\"]"

  reference {
    name = "akc.test.event"
    subject = schemaregistry_schema.referenced_event.subject
    // version will always be upgraded with the referenced event schema version  
    version = schemaregistry_schema.referenced_event.version
  }

  reference {
    name = "akc.test.other_event"
    subject = schemaregistry_schema.user_added.subject
    // version will always be upgraded with the referenced event schema version  
    version = schemaregistry_schema.referenced_event.version
  }
}
```

### Stick reference version to a given version

Use a `dataSource` to stick a reference to a **given version**, while upgrading the referenced event schema.

```
resource "schemaregistry_schema" "referenced_event_latest" {
  subject = "referenced_event_subject"
  schema  = file("<avro_schema_file_updated>")
}

data "schemaregistry_schema" "referenced_event_v1" {
  subject = "other_referenced_event_subject"
  schema  = "{\"type\":\"record\",\"name\":\"other_event\",\"namespace\":\"akc.test\",\"fields\":[{\"name\":\"bar\",\"type\":\"string\"}]}"
}

resource "schemaregistry_schema" "with_reference_to_v1" {
  subject = "with_reference_subject"
  schema = "[\"akc.test.event\", \"akc.test.other_event\"]"

  references {
    name = "akc.test.event"
    subject = data.schemaregistry_schema.referenced_event_v1.subject
    version = data.schemaregistry_schema.referenced_event_v1.version
  }
}
```

## The schema data source
```
data "schemaregistry_schema" "main" {
  subject = "<subject_name>"
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

## Importing an existing schema
`
terraform import schemaregistry_schema.main <subject_name>
`
