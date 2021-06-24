package ros

import (
	"context"
	"encoding/json"
	"fmt"
)

// Automatically generated by github.com/q3k/ros7api/gen, do not edit.

type InterfaceBridgeVlan struct {
	Record

	Bridge          string     `json:"bridge"`
	Disabled        Boolean    `json:"disabled"`
	Tagged          StringList `json:"tagged"`
	Untagged        StringList `json:"untagged"`
	VlanIDs         Number     `json:"vlan-ids"`
	CurrentTagged   StringList `json:"current-tagged"`
	CurrentUntagged StringList `json:"current-untagged"`
	Dynamic         Boolean    `json:"dynamic"`
}

type InterfaceBridgeVlan_Update struct {
	Bridge   *string     `json:"bridge,omitempty"`
	Disabled *Boolean    `json:"disabled,omitempty"`
	Tagged   *StringList `json:"tagged,omitempty"`
	Untagged *StringList `json:"untagged,omitempty"`
	VlanIDs  *Number     `json:"vlan-ids,omitempty"`
}

func (c *Client) InterfaceBridgeVlanList(ctx context.Context) ([]InterfaceBridgeVlan, error) {
	body, err := c.doGET(ctx, "interface/bridge/vlan")
	if err != nil {
		return nil, fmt.Errorf("could not GET: %w", err)
	}
	defer body.Close()

	var target []InterfaceBridgeVlan
	if err := json.NewDecoder(body).Decode(&target); err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}
	return target, nil
}

func (c *Client) InterfaceBridgeVlanPatch(ctx context.Context, id RecordID, u *InterfaceBridgeVlan_Update) (*InterfaceBridgeVlan, error) {
	rdata, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("could not marshal update: %w", err)
	}
	body, err := c.doPATCH(ctx, "interface/bridge/vlan/"+string(id), rdata)
	if err != nil {
		return nil, fmt.Errorf("could not PATCH: %w", err)
	}
	defer body.Close()

	var target struct {
		InterfaceBridgeVlan
		Error   int64  `json:"error"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}
	if err := json.NewDecoder(body).Decode(&target); err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}
	if target.Error != 0 {
		return nil, fmt.Errorf("server error: %s: %s", target.Message, target.Detail)
	}
	return &target.InterfaceBridgeVlan, nil
}
