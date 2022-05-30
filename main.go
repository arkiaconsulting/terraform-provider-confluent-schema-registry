package main

import (
	"terraform-provider-confluent-schema-registry/schemaregistry"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return schemaregistry.Provider()
		},
	})
}
