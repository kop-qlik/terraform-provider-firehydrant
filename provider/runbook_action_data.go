package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRunbookAction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantRunbookAction,
		Schema: map[string]*schema.Schema{
			// Required
			"integration_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantRunbookAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the runbook action
	runbookType := d.Get("type").(string)
	actionSlug := d.Get("slug").(string)
	integrationSlug := d.Get("integration_slug").(string)
	runbookActionResponse, err := firehydrantAPIClient.RunbookActions().Get(ctx, runbookType, integrationSlug, actionSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the attributes in state to the values we got from the API
	attributes := map[string]string{
		"name": runbookActionResponse.Name,
		"slug": runbookActionResponse.Slug,
	}

	if runbookActionResponse.Integration != nil {
		attributes["integration_slug"] = runbookActionResponse.Integration.Slug
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set the runbook action's ID in state
	d.SetId(runbookActionResponse.ID)

	return diag.Diagnostics{}
}
