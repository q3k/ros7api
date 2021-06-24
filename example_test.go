package ros7api

import (
	"context"
	"log"
	"os"

	"github.com/q3k/ros7api/ros"
)

func Example() {
	ctx := context.Background()
	c := ros.Client{
		Username: "admin",
		Password: os.Getenv("ROS7API_EXAMPLE_PASSWORD"),
		Address:  os.Getenv("ROS7API_EXAMPLE_ADDRESS"),
		HTTP:     ros.LetsEncryptClient,
	}

	// Get all VLANs.
	vlist, err := c.InterfaceBridgeVlanList(ctx)
	if err != nil {
		log.Fatalf("Could not list vlans: %v", err)
	}

	// Find VLAN 3005.
	var vl3005 *ros.InterfaceBridgeVlan
	for _, vlan := range vlist {
		if vlan.VlanIDs != 3005 {
			continue
		}
		vlan := vlan
		vl3005 = &vlan
	}
	if vl3005 == nil {
		log.Fatalf("No vlan 3005")
	}

	// Add ether8 if needed.
	add := true
	tagged := vl3005.Tagged
	for _, t := range tagged {
		if t == "ether8" {
			add = false
			break
		}
	}

	if add {
		tagged = append(tagged, "ether8")
		_, err = c.InterfaceBridgeVlanPatch(ctx, vl3005.ID, &ros.InterfaceBridgeVlan_Update{
			Tagged: &tagged,
		})
		if err != nil {
			log.Fatalf("Could not update vlan: %v", err)
		}
	}

	// Output:
}
