package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/logging/v2"
)

var LoggingExclusionBaseSchema = map[string]*schema.Schema{
	"filter": {
		Type:     schema.TypeString,
		Required: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"description": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"disabled": {
		Type:     schema.TypeBool,
		Optional: true,
	},
}

func ResourceLoggingExclusion(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) *schema.Resource {
	return &schema.Resource{
		Create: ResourceLoggingExclusionCreate(newUpdaterFunc),
		Read:   ResourceLoggingExclusionRead(newUpdaterFunc),
		Update: ResourceLoggingExclusionUpdate(newUpdaterFunc),
		Delete: ResourceLoggingExclusionDelete(newUpdaterFunc),

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: mergeSchemas(LoggingExclusionBaseSchema, parentSpecificSchema),
	}
}

func ResourceLoggingExclusionCreate(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		id, exclusion := expandResourceLoggingExclusion(d, updater.GetResourceType(), updater.GetResourceId())

		err = updater.CreateLoggingExclusion(id.parent(), exclusion)
		if err != nil {
			return err
		}

		d.SetId(id.canonicalId())

		return ResourceLoggingExclusionRead(newUpdaterFunc)(d, meta)
	}
}

func ResourceLoggingExclusionRead(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		exclusion, err := updater.ReadLoggingExclusion(d.Id())

		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Logging Exclusion %s", d.Get("name").(string)))
		}

		flattenResourceLoggingExclusion(d, exclusion)

		if updater.GetResourceType() == "projects" {
			d.Set("project", updater.GetResourceId())
		}

		return nil
	}
}

func ResourceLoggingExclusionUpdate(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.UpdateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		exclusion, updateMask := expandResourceLoggingExclusionForUpdate(d)

		err = updater.UpdateLoggingExclusion(d.Id(), exclusion, updateMask)
		if err != nil {
			return err
		}

		return ResourceLoggingExclusionRead(newUpdaterFunc)(d, meta)
	}
}

func ResourceLoggingExclusionDelete(newUpdaterFunc newResourceLoggingExclusionUpdaterFunc) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		err = updater.DeleteLoggingExclusion(d.Id())
		if err != nil {
			return err
		}

		d.SetId("")
		return nil
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

// The ResourceLoggingExclusionUpdater interface is implemented for each GCP
// resource supporting log exclusions.
//
// Implementations should keep track of the resource identifier.
type ResourceLoggingExclusionUpdater interface {
	CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error
	ReadLoggingExclusion(id string) (*logging.LogExclusion, error)
	UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error
	DeleteLoggingExclusion(id string) error

	GetResourceType() string

	// Returns the unique resource identifier.
	GetResourceId() string

	// Textual description of this resource to be used in error message.
	// The description should include the unique resource identifier.
	DescribeResource() string
}

type newResourceLoggingExclusionUpdaterFunc func(d *schema.ResourceData, config *Config) (ResourceLoggingExclusionUpdater, error)

// loggingExclusionResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingExclusionResourceTypes = []string{
	"billingAccount",
	"folders",
	"organizations",
	"projects",
}

// LoggingExclusionId represents the parts that make up the canonical id used within terraform for a logging resource.
type LoggingExclusionId struct {
	resourceType string
	resourceId   string
	name         string
}

// canonicalId returns the LoggingExclusionId as the canonical id used within terraform.
func (l LoggingExclusionId) canonicalId() string {
	return fmt.Sprintf("%s/%s/exclusions/%s", l.resourceType, l.resourceId, l.name)
}

// parent returns the "parent-level" resource that the exclusion is in (e.g. `folders/foo` for id `folders/foo/exclusions/bar`)
func (l LoggingExclusionId) parent() string {
	return fmt.Sprintf("%s/%s", l.resourceType, l.resourceId)
}
