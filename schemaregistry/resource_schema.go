package schemaregistry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/riferrei/srclient"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: schemaCreate,
		UpdateContext: schemaUpdate,
		ReadContext:   schemaRead,
		DeleteContext: schemaDelete,
		Schema: map[string]*schema.Schema{
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subject related to the schema",
				ForceNew:    true,
			},
			"schema": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The schema string",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newJSON, _ := structure.NormalizeJsonString(new)
					oldJSON, _ := structure.NormalizeJsonString(old)
					return newJSON == oldJSON
				},
			},
			"schema_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the schema",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The schema string",
			},
		},
	}
}

func schemaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	subject := d.Get("subject").(string)
	schemaString := d.Get("schema").(string)

	client := meta.(*srclient.SchemaRegistryClient)

	schema, err := client.CreateSchemaWithArbitrarySubject(subject, schemaString, srclient.Avro)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(formatSchemaVersionID(subject))
	d.Set("schema_id", schema.ID())
	d.Set("schema", schema.Schema())
	d.Set("version", schema.Version())

	return diags
}

func schemaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	subject := d.Get("subject").(string)
	schemaString := d.Get("schema").(string)

	client := meta.(*srclient.SchemaRegistryClient)
	schema, err := client.CreateSchemaWithArbitrarySubject(subject, schemaString, srclient.Avro)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("schema_id", schema.ID())
	d.Set("schema", schema.Schema())
	d.Set("version", schema.Version())

	return diags
}

func schemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*srclient.SchemaRegistryClient)
	subject := extractSchemaVersionID(d.Id())

	schema, err := client.GetLatestSchemaWithArbitrarySubject(subject)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("schema", schema.Schema())
	d.Set("schema_id", schema.ID())
	d.Set("subject", subject)
	d.Set("version", schema.Version())

	return diags
}

func schemaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*srclient.SchemaRegistryClient)
	subject := extractSchemaVersionID(d.Id())

	err := client.DeleteSubject(subject, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
