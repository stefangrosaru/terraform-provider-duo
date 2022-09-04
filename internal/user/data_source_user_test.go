package user

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO - add test cases for data source user
func TestAccDataSourceScaffolding(t *testing.T) {
	t.Skip("data source not yet implemented")

	resource.UnitTest(t, resource.TestCase{
		//PreCheck:          func() { provider.TestAccPreCheck(t) },
		//ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceScaffolding,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.duo_user.test", "user_id", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccDataSourceScaffolding = `
data "duo_user" "test" {
  user_id = "XXXXXXXXXXXXXXXXXXXX"
}
`
