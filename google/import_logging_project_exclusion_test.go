package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLoggingProjectExclusion_importBasic(t *testing.T) {
	t.Parallel()

	exclusionName := "tf-test-exclusion-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccLoggingProjectExclusion_basic(exclusionName),
			},

			resource.TestStep{
				ResourceName:      "google_logging_project_exclusion.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
