package ros

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type RecordID string

type Record struct {
	ID RecordID `json:".id"`
}

func StringPtr(s string) *string {
	return &s
}

type Number int64

func NumberPtr(n int64) *Number {
	v := Number(n)
	return &v
}

func (n *Number) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var v int64
	var err error
	if strings.HasPrefix(s, "0x") {
		v, err = strconv.ParseInt(s[2:], 16, 64)
	} else {
		v, err = strconv.ParseInt(s, 10, 64)
	}
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", s, err)
	}
	*n = Number(v)
	return nil
}

func (n *Number) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, *n)), nil
}

type Boolean bool

func BooleanPtr(b bool) *Boolean {
	v := Boolean(b)
	return &v
}

func (n *Boolean) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "true":
		*n = Boolean(true)
	case "false":
		*n = Boolean(false)
	default:
		return fmt.Errorf("invalid boolean: %q", s)
	}
	return nil
}

func (n *Boolean) MarshalJSON() ([]byte, error) {
	if *n {
		return []byte(`"true"`), nil
	} else {
		return []byte(`"false"`), nil
	}
}

type IPNet struct {
	Address net.IP
	Network net.IPNet
}

func (n *IPNet) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ip, net, err := net.ParseCIDR(s)
	if err != nil {
		return fmt.Errorf("invalid CIDR %q: %w", s, err)
	}
	*n = IPNet{
		Address: ip,
		Network: *net,
	}
	return nil
}

func (n *IPNet) MarshalJSON() ([]byte, error) {
	ones, _ := n.Network.Mask.Size()
	v := fmt.Sprintf(`"%s/%d"`, n.Address.String(), ones)
	return []byte(v), nil
}

type IP net.IP

func IPPtr(i net.IP) *IP {
	v := IP(i)
	return &v
}

func (n *IP) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ip := net.ParseIP(s)
	if ip == nil || ip.To4() == nil {
		return fmt.Errorf("invalid IP %q", s)
	}
	*n = IP(ip)
	return nil
}

func (n *IP) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", net.IP(*n).String())), nil
}

type StringList []string

func StringListPtr(l ...string) *StringList {
	v := StringList(l)
	return &v
}

func (n *StringList) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parts := strings.Split(s, ",")
	*n = parts
	return nil
}

func (n *StringList) MarshalJSON() ([]byte, error) {
	for i, el := range *n {
		if strings.Contains(el, ",") || strings.Contains(el, `"`) {
			return nil, fmt.Errorf("element %d of string list (%q) contains invalid character", i, el)
		}
	}
	res := fmt.Sprintf("%q", strings.Join(*n, ","))
	return []byte(res), nil
}
