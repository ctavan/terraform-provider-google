package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/logging/v2"

	"strings"
)

func resourceLoggingExclusionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"filter": {
			Type:     schema.TypeString,
			Required: true,
		},

		"disabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func expandResourceLoggingExclusion(d *schema.ResourceData, resourceType, resourceId string) (LoggingExclusionId, *logging.LogExclusion) {
	id := LoggingExclusionId{
		resourceType: resourceType,
		resourceId:   resourceId,
		name:         d.Get("name").(string),
	}

	exclusion := logging.LogExclusion{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Filter:      d.Get("filter").(string),
		Disabled:    d.Get("disabled").(bool),
	}
	return id, &exclusion
}

func flattenResourceLoggingExclusion(d *schema.ResourceData, exclusion *logging.LogExclusion) {
	d.Set("name", exclusion.Name)
	d.Set("description", exclusion.Description)
	d.Set("filter", exclusion.Filter)
	d.Set("disabled", exclusion.Disabled)
}

func expandResourceLoggingExclusionForUpdate(d *schema.ResourceData) (*logging.LogExclusion, string) {
	// Can update description/filter/disabled right now.
	exclusion := logging.LogExclusion{}

	var updateMaskArr []string

	if d.HasChange("description") {
		exclusion.Description = d.Get("description").(string)
		exclusion.ForceSendFields = append(exclusion.ForceSendFields, "Description")
		updateMaskArr = append(updateMaskArr, "description")
	}

	if d.HasChange("filter") {
		exclusion.Filter = d.Get("filter").(string)
		exclusion.ForceSendFields = append(exclusion.ForceSendFields, "Filter")
		updateMaskArr = append(updateMaskArr, "filter")
	}

	if d.HasChange("disabled") {
		exclusion.Disabled = d.Get("disabled").(bool)
		exclusion.ForceSendFields = append(exclusion.ForceSendFields, "Disabled")
		updateMaskArr = append(updateMaskArr, "disabled")
	}

	updateMask := strings.Join(updateMaskArr, ",")
	return &exclusion, updateMask
}
