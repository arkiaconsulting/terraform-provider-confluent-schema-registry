package schemaregistry

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/riferrei/srclient"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"schema_registry_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCHEMA_REGISTRY_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCHEMA_REGISTRY_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCHEMA_REGISTRY_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"schemaregistry_schema": resourceSchema(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"schemaregistry_schema": dataSourceSchema(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("schema_registry_url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (url != "") {
		client := srclient.CreateSchemaRegistryClient(url)

		if (username != "") && (password != "") {		
			client.SetCredentials(username, password)	
		}
		
		return client, diags
	}

	return nil, diag.FromErr(errors.New("invalid credential parameters"))
}
