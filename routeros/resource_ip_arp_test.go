package routeros

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testIpArpAddress = "routeros_ip_arp.test"

func TestAccIpArpTest_basic(t *testing.T) {
	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				CheckDestroy:      testCheckResourceDestroy("/ip/arp", "routeros_ip_arp"),
				Steps: []resource.TestStep{
					{
						Config: testAccIpArpConfig("02:00:00:00:00:0a", "test", false),
						Check: resource.ComposeTestCheckFunc(
							testResourcePrimaryInstanceId(testIpArpAddress),
							resource.TestCheckResourceAttr(testIpArpAddress, "address", "192.168.88.10"),
							resource.TestCheckResourceAttr(testIpArpAddress, "interface", "bridge"),
							resource.TestCheckResourceAttr(testIpArpAddress, "dynamic", "false"),
						),
					},
					{
						Config: testAccIpArpConfig("02:00:00:00:00:0B", "changed", true),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(testIpArpAddress, "mac_address", "02:00:00:00:00:0B"),
							resource.TestCheckResourceAttr(testIpArpAddress, "comment", "changed"),
							resource.TestCheckResourceAttr(testIpArpAddress, "disabled", "true"),
						),
					},
					{
						ResourceName:            testIpArpAddress,
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"complete", "invalid", "status"},
					},
				},
			})
		})
	}
}

func testAccIpArpConfig(mac, comment string, disabled bool) string {
	return providerConfig + `
resource "routeros_ip_arp" "test" {
  address     = "192.168.88.10"
  interface   = "bridge"
  mac_address = "` + mac + `"
  comment     = "` + comment + `"
  disabled    = ` + strconv.FormatBool(disabled) + `
}
`
}
