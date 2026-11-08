package routeros

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testSystemLoggingAction = "routeros_system_logging_action.action"

func TestAccSystemLoggingActionTest_basic(t *testing.T) {
	// RouterOS 7.18 replaced bsd-syslog with remote-log-format.
	format, formatValue := "bsd_syslog", "true"
	if testCheckMinVersion(t, "7.18") {
		format, formatValue = "remote_log_format", "syslog"
	}

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
						Config: testAccSystemLoggingActionConfig(format, formatValue),
						Check: resource.ComposeTestCheckFunc(
							testResourcePrimaryInstanceId(testSystemLoggingAction),
							resource.TestCheckResourceAttr(testSystemLoggingAction, "name", "action1"),
							resource.TestCheckResourceAttr(testSystemLoggingAction, "default", "false"),
							resource.TestCheckResourceAttr(testSystemLoggingAction, "target", "remote"),
							resource.TestCheckResourceAttr(testSystemLoggingAction, "remote", "192.168.1.1"),
							resource.TestCheckResourceAttr(testSystemLoggingAction, format, formatValue),
							resource.TestCheckResourceAttr(testSystemLoggingAction, "syslog_facility", "user"),
							resource.TestCheckResourceAttr(testSystemLoggingAction, "syslog_severity", "notice"),
							resource.TestCheckResourceAttr(testSystemLoggingAction, "syslog_time_format", "iso8601"),
						),
					},
				},
			})

		})
	}
}

func testAccSystemLoggingActionConfig(format, formatValue string) string {
	value := formatValue
	if format == "remote_log_format" {
		value = `"` + formatValue + `"`
	}
	return fmt.Sprintf(`%v
resource "routeros_system_logging_action" "action" {
	name               = "action1"
	target             = "remote"
	remote             = "192.168.1.1"
	%v = %v
	syslog_facility    = "user"
	syslog_severity    = "notice"
	syslog_time_format = "iso8601"
}
`, providerConfig, format, value)
}
