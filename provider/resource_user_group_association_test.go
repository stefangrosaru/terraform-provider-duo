package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUserGroupAssociation(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserGroupAssociation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("duo_group.test", "name", "test"),
					resource.TestCheckResourceAttr("duo_user.test", "username", "test@test.com"),
				),
			},
		},
	})
}

const testAccResourceUserGroupAssociation = `
resource "duo_user" "test" {
	username = "test@test.com"
}
  
resource "duo_group" "test" {
	name = "test"
}
  
resource "duo_user_group_association" "test" {
	group_id = duo_group.test.id
	user_id = duo_user.test.id
}
`
