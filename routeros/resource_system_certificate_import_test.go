package routeros

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testCertificateImportKey = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIHeMEkGCSqGSIb3DQEFDTA8MBsGCSqGSIb3DQEFDDAOBAiy/wEW6/MglgICCAAw
HQYJYIZIAWUDBAEqBBD6v8dLA2FjPn62Xz57pcu9BIGQhclivPw1eC2b14ea58Tw
nzDdbYN6/yUiMqapW2xZaT7ZFnbEai4n9/utgtEDnfKHlZvZj2kRhvYoWrvTkt/W
1mkd5d/runsn+B5GO+CMHFHh4t41WMpZysmg+iP8FiiehOQEsWyEZFaedxfYYtSL
Sk+abxJ+NMQoh+S5d73niu1CO8uqQjOd8BoSOurURsOh
-----END ENCRYPTED PRIVATE KEY-----
`

const testCertificateImportCert = `-----BEGIN CERTIFICATE-----
MIIBlTCCATugAwIBAgIINLsws71B5zIwCgYIKoZIzj0EAwIwHzEdMBsGA1UEAwwU
RXh0ZXJuYWwgQ2VydGlmaWNhdGUwHhcNMjQwNTE3MjEyOTUzWhcNMjUwNTE3MjEy
OTUzWjAfMR0wGwYDVQQDDBRFeHRlcm5hbCBDZXJ0aWZpY2F0ZTBZMBMGByqGSM49
AgEGCCqGSM49AwEHA0IABKE1g0Qj4ujIold9tklu2z4BUu/K7xDFF5YmedtOfJyM
1/80APNboqn71y4m4XNE1JNtQuR2bSZPHVrzODkR16ujYTBfMA8GA1UdEwEB/wQF
MAMBAf8wDgYDVR0PAQH/BAQDAgG2MB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEF
BQcDAjAdBgNVHQ4EFgQUNXd5bvluIV9YAhGc5yMHc6OzXpMwCgYIKoZIzj0EAwID
SAAwRQIhAODte/qS6CE30cvnQpxP/ObWBPIPZnHtkFHIIC1AOSXwAiBGCGQE+aJY
W72Rw0Y1ckvlt6sU0urkzGuj5wxVF/gSYA==
-----END CERTIFICATE-----
`

var reTemporaryFileName = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

const testSystemCertificatesImportAddress = "routeros_system_certificate.external"

func TestAccSystemCertificatesTest_import(t *testing.T) {
	// ROS 7.12.x does not return IDs for created files.
	if !testCheckMinVersion(t, testFileMinVersion) {
		t.Logf("Test skipped, the minimum required version is %v", testFileMinVersion)
		return
	}

	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			var externalCrt MikrotikItem

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				CheckDestroy: resource.ComposeTestCheckFunc(
					testCheckResourceDestroy("/certificate", "routeros_system_certificate"),
					testCheckNoTemporaryFiles,
				),
				Steps: []resource.TestStep{
					{
						Config: testAccSystemCertificatesImportConfig(),
						Check: resource.ComposeTestCheckFunc(
							testResourcePrimaryInstanceId(testSystemCertificatesImportAddress),
							// external_crt
							testCheckResourceExists("routeros_system_certificate.external", "/certificate", &externalCrt),
							testCheckMikrotikItemAttr("routeros_system_certificate.external", &externalCrt, "name", "external.crt"),
							resource.TestCheckResourceAttr("routeros_system_certificate.external", "name", "external.crt"),
							resource.TestCheckResourceAttr("routeros_system_certificate.external", "common_name", "External Certificate"),
							resource.TestCheckResourceAttr("routeros_system_certificate.external", "private_key", "true"),
							testCheckNoTemporaryFiles,
						),
					},
				},
			})

		})
	}
}

func TestAccSystemCertificatesTest_importFileName(t *testing.T) {
	if !testCheckMinVersion(t, testFileMinVersion) {
		t.Logf("Test skipped, the minimum required version is %v", testFileMinVersion)
		return
	}

	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				CheckDestroy:      testCheckResourceDestroy("/certificate", "routeros_system_certificate"),
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							testUploadFile(t, "external-name.crt", testCertificateImportCert)
							testUploadFile(t, "external-name.key", testCertificateImportKey)
						},
						Config: testAccSystemCertificatesImportFileNameConfig(),
						Check: resource.ComposeTestCheckFunc(
							testResourcePrimaryInstanceId(testSystemCertificatesImportAddress),
							resource.TestCheckResourceAttr(testSystemCertificatesImportAddress, "common_name", "External Certificate"),
							resource.TestCheckResourceAttr(testSystemCertificatesImportAddress, "private_key", "true"),
						),
					},
				},
			})
		})
	}
}

func TestAccSystemCertificatesTest_importWrongPassphrase(t *testing.T) {
	if !testCheckMinVersion(t, testFileMinVersion) {
		t.Logf("Test skipped, the minimum required version is %v", testFileMinVersion)
		return
	}

	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				CheckDestroy: resource.ComposeTestCheckFunc(
					testCheckResourceDestroy("/certificate", "routeros_system_certificate"),
					testCheckNoTemporaryFiles,
				),
				Steps: []resource.TestStep{
					{
						Config:      testAccSystemCertificatesImportWrongPassphraseConfig(),
						ExpectError: regexp.MustCompile(`nothing was imported from file .*private-keys-imported=0`),
					},
				},
			})
		})
	}
}

func testUploadFile(t *testing.T, name, contents string) {
	host := reHost.FindStringSubmatch(os.Getenv("ROS_HOSTURL"))[1]
	c := newRestClient(context.Background(), "https://"+host+":443", os.Getenv("ROS_USERNAME"), os.Getenv("ROS_PASSWORD"))
	if _, err := CreateItem(context.Background(), MikrotikItem{"name": name, "contents": contents}, "/file", c); err != nil {
		t.Fatalf("upload %v: %v", name, err)
	}
}

func testCheckNoTemporaryFiles(*terraform.State) error {
	res, err := ReadItems(nil, "/file", testAccProvider.Meta().(Client))
	if err != nil {
		return err
	}
	for _, f := range *res {
		if reTemporaryFileName.MatchString(f["name"]) {
			return fmt.Errorf("temporary file %q is left on the router", f["name"])
		}
	}
	return nil
}

func testAccSystemCertificatesImportConfig() string {
	return providerConfig + `
data "routeros_x509" "cert" {
	data = <<EOT
` + testCertificateImportCert + `EOT
}

resource "routeros_system_certificate" "external" {
	name        = "external.crt"
	common_name = data.routeros_x509.cert.common_name
	import {
		cert_file_content = data.routeros_x509.cert.pem
		key_file_content  = <<EOT
` + testCertificateImportKey + `EOT
		passphrase        = "11111111"
	}
}
`
}

func testAccSystemCertificatesImportFileNameConfig() string {
	return providerConfig + `
resource "routeros_system_certificate" "external" {
	name        = "external.crt"
	common_name = "External Certificate"
	import {
		cert_file_name = "external-name.crt"
		key_file_name  = "external-name.key"
		passphrase     = "11111111"
	}
}
`
}

func testAccSystemCertificatesImportWrongPassphraseConfig() string {
	return providerConfig + `
resource "routeros_system_certificate" "external" {
	name        = "external.crt"
	common_name = "External Certificate"
	import {
		cert_file_content = <<EOT
` + testCertificateImportCert + `EOT
		key_file_content  = <<EOT
` + testCertificateImportKey + `EOT
		passphrase        = "wrong"
	}
}
`
}
