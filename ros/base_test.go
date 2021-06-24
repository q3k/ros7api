package ros

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type update struct {
	Bridge   *string     `json:"bridge,omitempty"`
	Disabled *Boolean    `json:"disabled,omitempty"`
	Tagged   *StringList `json:"tagged,omitempty"`
	Untagged *StringList `json:"untagged,omitempty"`
	VlanIDs  *NumberList `json:"vlan-ids,omitempty"`
	Network  *IP         `json:"network,omitempty"`
	Address  *IPNet      `json:"address,omitempty"`
}

func numberList(t *testing.T, s string) *NumberList {
	t.Helper()
	v, err := ParseNumberList(s)
	if err != nil {
		t.Fatalf("could not parse list: %v", err)
	}
	return v
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
			VlanIDs: numberList(t, "1337"),
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

func TestNumberList(t *testing.T) {
	opt := cmp.AllowUnexported(NumberList{}, numberListRange{})

	got := NumberList{}
	if err := got.UnmarshalJSON([]byte(`"123,280-290,200-300,1004,1005"`)); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if want2, got2 := "123,280-290,200-300,1004,1005", got.String(); want2 != got2 {
		t.Errorf("serialized range should be %q, got %q", want2, got2)
	}
	want := NumberList{
		ranges: []numberListRange{
			{123, 123},
			{280, 290},
			{200, 300},
			{1004, 1004},
			{1005, 1005},
		},
	}
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Fatalf("diff: %s", diff)
	}

	got.optimize()
	want = NumberList{
		ranges: []numberListRange{
			{123, 123},
			{200, 300},
			{1004, 1005},
		},
	}
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Fatalf("diff: %s", diff)
	}

	got.Add(150)
	got.Add(124)
	got.Add(150)
	got.Add(150)
	got.Add(150)
	got.Add(124)
	got.Add(150)
	got.AddRange(299, 303)
	want = NumberList{
		ranges: []numberListRange{
			{123, 124},
			{150, 150},
			{200, 303},
			{1004, 1005},
		},
	}
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Fatalf("diff: %s", diff)
	}

	got.Remove(254)
	got.Add(254)
	got.Remove(254)
	want = NumberList{
		ranges: []numberListRange{
			{123, 124},
			{150, 150},
			{200, 253},
			{255, 303},
			{1004, 1005},
		},
	}
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Fatalf("diff: %s", diff)
	}

	if !got.Contains(253) {
		t.Errorf("should contain 253")
	}
	if got.Contains(254) {
		t.Errorf("should not contain 254")
	}
	if !got.Contains(255) {
		t.Errorf("should contain 255")
	}

	if want2, got2 := "123-124,150,200-253,255-303,1004-1005", got.String(); want2 != got2 {
		t.Errorf("serialized range should be %q, got %q", want2, got2)
	}
}
