package routeros

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
  {
    ".id": "*2",
    "address": "192.168.88.5",
    "comment": "static",
    "complete": "false",
    "dhcp": "false",
    "disabled": "false",
    "dynamic": "false",
    "interface": "bridge",
    "invalid": "true",
    "mac-address": "02:00:00:00:00:05",
    "published": "false",
    "vrf": "main"
  }
*/

// ResourceIpArp https://help.mikrotik.com/docs/spaces/ROS/pages/328222/ARP
func ResourceIpArp() *schema.Resource {
	resSchema := map[string]*schema.Schema{
		MetaResourcePath: PropResourcePath("/ip/arp"),
		MetaId:           PropId(Id),

		"address": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "IP address to be mapped.",
			ValidateFunc: validation.IsIPv4Address,
		},
		KeyComment: PropCommentRw,
		"complete": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether the entry is resolved.",
		},
		"dhcp": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether the entry was added by the DHCP server.",
		},
		KeyDisabled: PropDisabledRw,
		KeyDynamic:  PropDynamicRo,
		KeyInterface: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Interface name the IP address is assigned to.",
		},
		KeyInvalid: PropInvalidRo,
		KeyMacAddress: {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "MAC address to be mapped to.",
			ValidateFunc:     validation.IsMACAddress,
			DiffSuppressFunc: MacAddressEqual,
		},
		"published": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Static proxy-arp entry for individual IP address.",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Neighbor state of the entry.",
		},
		KeyVrf: PropVrfRo,
	}

	return &schema.Resource{
		CreateContext: DefaultCreate(resSchema),
		ReadContext:   DefaultRead(resSchema),
		UpdateContext: DefaultUpdate(resSchema),
		DeleteContext: DefaultDelete(resSchema),

		Importer: &schema.ResourceImporter{
			StateContext: ImportStateCustomContext(resSchema),
		},

		Schema: resSchema,
	}
}
