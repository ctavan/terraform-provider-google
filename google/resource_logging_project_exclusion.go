package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLoggingProjectExclusion() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingProjectExclusionCreate,
		Read:   resourceLoggingProjectExclusionRead,
		Delete: resourceLoggingProjectExclusionDelete,
		Update: resourceLoggingProjectExclusionUpdate,
		Schema: resourceLoggingExclusionSchema(),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	schm.Schema["project"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	}
	return schm
}

func resourceLoggingProjectExclusionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	id, exclusion := expandResourceLoggingExclusion(d, "projects", project)

	_, err = config.clientLogging.Projects.Exclusions.Create(id.parent(), exclusion).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())

	return resourceLoggingProjectExclusionRead(d, meta)
}

func resourceLoggingProjectExclusionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	exclusion, err := config.clientLogging.Projects.Exclusions.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project Logging Exclusion %s", d.Get("name").(string)))
	}

	flattenResourceLoggingExclusion(d, exclusion)
	return nil
}

func resourceLoggingProjectExclusionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	exclusion, updateMask := expandResourceLoggingExclusionForUpdate(d)

	_, err := config.clientLogging.Projects.Exclusions.Patch(d.Id(), exclusion).UpdateMask(updateMask).Do()
	if err != nil {
		return err
	}

	return resourceLoggingProjectExclusionRead(d, meta)
}

func resourceLoggingProjectExclusionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientLogging.Projects.Exclusions.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
