package routeros

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testInterfaceVethAddress = "routeros_interface_veth.test"

func TestAccInterfaceVethTest_basic(t *testing.T) {
	if !testCheckCommandExists(t, "/interface/veth") {
		t.Skip("/interface/veth is not available (RouterOS 7.21 and later need the container package)")
	}

	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				CheckDestroy:      testCheckResourceDestroy("/interface/veth", "routeros_interface_veth"),
				Steps: []resource.TestStep{
					{
						Config: testAccInterfaceVethConfig(),
						Check: resource.ComposeTestCheckFunc(
							testResourcePrimaryInstanceId(testInterfaceVethAddress),
							resource.TestCheckResourceAttr(testInterfaceVethAddress, "name", "veth-test"),
						),
					},
				},
			})
		})
	}
}

func testAccInterfaceVethConfig() string {
	return providerConfig + `

resource "routeros_interface_veth" "test" {
  name    = "veth-test"
  address = ["192.168.120.2/24"]
  gateway = "192.168.120.1"
  comment = "Virtual interface"
}
`
}
