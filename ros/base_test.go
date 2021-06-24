package ros

import (
	"encoding/json"
	"net"
	"testing"
)

type update struct {
	Bridge   *string     `json:"bridge,omitempty"`
	Disabled *Boolean    `json:"disabled,omitempty"`
	Tagged   *StringList `json:"tagged,omitempty"`
	Untagged *StringList `json:"untagged,omitempty"`
	VlanIDs  *Number     `json:"vlan-ids,omitempty"`
	Network  *IP         `json:"network,omitempty"`
	Address  *IPNet      `json:"address,omitempty"`
}

// TestSerialize ensures that all ROS types serialize as expected.
func TestSerialize(t *testing.T) {
	for i, te := range []struct {
		u    update
		want string
	}{
		{update{}, `{}`},
		{update{
			Bridge: StringPtr("bridge1"),
		}, `{"bridge":"bridge1"}`},
		{update{
			Bridge:   StringPtr("bridge1"),
			Disabled: BooleanPtr(false),
		}, `{"bridge":"bridge1","disabled":"false"}`},
		{update{
			Bridge: StringPtr("bridge1"),
			Tagged: StringListPtr("ether1", "ether2"),
		}, `{"bridge":"bridge1","tagged":"ether1,ether2"}`},
		{update{
			Bridge:  StringPtr("bridge1"),
			VlanIDs: NumberPtr(1337),
		}, `{"bridge":"bridge1","vlan-ids":"1337"}`},
		{update{
			Bridge:  StringPtr("bridge1"),
			Network: IPPtr(net.ParseIP("1.2.3.4")),
		}, `{"bridge":"bridge1","network":"1.2.3.4"}`},
		{update{
			Bridge: StringPtr("bridge1"),
			Address: &IPNet{
				Address: net.ParseIP("1.2.3.4"),
				Network: net.IPNet{
					IP:   net.ParseIP("1.2.3.0"),
					Mask: net.CIDRMask(24, 32),
				},
			},
		}, `{"bridge":"bridge1","address":"1.2.3.4/24"}`},
	} {
		gotBytes, err := json.Marshal(te.u)
		if err != nil {
			t.Errorf("%d: marshal: %v", i, err)
			continue
		}
		if want, got := te.want, string(gotBytes); want != got {
			t.Errorf("%d: wanted %q, got %q", i, want, got)
		}
	}
}
