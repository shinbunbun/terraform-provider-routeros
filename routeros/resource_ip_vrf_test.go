package routeros

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpVrfTest_basic(t *testing.T) {
	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				CheckDestroy:      testCheckResourceDestroy("/ip/vrf", "routeros_ip_vrf"),
				Steps: []resource.TestStep{
					{
						Config: testAccIpVrfConfig(),
						Check: resource.ComposeTestCheckFunc(
							// A
							testResourcePrimaryInstanceId("routeros_ip_vrf.test_vrf_a"),
							resource.TestCheckResourceAttr("routeros_ip_vrf.test_vrf_a", "disabled", "true"),
							resource.TestCheckResourceAttr("routeros_ip_vrf.test_vrf_a", "name", "vrf_1"),
						),
					},
				},
			})

		})
	}
}

func testAccIpVrfConfig() string {
	return providerConfig + `
resource "routeros_interface_bridge" "vrf1" {
	name    = "vrf-br1"
}

resource "routeros_interface_bridge" "vrf2" {
	name    = "vrf-br2"
}

resource "routeros_ip_vrf" "test_vrf_a" {
	disabled 	= true
	name 		= "vrf_1"
	interfaces 	= ["vrf-br1", "vrf-br2"]
	depends_on  = [routeros_interface_bridge.vrf1, routeros_interface_bridge.vrf2]
}
`
}
