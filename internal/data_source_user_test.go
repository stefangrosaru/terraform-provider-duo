package internal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("duo_user.test", "username", "data.duo_user.test", "username"),
				),
			},
		},
	})
}

const testAccDataSourceUser = `
resource "duo_user" "test" {
	username = "test@test.com"
}

data "duo_user" "test" {
	user_id = duo_user.test.id
}
`
