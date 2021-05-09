package schemaregistry

import (
	"fmt"
)

const IDSeparator = "___"

func formatSchemaVersionID(subject string) string {
	return fmt.Sprintf("%s", subject)
}

func extractSchemaVersionID(id string) string {
	return id
}
