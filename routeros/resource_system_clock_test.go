package routeros

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testSystemClockTask = "routeros_system_clock.set"

func TestAccSystemClockTest_basic(t *testing.T) {
	for _, name := range testNames {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					testAccPreCheck(t)
					testSetTransportEnv(t, name)
				},
				ProviderFactories: testAccProviderFactories,
				Steps:             makeSteps(name),
			})

		})
	}
}

func makeSteps(name string) (res []resource.TestStep) {
	// Newer RouterOS versions refuse a time before the package build time.
	now := time.Now().UTC()
	params := map[string]map[string]string{
		"API": {
			"date":           now.AddDate(0, 0, 1).Format("2006-01-02"),
			"time":           `17:58:11`,
			"time_zone_name": `EST`,
		},
		"REST": {
			"date":           now.AddDate(0, 0, 2).Format("2006-01-02"),
			"time":           `18:58:11`,
			"time_zone_name": `UTC`,
		},
	}

	for k, v := range params[name] {
		res = append(res, resource.TestStep{
			Config: fmt.Sprintf(`%v
			resource "routeros_system_clock" "set" {
				%v = "%v"
			}`, providerConfig, k, v),
			Check: resource.ComposeTestCheckFunc(
				testResourcePrimaryInstanceId(testSystemClockTask),
				resource.TestCheckResourceAttr(testSystemClockTask, k, v),
			),
		})

	}
	return
}
