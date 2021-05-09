package schemaregistry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/riferrei/srclient"
)

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubjectRead,
		Schema: map[string]*schema.Schema{
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subject related to the schema",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the schema",
			},
			"schema_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The schema ID",
			},
			"schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The schema string",
			},
		},
	}
}

func dataSourceSubjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	subject := d.Get("subject").(string)

	client := m.(*srclient.SchemaRegistryClient)

	schema, err := client.GetLatestSchemaWithArbitrarySubject(subject)
	if err != nil {
		return diag.FromErr(err)
		// return diag.FromErr(fmt.Errorf("unknown schema for subject '%s'", subject))
	}

	d.Set("schema", schema.Schema())
	d.Set("schema_id", schema.ID())
	d.Set("version", schema.Version())

	d.SetId(formatSchemaVersionID(subject))

	return diags
}
