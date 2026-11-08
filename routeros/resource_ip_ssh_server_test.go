package routeros

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testIpSSHServerSettings = "routeros_ip_ssh_server.test"

func TestAccIpSSHServerSettingsTest_basic(t *testing.T) {
	// RouterOS 7.17 removed allow-none-crypto.
	crypto := "allow_none_crypto"
	if testCheckMinVersion(t, "7.17") {
		crypto = "strong_crypto"
	}
	// RouterOS 7.21 replaced always-allow-password-login with the password-authentication enum.
	passwordLogin := !testCheckMinVersion(t, "7.21")

	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: testAccIpSSHServerSettingsConfig(crypto, passwordLogin, true, "both", 1024),
						Check: resource.ComposeTestCheckFunc(
							testResourcePrimaryInstanceId(testIpSSHServerSettings),
							resource.TestCheckResourceAttr(testIpSSHServerSettings, crypto, "true"),
							testCheckResourceAttrIf(passwordLogin, testIpSSHServerSettings, "always_allow_password_login", "true"),
							resource.TestCheckResourceAttr(testIpSSHServerSettings, "forwarding_enabled", "both"),
							resource.TestCheckResourceAttr(testIpSSHServerSettings, "host_key_size", "1024"),
						),
					},
					{
						Config: testAccIpSSHServerSettingsConfig(crypto, passwordLogin, false, "no", 2048),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(testIpSSHServerSettings, crypto, "false"),
							testCheckResourceAttrIf(passwordLogin, testIpSSHServerSettings, "always_allow_password_login", "false"),
							resource.TestCheckResourceAttr(testIpSSHServerSettings, "forwarding_enabled", "no"),
							resource.TestCheckResourceAttr(testIpSSHServerSettings, "host_key_size", "2048"),
						),
					},
				},
			})
		})
	}
}

func testAccIpSSHServerSettingsConfig(crypto string, passwordLogin, b bool, fwd string, kSize int) string {
	login := ""
	if passwordLogin {
		login = fmt.Sprintf("always_allow_password_login = %v", b)
	}
	return fmt.Sprintf(`%v

resource "routeros_ip_ssh_server" "test" {
	%v = %v
	%v
	forwarding_enabled          = "%v"
	host_key_size               = %v
}
`, providerConfig, crypto, b, login, fwd, kSize)
}

func testCheckResourceAttrIf(cond bool, name, key, value string) resource.TestCheckFunc {
	if !cond {
		return func(*terraform.State) error { return nil }
	}
	return resource.TestCheckResourceAttr(name, key, value)
}
