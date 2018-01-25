package google

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/logging/v2"
)

var ProjectLoggingExclusionSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	},
}

type ProjectLoggingExclusionUpdater struct {
	resourceType string
	resourceId   string
	Config       *Config
}

func NewProjectLoggingExclusionUpdater(d *schema.ResourceData, config *Config) (ResourceLoggingExclusionUpdater, error) {
	pid, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	return &ProjectLoggingExclusionUpdater{
		resourceType: "projects",
		resourceId:   pid,
		Config:       config,
	}, nil
}

func (u *ProjectLoggingExclusionUpdater) CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error {
	_, err := u.Config.clientLogging.Projects.Exclusions.Create(parent, exclusion).Do()
	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error creating logging exclusion for %s.", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectLoggingExclusionUpdater) ReadLoggingExclusion(id string) (*logging.LogExclusion, error) {
	exclusion, err := u.Config.clientLogging.Projects.Exclusions.Get(id).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving logging exclusion for %s: %s", u.DescribeResource(), err)
	}

	return exclusion, nil
}

func (u *ProjectLoggingExclusionUpdater) UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error {
	_, err := u.Config.clientLogging.Projects.Exclusions.Patch(id, exclusion).UpdateMask(updateMask).Do()
	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error updating logging exclusion for %s.", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectLoggingExclusionUpdater) DeleteLoggingExclusion(id string) error {
	_, err := u.Config.clientLogging.Projects.Exclusions.Delete(id).Do()
	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error deleting logging exclusion for %s.", u.DescribeResource()), err)
	}

	return nil
}

func (u *ProjectLoggingExclusionUpdater) GetResourceType() string {
	return u.resourceType
}

func (u *ProjectLoggingExclusionUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *ProjectLoggingExclusionUpdater) DescribeResource() string {
	return fmt.Sprintf("%q %q", u.resourceType, u.resourceId)
}
