package user

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stefangrosaru/terraform-provider-duo/internal/acctest"
)

// TDOD - add test cases for data source user
func TestAccResourceUser(t *testing.T) {
	t.Skip("resource not yet implemented")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"duo_user.test", "username", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccResourceUser = `
resource "duo_user" "test" {
  username = "test@test.com"
}
`
